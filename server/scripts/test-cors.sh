 #!/bin/bash

# Test CORS functionality
echo "🧪 Testing CORS Configuration..."

# Test OPTIONS preflight request
echo "📡 Testing OPTIONS preflight request..."
curl -X OPTIONS \
  -H "Origin: https://opulent-potato-q7pq7jg65p9jhg57-3000.app.github.dev" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization" \
  -v \
  https://api.lornian.com/auth/register

echo -e "\n📡 Testing actual POST request..."
curl -X POST \
  -H "Origin: https://opulent-potato-q7pq7jg65p9jhg57-3000.app.github.dev" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123","username":"testuser","name":"Test User"}' \
  -v \
  https://api.lornian.com/auth/register

echo -e "\n📡 Testing from localhost..."
curl -X OPTIONS \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization" \
  -v \
  https://api.lornian.com/auth/register

echo -e "\n✅ CORS test completed!"
