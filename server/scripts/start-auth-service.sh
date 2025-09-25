#!/bin/bash

# Start Auth Service with JWT Management
# This script sets up the environment and starts the auth service

set -e

echo "🚀 Starting Lornian Auth Service with JWT Management"

# Set default environment variables if not set
export JWT_SECRET="${JWT_SECRET:-your-super-secure-jwt-secret-change-this-in-production}"
export JWT_EXP="${JWT_EXP:-3600}"                    # 1 hour
export REFRESH_TOKEN_EXP="${REFRESH_TOKEN_EXP:-2592000}"  # 30 days
export PORT="${PORT:-8080}"
export DATABASE_URL="${DATABASE_URL:-postgresql://user:pass@localhost:5432/lornian}"

echo "📊 Configuration:"
echo "  - JWT Expiration: ${JWT_EXP} seconds"
echo "  - Refresh Token Expiration: ${REFRESH_TOKEN_EXP} seconds"
echo "  - Port: ${PORT}"
echo "  - Database URL: ${DATABASE_URL}"

# Build the service
echo "🔨 Building auth service..."
go build -o bin/auth-service ./services/auth-service/cmd/main.go

# Run migrations (if available)
echo "🗃️  Running database migrations..."
if [ -f "./run-migrations.sh" ]; then
    ./run-migrations.sh
else
    echo "   No migration script found, skipping..."
fi

# Start the service
echo "🎯 Starting auth service..."
./bin/auth-service

echo "✅ Auth service started successfully!"
