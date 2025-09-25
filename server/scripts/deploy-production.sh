#!/bin/bash

# Lornian.com Production Deployment Script
# Run this script on your Azure VM (20.108.32.156)

set -e

echo "🚀 Deploying Lornian Backend to Production"
echo "Domain: lornian.com"
echo "IP: 20.108.32.156"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   print_error "This script should not be run as root"
   exit 1
fi

# Check if we're on the Azure VM
CURRENT_IP=$(curl -s ifconfig.me)
if [[ "$CURRENT_IP" != "20.108.32.156" ]]; then
    print_warning "Current IP ($CURRENT_IP) doesn't match expected Azure VM IP (20.108.32.156)"
    read -p "Continue anyway? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Step 1: Check prerequisites
echo "🔍 Checking prerequisites..."

# Check Docker
if ! command -v docker &> /dev/null; then
    print_error "Docker not found. Please install Docker first."
    echo "Run: curl -fsSL https://get.docker.com | sh"
    exit 1
fi

# Check Docker Compose
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose not found. Please install Docker Compose first."
    exit 1
fi

# Check if user is in docker group
if ! groups $USER | grep &>/dev/null '\bdocker\b'; then
    print_error "User is not in docker group."
    echo "Run: sudo usermod -aG docker $USER"
    echo "Then logout and login again."
    exit 1
fi

print_status "Prerequisites check passed"

# Step 2: Check DNS
echo "🌐 Checking DNS configuration..."
if getent hosts lornian.com | grep -q "20.108.32.156"; then
    print_status "Main domain DNS is correctly configured"
else
    print_warning "Main domain DNS might not be propagated yet or configured incorrectly"
    echo "Expected: lornian.com -> 20.108.32.156"
    echo "Current:"
    getent hosts lornian.com || echo "No A record found"
fi

if getent hosts api.lornian.com | grep -q "20.108.32.156"; then
    print_status "API subdomain DNS is correctly configured"
else
    print_warning "API subdomain DNS might not be propagated yet or configured incorrectly"
    echo "Expected: api.lornian.com -> 20.108.32.156"
    echo "Current:"
    getent hosts api.lornian.com || echo "No A record found"
fi

# Step 3: Check firewall
echo "🔥 Checking firewall configuration..."
if command -v ufw >/dev/null 2>&1; then
    if sudo ufw status | grep -q "80/tcp.*ALLOW"; then
        print_status "Port 80 is open"
    else
        print_warning "Port 80 might not be open"
        echo "Run: sudo ufw allow 80/tcp"
    fi

    if sudo ufw status | grep -q "443/tcp.*ALLOW"; then
        print_status "Port 443 is open"
    else
        print_warning "Port 443 might not be open"
        echo "Run: sudo ufw allow 443/tcp"
    fi
else
    print_warning "ufw not available - cannot check firewall status"
    echo "   Make sure ports 80 and 443 are open"
fi

# Step 4: Environment setup
echo "⚙️  Setting up environment..."

# Check if .env exists
if [[ ! -f .env ]]; then
    print_warning ".env file not found. Creating template..."
    cat > .env << EOF
# Database Configuration
DATABASE_URL=postgresql://username:password@host/database

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_minimum_32_characters

# OAuth Configuration (Optional)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
EOF
    chmod 600 .env
    print_error "Please edit .env file with your actual configuration"
    echo "nano .env"
    exit 1
else
    print_status ".env file found"
fi

# Step 5: SSL Setup
echo "🔒 Setting up SSL certificates..."

# Check if SSL certs exist
if [[ -f nginx/ssl/cert.pem && -f nginx/ssl/key.pem ]]; then
    print_status "SSL certificates found"
    
    # Check if they're about to expire (Let's Encrypt certs)
    if openssl x509 -in nginx/ssl/cert.pem -checkend 604800 -noout; then
        print_status "SSL certificates are valid"
    else
        print_warning "SSL certificates expire within 7 days"
    fi
else
    print_warning "SSL certificates not found. Setting up Let's Encrypt..."
    
    # Check if domain resolves to this server
    if getent hosts api.lornian.com | grep -q "20.108.32.156"; then
        echo "Setting up Let's Encrypt certificates..."
        if [[ -x ./scripts/setup-ssl.sh ]]; then
            ./scripts/setup-ssl.sh lornian.com admin@lornian.com
        else
            print_error "SSL setup script not found or not executable"
            exit 1
        fi
    else
        print_warning "DNS not pointing to this server. Creating self-signed certificates for testing..."
        if [[ -x ./scripts/setup-ssl.sh ]]; then
            ./scripts/setup-ssl.sh lornian.com
        else
            print_error "SSL setup script not found or not executable"
            exit 1
        fi
    fi
fi

# Step 6: Deploy application
echo "🚀 Deploying application..."

# Stop existing containers
if docker ps -q --filter "name=lornian-" | grep -q .; then
    print_status "Stopping existing containers..."
    docker-compose -f docker/docker-compose.prod.yml down
    
    # Wait for containers to fully stop and ports to be released
    echo "⏳ Waiting for containers to fully stop..."
    sleep 5
    
    # Ensure no containers are still using port 80
    if docker ps --filter "publish=80" --format '{{.Names}}' | grep -q .; then
        print_warning "Some containers are still using port 80. Forcing stop..."
        docker ps --filter "publish=80" --format '{{.Names}}' | xargs -r docker stop
        sleep 2
    fi
fi

# Build and start containers
print_status "Building and starting containers..."
docker-compose -f docker/docker-compose.prod.yml up --build -d

# Wait for containers to start
echo "⏳ Waiting for containers to start..."
sleep 10

# Step 7: Health checks
echo "🏥 Running health checks..."

# Check if containers are running
if docker ps --filter "name=lornian-nginx" --filter "status=running" | grep -q lornian-nginx; then
    print_status "Nginx container is running"
else
    print_error "Nginx container is not running"
    docker logs lornian-nginx
    exit 1
fi

if docker ps --filter "name=lornian-api-gateway" --filter "status=running" | grep -q lornian-api-gateway; then
    print_status "API Gateway container is running"
else
    print_error "API Gateway container is not running"
    docker logs lornian-api-gateway
    exit 1
fi

if docker ps --filter "name=lornian-auth-service" --filter "status=running" | grep -q lornian-auth-service; then
    print_status "Auth Service container is running"
else
    print_error "Auth Service container is not running"
    docker logs lornian-auth-service
    exit 1
fi

# Test HTTP to HTTPS redirect
echo "🔄 Testing HTTP to HTTPS redirect..."
if curl -s -I http://lornian.com | grep -q "301"; then
    print_status "HTTP to HTTPS redirect is working"
else
    print_warning "HTTP to HTTPS redirect might not be working"
fi

# Test HTTPS health endpoint
echo "🏥 Testing HTTPS health endpoint..."
if curl -s -k https://api.lornian.com/health | grep -q "healthy"; then
    print_status "HTTPS health endpoint is working"
else
    print_warning "HTTPS health endpoint might not be working"
    echo "Response:"
    curl -s -k https://api.lornian.com/health || echo "Failed to connect"
fi

# Step 8: Setup monitoring
echo "📊 Setting up monitoring..."

# Create log rotation for nginx
sudo tee /etc/logrotate.d/docker-nginx << EOF
/var/log/nginx/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 0644 nginx nginx
    postrotate
        docker kill -s USR1 lornian-nginx 2>/dev/null || true
    endscript
}
EOF

print_status "Log rotation configured"

# Setup certificate renewal cron job
if command -v crontab >/dev/null 2>&1; then
    if ! crontab -l 2>/dev/null | grep -q "certbot renew"; then
        (crontab -l 2>/dev/null; echo "0 2 * * * cd $(pwd) && docker run --rm -v \$(pwd)/nginx/ssl:/etc/letsencrypt -v certbot_webroot:/var/www/certbot certbot/certbot renew --quiet && docker exec lornian-nginx nginx -s reload") | crontab -
        print_status "Certificate renewal cron job added"
    fi
else
    print_warning "crontab not available - certificate auto-renewal not configured"
    echo "   You'll need to manually renew certificates"
fi

# Final status
echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📋 Summary:"
echo "  Main Domain: https://lornian.com (redirects to API)"
echo "  API Domain: https://api.lornian.com"
echo "  IP: 20.108.32.156"
echo "  API Health: https://api.lornian.com/health"
echo "  Nginx Health: https://api.lornian.com/nginx-health"
echo ""
echo "🔗 API Endpoints:"
echo "  Main API: https://api.lornian.com"
echo "  Auth: https://api.lornian.com/auth/*"
echo "  User: https://api.lornian.com/user/*"
echo "  AI: https://api.lornian.com/ai/*"
echo ""
echo "📊 Monitoring:"
echo "  Container status: docker ps"
echo "  Nginx logs: docker logs lornian-nginx"
echo "  App logs: docker logs lornian-api-gateway"
echo "  Access logs: docker exec lornian-nginx tail -f /var/log/nginx/access.log"
echo ""
echo "🔒 SSL certificate will auto-renew via cron job"
echo ""
print_status "Lornian.com is now live! 🚀"
