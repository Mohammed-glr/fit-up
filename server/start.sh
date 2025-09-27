#!/bin/bash

# Fit-Up Server Startup Script
echo "🚀 Starting Fit-Up Server Services..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    if ! command -v docker &> /dev/null || ! docker --help | grep -q "compose"; then
        echo "❌ docker-compose not found. Please install docker-compose."
        exit 1
    fi
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

# Set default environment variables if not set
export GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID:-"your-google-client-id"}
export GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET:-"your-google-client-secret"}
export GITHUB_CLIENT_ID=${GITHUB_CLIENT_ID:-"your-github-client-id"}
export GITHUB_CLIENT_SECRET=${GITHUB_CLIENT_SECRET:-"your-github-client-secret"}

echo "📦 Building and starting services..."

# Build and start all services
$COMPOSE_CMD up --build -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check service health
echo "🔍 Checking service health..."

services=("api-gateway:8080" "auth-service:8081" "message-service:8082" "schema-service:8083")
all_healthy=true

for service in "${services[@]}"; do
    name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    
    if curl -f -s "http://localhost:$port/health" > /dev/null; then
        echo "✅ $name is healthy"
    else
        echo "❌ $name is not responding"
        all_healthy=false
    fi
done

if $all_healthy; then
    echo ""
    echo "🎉 All services are running successfully!"
    echo ""
    echo "📍 Service URLs:"
    echo "   API Gateway:     http://localhost:8080"
    echo "   Auth Service:    http://localhost:8081"
    echo "   Message Service: http://localhost:8082"
    echo "   Schema Service:  http://localhost:8083"
    echo "   PostgreSQL:      localhost:5432"
    echo ""
    echo "🔧 Useful commands:"
    echo "   View logs:       $COMPOSE_CMD logs -f [service-name]"
    echo "   Stop services:   $COMPOSE_CMD down"
    echo "   Restart:         $COMPOSE_CMD restart [service-name]"
    echo ""
    echo "📖 Test the API:"
    echo "   Health check:    curl http://localhost:8080/health"
    echo "   Auth endpoints:  curl http://localhost:8080/auth/login"
else
    echo ""
    echo "⚠️  Some services are not healthy. Check logs with:"
    echo "   $COMPOSE_CMD logs"
fi