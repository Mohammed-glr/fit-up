#!/bin/bash

# Simplified Azure VM Deployment Script for Lornian Backend
# Use this if the main deployment script fails

set -e

echo "🚀 Simple Deployment for Azure VM"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Step 1: Basic checks
echo "🔍 Basic checks..."
if ! command -v docker &> /dev/null; then
    print_error "Docker not found. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose not found. Please install Docker Compose first."
    exit 1
fi

if ! docker info >/dev/null 2>&1; then
    print_error "Docker daemon not running or not accessible."
    echo "Try: sudo systemctl start docker"
    echo "Or add user to docker group: sudo usermod -aG docker \$USER"
    exit 1
fi

print_status "Docker is available and running"

# Step 2: Environment check
echo "⚙️  Environment setup..."
if [[ ! -f .env ]]; then
    print_warning ".env file not found. Creating minimal template..."
    cat > .env << EOF
# Minimal configuration for Lornian Backend
DATABASE_URL=postgresql://user:pass@localhost/lornian
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production
EOF
    chmod 600 .env
    print_status ".env file created"
else
    print_status ".env file found"
fi

# Step 3: SSL certificates
echo "🔒 SSL setup..."
if [[ ! -f nginx/ssl/cert.pem || ! -f nginx/ssl/key.pem ]]; then
    print_warning "SSL certificates not found. Creating self-signed certificates..."
    
    # Create SSL directory
    mkdir -p nginx/ssl
    
    # Generate self-signed certificates directly without Docker
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=US/ST=State/L=City/O=Lornian/CN=lornian.com" \
        -config <(
            echo '[distinguished_name]'
            echo '[req]'
            echo 'distinguished_name = distinguished_name'
            echo '[v3_req]'
            echo 'keyUsage = keyEncipherment, dataEncipherment'
            echo 'extendedKeyUsage = serverAuth'
            echo "subjectAltName = @alt_names"
            echo '[alt_names]'
            echo "DNS.1 = lornian.com"
            echo "DNS.2 = www.lornian.com"
            echo "DNS.3 = api.lornian.com"
            echo "DNS.4 = localhost"
        ) -extensions v3_req
    
    chmod 600 nginx/ssl/key.pem
    chmod 644 nginx/ssl/cert.pem
    
    print_status "Self-signed SSL certificates created"
else
    print_status "SSL certificates found"
fi

# Step 4: Stop existing containers
echo "🛑 Stopping existing containers..."
if docker ps -q --filter "name=lornian-" | grep -q .; then
    print_status "Stopping existing containers..."
    docker-compose -f docker-compose.nginx.yml down 2>/dev/null || true
    sleep 3
else
    print_status "No existing containers found"
fi

# Step 5: Build and deploy
echo "🚀 Building and deploying..."
print_status "Starting build process..."
if docker-compose -f docker-compose.nginx.yml up --build -d; then
    print_status "Containers started successfully"
else
    print_error "Failed to start containers"
    echo "Check logs with: docker-compose -f docker-compose.nginx.yml logs"
    exit 1
fi

# Step 6: Wait and verify
echo "⏳ Waiting for services to start..."
sleep 15

echo "🏥 Health checks..."
failed=0

# Check if containers are running
containers=("lornian-nginx" "lornian-api-gateway" "lornian-auth-service")
for container in "${containers[@]}"; do
    if docker ps --filter "name=$container" --filter "status=running" | grep -q "$container"; then
        print_status "$container is running"
    else
        print_error "$container is not running"
        failed=1
    fi
done

if [ $failed -eq 0 ]; then
    echo ""
    echo "🎉 Deployment completed successfully!"
    echo ""
    echo "📋 Service URLs:"
    echo "  Main site: https://localhost (or your domain)"
    echo "  API: https://localhost/health"
    echo "  Container status: docker ps"
    echo ""
    echo "🔍 Troubleshooting:"
    echo "  View logs: docker-compose -f docker-compose.nginx.yml logs"
    echo "  Restart: docker-compose -f docker-compose.nginx.yml restart"
    echo "  Stop: docker-compose -f docker-compose.nginx.yml down"
    echo ""
    print_status "Lornian Backend is now running! 🚀"
else
    echo ""
    print_error "Some services failed to start"
    echo "Check the logs for details:"
    echo "  docker-compose -f docker-compose.nginx.yml logs"
    exit 1
fi
