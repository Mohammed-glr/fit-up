#!/bin/bash

# SSL Certificate Setup Script for Lornian Backend

set -e

DOMAIN=${1:-"your-domain.com"}
EMAIL=${2:-"your-email@domain.com"}
SSL_DIR="./nginx/ssl"

# Cleanup function
cleanup() {
    echo "Cleaning up temporary nginx container..."
    docker stop temp-nginx >/dev/null 2>&1 || true
    docker rm temp-nginx >/dev/null 2>&1 || true
    rm -f "/tmp/temp-nginx-${DOMAIN}.conf" >/dev/null 2>&1 || true
    rm -rf "/tmp/letsencrypt-${DOMAIN}" >/dev/null 2>&1 || true
}

# Set trap to cleanup on exit
trap cleanup EXIT

echo "Setting up SSL certificates for domain: $DOMAIN"

# Create SSL directory if it doesn't exist
mkdir -p "$SSL_DIR"

# Check if we're in production or development
# If no email provided, create self-signed certificates (even for production domains)
if [[ "$DOMAIN" == "localhost" || "$DOMAIN" == "127.0.0.1" || "$DOMAIN" == "your-domain.com" || -z "$EMAIL" || "$EMAIL" == "your-email@domain.com" ]]; then
    echo "Creating self-signed certificate..."
    
    # Generate self-signed certificate for development
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout "$SSL_DIR/key.pem" \
        -out "$SSL_DIR/cert.pem" \
        -subj "/C=US/ST=State/L=City/O=Development/CN=$DOMAIN" \
        -config <(
            echo '[distinguished_name]'
            echo '[req]'
            echo 'distinguished_name = distinguished_name'
            echo '[v3_req]'
            echo 'keyUsage = keyEncipherment, dataEncipherment'
            echo 'extendedKeyUsage = serverAuth'
            echo "subjectAltName = @alt_names"
            echo '[alt_names]'
            echo "DNS.1 = $DOMAIN"
            echo "DNS.2 = www.$DOMAIN"
            echo "DNS.3 = api.$DOMAIN"
        ) -extensions v3_req
    
    echo "✅ Self-signed certificate created"
    echo "⚠️  Warning: This certificate is not trusted by browsers"
    echo "   For production with proper DNS, provide an email address"
    
else
    echo "Production mode: Setting up Let's Encrypt certificate..."
    
    # Create temporary nginx config for Let's Encrypt challenge
    TEMP_CONFIG_FILE="/tmp/temp-nginx-${DOMAIN}.conf"
    cat > "$TEMP_CONFIG_FILE" << EOF
events {
    worker_connections 1024;
}

http {
    server {
        listen 80;
        server_name $DOMAIN;

        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }

        location / {
            return 301 https://\$server_name\$request_uri;
        }
    }
}
EOF

    echo "Starting temporary nginx for Let's Encrypt challenge..."
    
    # Clean up any existing temp-nginx container
    if docker ps -a --format '{{.Names}}' | grep -q '^temp-nginx$'; then
        echo "Removing existing temp-nginx container..."
        docker stop temp-nginx >/dev/null 2>&1 || true
        docker rm temp-nginx >/dev/null 2>&1 || true
    fi
    
    # Check if port 80 is available
    PORT_CHECK_FAILED=false
    if command -v netstat >/dev/null 2>&1; then
        if netstat -tlnp 2>/dev/null | grep ':80 ' >/dev/null; then
            echo "⚠️  Port 80 appears to be in use"
            PORT_CHECK_FAILED=true
        fi
    elif command -v ss >/dev/null 2>&1; then
        if ss -tlnp 2>/dev/null | grep ':80 ' >/dev/null; then
            echo "⚠️  Port 80 appears to be in use"
            PORT_CHECK_FAILED=true
        fi
    else
        echo "⚠️  Cannot check port availability (netstat/ss not found)"
    fi
    
    if [ "$PORT_CHECK_FAILED" = true ]; then
        echo "⚠️  Port 80 might be in use. Attempting to start anyway..."
        echo "   If this fails, try stopping any services using port 80"
    fi
    
    # Wait a moment for any port to be fully released
    sleep 2
    
    # Start nginx with temporary config
    echo "Starting temporary nginx for Let's Encrypt challenge..."
    if ! docker run -d --name temp-nginx \
        -p 80:80 \
        -v "$TEMP_CONFIG_FILE:/etc/nginx/nginx.conf:ro" \
        -v certbot_webroot:/var/www/certbot \
        nginx:alpine; then
        
        echo "❌ Failed to start temporary nginx container."
        echo "   This might be due to port 80 being in use or Docker issues."
        echo "   Falling back to self-signed certificate..."
        
        # Generate self-signed certificate as fallback
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout "$SSL_DIR/key.pem" \
            -out "$SSL_DIR/cert.pem" \
            -subj "/C=US/ST=State/L=City/O=Lornian/CN=$DOMAIN" \
            -config <(
                echo '[distinguished_name]'
                echo '[req]'
                echo 'distinguished_name = distinguished_name'
                echo '[v3_req]'
                echo 'keyUsage = keyEncipherment, dataEncipherment'
                echo 'extendedKeyUsage = serverAuth'
                echo "subjectAltName = @alt_names"
                echo '[alt_names]'
                echo "DNS.1 = $DOMAIN"
                echo "DNS.2 = www.$DOMAIN"
                echo "DNS.3 = api.$DOMAIN"
            ) -extensions v3_req
        
        echo "✅ Self-signed certificate created as fallback"
        echo "⚠️  Warning: This certificate is not trusted by browsers"
        
        # Skip the Let's Encrypt process
        chmod 600 "$SSL_DIR/key.pem"
        chmod 644 "$SSL_DIR/cert.pem"
        echo "🔒 SSL certificates are ready!"
        exit 0
    fi

    echo "Waiting for nginx to start..."
    sleep 5

    # Get Let's Encrypt certificate
    echo "Attempting to get Let's Encrypt certificate..."
    
    # Use a temporary directory outside the build context for Let's Encrypt
    LETSENCRYPT_DIR="/tmp/letsencrypt-${DOMAIN}"
    mkdir -p "$LETSENCRYPT_DIR"
    
    if docker run --rm \
        -v "$LETSENCRYPT_DIR:/etc/letsencrypt" \
        -v certbot_webroot:/var/www/certbot \
        certbot/certbot \
        certonly --webroot \
        --webroot-path=/var/www/certbot \
        --email "$EMAIL" \
        --agree-tos \
        --no-eff-email \
        -d "$DOMAIN" \
        -d "www.$DOMAIN" \
        -d "api.$DOMAIN" \
        --non-interactive; then
        
        # Success - copy certificates to expected location
        cp "$LETSENCRYPT_DIR/live/$DOMAIN/fullchain.pem" "$SSL_DIR/cert.pem"
        cp "$LETSENCRYPT_DIR/live/$DOMAIN/privkey.pem" "$SSL_DIR/key.pem"
        
        # Clean up temporary Let's Encrypt directory
        rm -rf "$LETSENCRYPT_DIR"
        
        echo "✅ Let's Encrypt certificate obtained successfully"
    else
        echo "❌ Let's Encrypt failed. Creating self-signed certificate as fallback..."
        
        # Clean up temporary Let's Encrypt directory
        rm -rf "$LETSENCRYPT_DIR" >/dev/null 2>&1 || true
        
        # Generate self-signed certificate as fallback
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout "$SSL_DIR/key.pem" \
            -out "$SSL_DIR/cert.pem" \
            -subj "/C=US/ST=State/L=City/O=Lornian/CN=$DOMAIN" \
            -config <(
                echo '[distinguished_name]'
                echo '[req]'
                echo 'distinguished_name = distinguished_name'
                echo '[v3_req]'
                echo 'keyUsage = keyEncipherment, dataEncipherment'
                echo 'extendedKeyUsage = serverAuth'
                echo "subjectAltName = @alt_names"
                echo '[alt_names]'
                echo "DNS.1 = $DOMAIN"
                echo "DNS.2 = www.$DOMAIN"
                echo "DNS.3 = api.$DOMAIN"
            ) -extensions v3_req
        
        echo "✅ Self-signed certificate created as fallback"
        echo "⚠️  Warning: This certificate is not trusted by browsers"
    fi

    # Stop temporary nginx
    echo "Cleaning up temporary nginx container..."
    docker stop temp-nginx >/dev/null 2>&1 || true
    docker rm temp-nginx >/dev/null 2>&1 || true
fi

# Set proper permissions
chmod 600 "$SSL_DIR/key.pem"
chmod 644 "$SSL_DIR/cert.pem"

echo "🔒 SSL certificates are ready!"
echo "Certificate: $SSL_DIR/cert.pem"
echo "Private Key: $SSL_DIR/key.pem"
echo ""
echo "Next steps:"
echo "1. Update your domain in nginx/nginx.conf"
echo "2. Run: docker-compose -f docker-compose.nginx.yml up -d"
