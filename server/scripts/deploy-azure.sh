#!/bin/bash

# Azure VM Deployment Script
# Run this script on your Azure VM to deploy the Lornian Backend

set -e

echo "🚀 Starting Azure VM deployment for Lornian Backend..."

# Update system
echo "📦 Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install Docker
if ! command -v docker &> /dev/null; then
    echo "🐳 Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
fi

# Install Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "🐳 Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create application directory
APP_DIR="/opt/lornian-backend"
echo "📁 Creating application directory at $APP_DIR..."
sudo mkdir -p $APP_DIR
sudo chown $USER:$USER $APP_DIR

# Clone or update repository (you'll need to modify this for your repo)
if [ ! -d "$APP_DIR/.git" ]; then
    echo "📥 Cloning repository..."
    # Replace with your actual repository URL
    # git clone https://github.com/yourusername/lornian-backend.git $APP_DIR
    echo "⚠️  Please clone your repository to $APP_DIR manually"
else
    echo "🔄 Updating repository..."
    cd $APP_DIR
    git pull origin main
fi

cd $APP_DIR

# Copy environment file
if [ ! -f ".env" ]; then
    echo "📝 Creating environment file..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your actual configuration values"
    echo "📝 Don't forget to update:"
    echo "   - DATABASE_URL (Azure PostgreSQL)"
    echo "   - JWT_SECRET"
    echo "   - API keys"
    echo "   - Domain names"
fi

# Create SSL directory
sudo mkdir -p nginx/ssl
echo "🔐 SSL certificates directory created at nginx/ssl/"
echo "📝 Please add your SSL certificates (cert.pem and key.pem) to nginx/ssl/"

# Build and start services
echo "🏗️  Building and starting services..."
docker-compose -f docker/docker-compose.azure.yml --env-file .env build
docker-compose -f docker/docker-compose.azure.yml --env-file .env up -d

# Setup logrotate for Docker logs
echo "📋 Setting up log rotation..."
sudo tee /etc/logrotate.d/docker-containers > /dev/null <<EOF
/var/lib/docker/containers/*/*.log {
  rotate 7
  daily
  compress
  size=1M
  missingok
  delaycompress
  copytruncate
}
EOF

# Create systemd service for auto-restart
echo "🔄 Creating systemd service..."
sudo tee /etc/systemd/system/lornian-backend.service > /dev/null <<EOF
[Unit]
Description=Lornian Backend Services
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$APP_DIR
ExecStart=/usr/local/bin/docker-compose -f docker/docker-compose.azure.yml --env-file .env up -d
ExecStop=/usr/local/bin/docker-compose -f docker/docker-compose.azure.yml --env-file .env down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable lornian-backend.service
sudo systemctl daemon-reload

# Setup firewall
echo "🔥 Configuring firewall..."
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw --force enable

echo "✅ Deployment script completed!"
echo ""
echo "📋 Next steps:"
echo "1. Edit .env file with your actual configuration"
echo "2. Add SSL certificates to nginx/ssl/ directory"
echo "3. Update nginx/nginx.conf with your domain name"
echo "4. Run: sudo systemctl start lornian-backend"
echo "5. Check status: docker-compose -f docker/docker-compose.azure.yml ps"
echo ""
echo "🔍 Useful commands:"
echo "   View logs: docker-compose -f docker/docker-compose.azure.yml --env-file .env logs -f"
echo "   Restart: sudo systemctl restart lornian-backend"
echo "   Stop: sudo systemctl stop lornian-backend"
echo "   Redeploy: git pull && docker-compose -f docker/docker-compose.azure.yml --env-file .env build && docker-compose -f docker/docker-compose.azure.yml --env-file .env up -d"
