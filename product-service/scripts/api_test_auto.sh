#!/usr/bin/env bash
# api_test_auto.sh - Automated quick validation (CI/CD ready)

BASE_URL="http://localhost:8080"
PASSED=0
FAILED=0

test_endpoint() {
  local name=$1
  local method=$2
  local endpoint=$3
  local data=$4
  local expected_code=$5

  if [ -n "$data" ]; then
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$BASE_URL$endpoint" \
      -H "Content-Type: application/json" -d "$data")
  else
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$BASE_URL$endpoint")
  fi

  if [ "$STATUS" == "$expected_code" ]; then
    echo "[PASS] $name ($method $endpoint) -> $STATUS"
    ((PASSED++))
  else
    echo "[FAIL] $name ($method $endpoint) -> Expected $expected_code, got $STATUS"
    ((FAILED++))
  fi
}

echo "=== QUICK API VALIDATION ==="
echo ""

# Check server
if ! curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products" | grep -q "200"; then
  echo "[ERROR] Server not running at $BASE_URL"
  exit 1
fi

# Run tests
test_endpoint "List products" "GET" "/products" "" "200"
test_endpoint "Get by ID" "GET" "/products/1" "" "200"
test_endpoint "Invalid ID (404)" "GET" "/products/99999" "" "404"
test_endpoint "Invalid ID (400)" "GET" "/products/abc" "" "400"
test_endpoint "Pagination" "GET" "/products/paged?page=1&limit=5" "" "200"
test_endpoint "Search" "GET" "/products/search?q=teh" "" "200"

# Create -> Update -> Delete flow
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name":"AutoTest","price":1000,"qty":5}')
TEST_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

if [ -n "$TEST_ID" ] && [ "$TEST_ID" != "null" ]; then
  echo "[PASS] Create product -> ID: $TEST_ID"
  ((PASSED++))

  # Update
  sleep 0.3 # Avoid rate limit
  test_endpoint "Update product" "PUT" "/products/$TEST_ID" \
    '{"name":"Updated","price":2000,"qty":10,"category":"test"}' "204"

  # Patch
  sleep 0.3 # Avoid rate limit
  test_endpoint "Patch product" "PATCH" "/products/$TEST_ID" \
    '{"price":3000}' "200"

  # Delete
  sleep 0.3 # Avoid rate limit
  test_endpoint "Delete product" "DELETE" "/products/$TEST_ID" "" "204"

  # Verify deletion
  sleep 0.3 # Avoid rate limit
  test_endpoint "Verify deletion" "GET" "/products/$TEST_ID" "" "404"
else
  echo "[FAIL] Create product failed"
  ((FAILED++))
fi

# Summary
echo ""
echo "=== RESULTS ==="
echo "Passed: $PASSED"
echo "Failed: $FAILED"

if [ "$FAILED" -eq 0 ]; then
  echo "✓ All tests passed!"
  exit 0
else
  echo "✗ Some tests failed!"
  exit 1
fi
