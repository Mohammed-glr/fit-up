#!/bin/bash

# Generate random user data for testing
TIMESTAMP=$(date +%s)
RANDOM_NUM=$((RANDOM % 1000))
TEST_EMAIL="testuser${TIMESTAMP}${RANDOM_NUM}@example.com"
TEST_USERNAME="testuser${TIMESTAMP}${RANDOM_NUM}"
TEST_PASSWORD="TestPassword123!"
TEST_NAME="Test User ${RANDOM_NUM}"

echo "=== Testing Complete Auth Flow ==="
echo "Test User: $TEST_USERNAME"
echo "Test Email: $TEST_EMAIL"
echo ""

# Step 1: Register new user
echo "1. Registering new user..."
REGISTER_RESPONSE=$(curl -s -X POST https://api.lornian.com/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\",
    \"name\": \"$TEST_NAME\"
  }")

echo "Register response: $REGISTER_RESPONSE"

# Check if registration was successful
if echo "$REGISTER_RESPONSE" | grep -q "error\|Error\|failed\|Failed"; then
  echo "❌ Registration failed!"
  echo "Response: $REGISTER_RESPONSE"
  exit 1
fi

echo "✅ Registration successful!"
echo ""

# Step 2: Login with the registered credentials
echo "2. Logging in with registered credentials..."
LOGIN_RESPONSE=$(curl -s -X POST https://api.lornian.com/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"identifier\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "Login response: $LOGIN_RESPONSE"

# Extract tokens from login response
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)

# Alternative token extraction methods (try different response formats)
if [ -z "$ACCESS_TOKEN" ]; then
  ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$ACCESS_TOKEN" ]; then
  ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$ACCESS_TOKEN" ]; then
  echo "❌ Failed to extract access token from login response"
  echo "Login response: $LOGIN_RESPONSE"
  echo ""
  echo "Trying to parse JSON differently..."
  # Try using jq if available
  if command -v jq &> /dev/null; then
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token // .token // .accessToken // empty')
    REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token // .refreshToken // empty')
  fi
fi

if [ -z "$ACCESS_TOKEN" ]; then
  echo "❌ Still couldn't extract token. Please check the login response format."
  exit 1
fi

echo "✅ Login successful!"
echo "Access Token: ${ACCESS_TOKEN:0:20}..."
if [ ! -z "$REFRESH_TOKEN" ]; then
  echo "Refresh Token: ${REFRESH_TOKEN:0:20}..."
fi
echo ""

# Step 3: Test logout
echo "3. Testing logout..."
echo "Making logout request with token..."

LOGOUT_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}\nTIME_TOTAL:%{time_total}" \
  -X POST https://api.lornian.com/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

# Parse the response
HTTP_STATUS=$(echo "$LOGOUT_RESPONSE" | grep "HTTP_STATUS:" | cut -d: -f2)
TIME_TOTAL=$(echo "$LOGOUT_RESPONSE" | grep "TIME_TOTAL:" | cut -d: -f2)
RESPONSE_BODY=$(echo "$LOGOUT_RESPONSE" | sed '/HTTP_STATUS:/,$d')

echo "Logout Response Body: $RESPONSE_BODY"
echo "HTTP Status: $HTTP_STATUS"
echo "Response Time: ${TIME_TOTAL}s"

# Check logout result
if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "204" ]; then
  echo "✅ Logout successful!"
elif [ "$HTTP_STATUS" = "500" ]; then
  echo "❌ Logout failed with 500 Internal Server Error"
  echo "This is the error you're experiencing!"
elif [ "$HTTP_STATUS" = "401" ]; then
  echo "❌ Logout failed with 401 Unauthorized"
  echo "Token might be invalid or expired"
else
  echo "❌ Logout failed with HTTP status: $HTTP_STATUS"
fi

echo ""

# Step 4: Verbose logout test for debugging
echo "4. Running verbose logout test for debugging..."
echo "=== VERBOSE CURL OUTPUT ==="
curl -X POST https://api.lornian.com/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -v

echo ""
echo "=== TEST COMPLETE ==="

# Optional: Test with invalid token
echo ""
echo "5. Testing logout with invalid token (should fail gracefully)..."
curl -s -w "HTTP Status: %{http_code}\n" \
  -X POST https://api.lornian.com/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid_token_12345" \
  -o /dev/null