#!/bin/bash

# Simple Azure VM startup script
echo "🚀 Starting Lornian Backend on Azure VM..."

# Stop any running services
echo "🛑 Stopping existing services..."
docker-compose down 2>/dev/null || true

# Build and start services
echo "🏗️  Building and starting services..."
docker-compose up --build -d

# Wait a moment for services to start
sleep 5

# Check service status
echo "📊 Service status:"
docker-compose ps

# Show logs
echo "📋 Recent logs:"
docker-compose logs --tail=10

echo "✅ Services started! Your API is available at:"
echo "   🌐 Main API: http://$(curl -s ifconfig.me):8080"
echo "   🔍 Health check: http://$(curl -s ifconfig.me):8080/health"
echo ""
echo "💡 Useful commands:"
echo "   View logs: docker-compose logs -f"
echo "   Stop services: docker-compose down"
echo "   Restart: docker-compose restart"
