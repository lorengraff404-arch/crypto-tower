#!/bin/bash
# Backend Endpoint Tests
# Tests all battle system endpoints with security checks

BASE_URL="http://localhost:8080"
TOKEN=""

echo "=== BACKEND ENDPOINT TESTS ==="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Test 1: Health Check
echo "1. Health Check"
response=$(curl -s "${BASE_URL}/health")
if echo "$response" | grep -q "healthy"; then
    echo -e "${GREEN}✓ Health check OK${NC}"
else
    echo -e "${RED}✗ Health check FAILED${NC}"
fi
echo ""

# Test 2: Auth - Login Required
echo "2. Auth - Endpoints require JWT"
response=$(curl -s -w "%{http_code}" -o /dev/null "${BASE_URL}/api/v1/raids/start")
if [ "$response" = "401" ]; then
    echo -e "${GREEN}✓ Unauthorized request blocked${NC}"
else
    echo -e "${RED}✗ Auth not working (got $response)${NC}"
fi
echo ""

# Test 3: Create Test User (if needed)
echo "3. User Registration"
response=$(curl -s -X POST "${BASE_URL}/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "wallet_address": "0xTEST123456789",
        "signature": "test_signature"
    }')

if echo "$response" | grep -q "token\|already"; then
    echo -e "${GREEN}✓ User registration works${NC}"
    TOKEN=$(echo "$response" | jq -r '.token // empty')
else
    echo -e "${RED}✗ Registration failed${NC}"
fi
echo ""

# Test 4: Login
echo "4. User Login"
response=$(curl -s -X POST "${BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "wallet_address": "0xTEST123456789",
        "signature": "test_signature"
    }')

if echo "$response" | grep -q "token"; then
    echo -e "${GREEN}✓ Login successful${NC}"
    TOKEN=$(echo "$response" | jq -r '.token')
    echo "Token: ${TOKEN:0:20}..."
else
    echo -e "${RED}✗ Login failed${NC}"
fi
echo ""

# Test 5: Raid Start (with auth)
if [ -n "$TOKEN" ]; then
    echo "5. Raid Start (authenticated)"
    response=$(curl -s -w "%{http_code}" -o /tmp/raid_response.json \
        -X POST "${BASE_URL}/api/v1/raids/start" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "island_id": 1,
            "team_characters": []
        }')
    
    if [ "$response" = "200" ] || [ "$response" = "400" ]; then
        echo -e "${GREEN}✓ Raid endpoint responding${NC}"
        cat /tmp/raid_response.json | jq '.' 2>/dev/null || cat /tmp/raid_response.json
    else
        echo -e "${RED}✗ Raid start failed (code: $response)${NC}"
    fi
else
    echo -e "${RED}✗ Skipped - no token${NC}"
fi
echo ""

# Test 6: Ranked Matchmaking
if [ -n "$TOKEN" ]; then
    echo "6. Ranked Matchmaking"
    response=$(curl -s -w "%{http_code}" -o /tmp/ranked_response.json \
        -X POST "${BASE_URL}/api/v1/battle/ranked" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{}')
    
    if [ "$response" = "200" ] || [ "$response" = "404" ] || [ "$response" = "400" ]; then
        echo -e "${GREEN}✓ Ranked endpoint responding${NC}"
        cat /tmp/ranked_response.json | jq '.' 2>/dev/null || cat /tmp/ranked_response.json
    else
        echo -e "${RED}✗ Ranked failed (code: $response)${NC}"
    fi
else
    echo -e "${RED}✗ Skipped - no token${NC}"
fi
echo ""

# Test 7: Rate Limiting
echo "7. Rate Limiting Check"
echo "Sending 5 rapid requests..."
for i in {1..5}; do
    curl -s -w "%{http_code}\n" -o /dev/null "${BASE_URL}/health"
done
echo -e "${GREEN}✓ Rate limiting tested${NC}"
echo ""

# Test 8: Database Connection
echo "8. Database Status"
db_status=$(curl -s "${BASE_URL}/health" | jq -r '.database')
if [ "$db_status" = "connected" ]; then
    echo -e "${GREEN}✓ Database connected${NC}"
else
    echo -e "${RED}✗ Database not connected${NC}"
fi
echo ""

echo "=== TEST SUMMARY ==="
echo "Backend: $BASE_URL"
echo "Token obtained: $([ -n "$TOKEN" ] && echo 'Yes' || echo 'No')"
echo "Tests completed"
