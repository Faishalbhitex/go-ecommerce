#!/usr/bin/env bash

# Colors
GREEN="\e[32m"
RED="\e[31m"
YELLOW="\e[33m"
BLUE="\e[34m"
CYAN="\e[36m"
RESET="\e[0m"

# Config
BASE_URL="http://localhost:8080"
PASS_COUNT=0
FAIL_COUNT=0

# Helper functions
print_header() {
  echo -e "\n${CYAN}========================================${RESET}"
  echo -e "${CYAN}$1${RESET}"
  echo -e "${CYAN}========================================${RESET}"
}

print_test() {
  echo -e "${BLUE}[TEST]${RESET} $1"
}

print_pass() {
  echo -e "${GREEN}[PASS]${RESET} $1"
  ((PASS_COUNT++))
}

print_fail() {
  echo -e "${RED}[FAIL]${RESET} $1"
  ((FAIL_COUNT++))
}

print_info() {
  echo -e "${YELLOW}[INFO]${RESET} $1"
}

# Test 1: CORS Headers
print_header "TEST 1: CORS Headers"
print_test "Checking CORS headers..."

CORS_RESPONSE=$(curl -sI "$BASE_URL/products")
if echo "$CORS_RESPONSE" | grep -q "Access-Control-Allow-Origin: \*" 2>/dev/null; then
  print_pass "CORS header 'Access-Control-Allow-Origin: *' found"
else
  print_fail "CORS header 'Access-Control-Allow-Origin' not found"
fi

if echo "$CORS_RESPONSE" | grep -q "Access-Control-Allow-Methods" 2>/dev/null; then
  print_pass "CORS header 'Access-Control-Allow-Methods' found"
else
  print_fail "CORS header 'Access-Control-Allow-Methods' not found"
fi

# Test 2: Get All Products
print_header "TEST 2: Get All Products"
print_test "Fetching all products..."

RESPONSE=$(curl -s "$BASE_URL/products")
TOTAL_PRODUCTS=$(echo "$RESPONSE" | jq '.data | length')

if [ "$TOTAL_PRODUCTS" -gt 0 ]; then
  print_pass "Retrieved $TOTAL_PRODUCTS products"
  print_info "First 3 products:"
  echo "$RESPONSE" | jq -r '.data[0:3] | .[] | "  - ID: \(.id) | \(.name) | Rp\(.price)"'
else
  print_fail "No products found"
fi

# Test 3: Get Single Product
print_header "TEST 3: Get Single Product"
print_test "Fetching product ID=1..."

SINGLE=$(curl -s "$BASE_URL/products/1")
PRODUCT_NAME=$(echo "$SINGLE" | jq -r '.name')

if [ "$PRODUCT_NAME" != "null" ] && [ -n "$PRODUCT_NAME" ]; then
  print_pass "Product found: $PRODUCT_NAME"
  echo "$SINGLE" | jq -r '"  ID: \(.id) | \(.name) | Rp\(.price) | Stock: \(.qty)"'
else
  print_fail "Product ID=1 not found"
fi

# Test 4: Invalid Product ID
print_header "TEST 4: Invalid Product ID"
print_test "Testing invalid ID (99999)..."

STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products/99999")
if [ "$STATUS_CODE" == "404" ]; then
  print_pass "Returns 404 for non-existent product"
else
  print_fail "Expected 404, got $STATUS_CODE"
fi

print_test "Testing non-numeric ID (abc)..."
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products/abc")
if [ "$STATUS_CODE" == "400" ]; then
  print_pass "Returns 400 for invalid ID format"
else
  print_fail "Expected 400, got $STATUS_CODE"
fi

# Test 5: Pagination
print_header "TEST 5: Pagination"
print_test "Testing pagination (page=1, limit=5)..."

PAGED=$(curl -s "$BASE_URL/products/paged?page=1&limit=5")
PAGED_COUNT=$(echo "$PAGED" | jq '.data | length')

if [ "$PAGED_COUNT" -eq 5 ] || [ "$PAGED_COUNT" -le 5 ]; then
  print_pass "Pagination working: got $PAGED_COUNT items"
  print_info "Paginated products:"
  echo "$PAGED" | jq -r '.data[] | "  - ID: \(.id) | \(.name)"'
else
  print_fail "Expected max 5 items, got $PAGED_COUNT"
fi

print_test "Testing invalid pagination params..."
INVALID_PAGE=$(curl -s "$BASE_URL/products/paged?page=abc&limit=xyz")
INVALID_COUNT=$(echo "$INVALID_PAGE" | jq '.data | length')

if [ "$INVALID_COUNT" -gt 0 ]; then
  print_pass "Falls back to default pagination: $INVALID_COUNT items"
else
  print_fail "Failed to handle invalid pagination params"
fi

# Test 6: Search
print_header "TEST 6: Search Functionality"
print_test "Searching for 'teh'..."

SEARCH=$(curl -s "$BASE_URL/products/search?q=teh")
SEARCH_COUNT=$(echo "$SEARCH" | jq '.data | length')

if [ "$SEARCH_COUNT" -gt 0 ]; then
  print_pass "Search found $SEARCH_COUNT result(s)"
  echo "$SEARCH" | jq -r '.data[] | "  - \(.name)"'
else
  print_fail "Search returned no results"
fi

print_test "Searching with wildcard '%'..."
WILDCARD=$(curl -s "$BASE_URL/products/search?q=%25")
WILDCARD_COUNT=$(echo "$WILDCARD" | jq '.data | length')

if [ "$WILDCARD_COUNT" -eq "$TOTAL_PRODUCTS" ]; then
  print_pass "Wildcard search returns all products: $WILDCARD_COUNT"
else
  print_fail "Wildcard search failed (expected: $TOTAL_PRODUCTS, got: $WILDCARD_COUNT)"
fi

# Test 7: Rate Limiting
print_header "TEST 7: Rate Limiting"
print_test "Sending 20 rapid requests..."

SUCCESS=0
RATE_LIMITED=0

for i in {1..20}; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products")
  if [ "$STATUS" == "200" ]; then
    ((SUCCESS++))
  elif [ "$STATUS" == "429" ]; then
    ((RATE_LIMITED++))
  fi
done

print_info "Results: $SUCCESS success, $RATE_LIMITED rate-limited"

if [ "$RATE_LIMITED" -gt 0 ]; then
  print_pass "Rate limiting is working (got $RATE_LIMITED 429 responses)"
else
  print_fail "Rate limiting not triggered"
fi

print_test "Waiting 1.5s and retrying..."
sleep 1.5
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products")
if [ "$STATUS" == "200" ]; then
  print_pass "Request successful after cooldown"
else
  print_fail "Request still blocked after cooldown"
fi

# Test 8: Response Structure
print_header "TEST 8: Response Structure"
print_test "Validating product JSON structure..."

FIRST_PRODUCT=$(curl -s "$BASE_URL/products" | jq '.data[0]')
REQUIRED_FIELDS=("id" "name" "description" "price" "qty" "category" "create_at" "update_at")

for field in "${REQUIRED_FIELDS[@]}"; do
  if echo "$FIRST_PRODUCT" | jq -e "has(\"$field\")" >/dev/null 2>&1; then
    print_pass "Field '$field' exists"
  else
    print_fail "Field '$field' missing"
  fi
done

# Summary
print_header "TEST SUMMARY"
TOTAL=$((PASS_COUNT + FAIL_COUNT))
echo -e "${GREEN}Passed: $PASS_COUNT${RESET}"
echo -e "${RED}Failed: $FAIL_COUNT${RESET}"
echo -e "Total:  $TOTAL"

if [ "$FAIL_COUNT" -eq 0 ]; then
  echo -e "\n${GREEN}✓ All tests passed!${RESET}"
  exit 0
else
  echo -e "\n${RED}✗ Some tests failed!${RESET}"
  exit 1
fi
