#!/bin/bash

# Simple Azure VM startup script for FitUp Backend
echo "🚀 Starting FitUp Backend on Azure VM..."

COMPOSE_FILE="docker/docker-compose.azure.yml"
ENV_FILE="docker/.env.azure"

# Check if environment file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Environment file $ENV_FILE not found"
    echo "💡 Please copy docker/.env.azure.template to docker/.env.azure and configure it"
    exit 1
fi

# Stop any running services
echo "🛑 Stopping existing services..."
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down 2>/dev/null || true

# Build and start services
echo "🏗️  Building and starting services..."
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up --build -d

# Wait a moment for services to start
sleep 5

# Check service status
echo "📊 Service status:"
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps

# Show logs
echo "📋 Recent logs:"
docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs --tail=10

echo "✅ Services started! Your API is available at:"
echo "   🌐 HTTPS: https://$(curl -s ifconfig.me)"
echo "   🌐 HTTP: http://$(curl -s ifconfig.me) (redirects to HTTPS)"
echo "   🔍 Health check: https://$(curl -s ifconfig.me)/health"
echo ""
echo "💡 Useful commands:"
echo "   View logs: docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f"
echo "   Stop services: docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE down"
echo "   Restart: docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE restart"
