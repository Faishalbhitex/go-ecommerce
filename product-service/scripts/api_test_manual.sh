#!/usr/bin/env bash
# api_test_manual.sh - Interactive API testing for humans

GREEN="\e[32m"
CYAN="\e[36m"
YELLOW="\e[33m"
RED="\e[31m"
RESET="\e[0m"

BASE_URL="http://localhost:8080"

# Check if server is running
if ! curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products" | grep -q "200"; then
  echo -e "${RED}[ERROR] Server not running!${RESET}"
  echo -e "${YELLOW}Start server first: ./scripts/dev.sh${RESET}"
  exit 1
fi

echo -e "${CYAN}=== MANUAL API TESTING ===${RESET}\n"

# 1. CREATE
echo -e "${GREEN}[1] CREATE Product${RESET}"
echo -e "${CYAN}POST /products${RESET}"
CREATED=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Product", "description": "Created via script", "price": 5000, "qty": 10, "category": "test"}')

CREATED_ID=$(echo "$CREATED" | jq -r '.id')
echo "$CREATED" | jq '.'
echo -e "${YELLOW}Created ID: $CREATED_ID${RESET}\n"

# 2. GET BY ID
echo -e "${GREEN}[2] GET Product by ID${RESET}"
echo -e "${CYAN}GET /products/$CREATED_ID${RESET}"
curl -s "$BASE_URL/products/$CREATED_ID" | jq '.'
echo ""

# 3. UPDATE (Full)
echo -e "${GREEN}[3] UPDATE Product (Full)${RESET}"
echo -e "${CYAN}PUT /products/$CREATED_ID${RESET}"
curl -s -X PUT "$BASE_URL/products/$CREATED_ID" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Updated Product\", \"description\": \"Full update\", \"price\": 6000, \"qty\": 15, \"category\": \"updated\"}" | jq '.'
echo ""

# 4. PATCH (Partial)
echo -e "${GREEN}[4] PATCH Product (Partial)${RESET}"
echo -e "${CYAN}PATCH /products/$CREATED_ID${RESET}"
curl -s -X PATCH "$BASE_URL/products/$CREATED_ID" \
  -H "Content-Type: application/json" \
  -d '{"price": 7000}' | jq '.'
echo ""

# 5. LIST
echo -e "${GREEN}[5] LIST All Products${RESET}"
echo -e "${CYAN}GET /products${RESET}"
curl -s "$BASE_URL/products" | jq '. | length as $count | "Total: \($count) products"'
echo ""

# 6. PAGINATION
echo -e "${GREEN}[6] PAGINATION${RESET}"
echo -e "${CYAN}GET /products/paged?page=1&limit=3${RESET}"
curl -s "$BASE_URL/products/paged?page=1&limit=3" | jq '.'
echo ""

# 7. SEARCH
echo -e "${GREEN}[7] SEARCH${RESET}"
echo -e "${CYAN}GET /products/search?q=test${RESET}"
curl -s "$BASE_URL/products/search?q=test" | jq '.'
echo ""

# 8. DELETE
echo -e "${GREEN}[8] DELETE Product${RESET}"
echo -e "${CYAN}DELETE /products/$CREATED_ID${RESET}"
curl -s -i -X DELETE "$BASE_URL/products/$CREATED_ID" | head -n 1
echo ""

# 9. Verify Deletion
echo -e "${GREEN}[9] VERIFY Deletion (should 404)${RESET}"
echo -e "${CYAN}GET /products/$CREATED_ID${RESET}"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products/$CREATED_ID")
if [ "$STATUS" == "404" ]; then
  echo -e "${GREEN}✓ Correctly returns 404${RESET}"
else
  echo -e "${RED}✗ Expected 404, got $STATUS${RESET}"
fi
echo ""

echo -e "${CYAN}=== TESTING COMPLETE ===${RESET}"
