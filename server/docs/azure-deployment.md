# FitUp Server - Complete Azure VM Deployment Guide

This comprehensive guide walks you through deploying the FitUp backend server to an Azure Virtual Machine using Docker Compose and Nginx.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Azure VM Setup](#azure-vm-setup)
- [Quick Deployment](#quick-deployment)
- [Manual Deployment](#manual-deployment)
- [SSL/TLS Setup](#ssltls-setup)
- [Configuration](#configuration)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)

---

## Prerequisites

### Required
- **Azure Account** with active subscription
- **SSH Client** (Terminal, PuTTY, etc.)
- **Git** installed locally

### Optional but Recommended
- **Domain Name** pointed to your Azure VM's public IP (for production)
  - If you don't have a domain, you can use the VM's public IP address
  - See [Deployment Without Domain](#deployment-without-domain) section below

### Recommended Knowledge
- Basic Linux command line
- Docker and Docker Compose basics
- Nginx configuration understanding

## Deployment Without Domain

If you **don't have a domain name yet**, you can still deploy using your Azure VM's public IP address:

### Limitations
- ⚠️ **No HTTPS/SSL** - Let's Encrypt requires a domain name
- ⚠️ **Self-signed certificate only** - Browsers will show security warnings
- ⚠️ **Not recommended for production** - Use for testing/development only

### IP-Based Deployment Steps

1. **Get your VM's public IP**:
   ```bash
   az vm show -d -g fitup-rg -n fitup-server --query publicIps -o tsv
   # Example output: 20.121.45.123
   ```

2. **Configure environment with IP address**:
   ```bash
   # Use IP address instead of domain
   PUBLIC_HOST=20.121.45.123
   FRONTEND_URL=http://20.121.45.123  # or your frontend IP
   CORS_ORIGINS=http://20.121.45.123,http://localhost:19006
   
   # OAuth redirects will use IP
   GOOGLE_REDIRECT_URI=http://20.121.45.123/api/v1/auth/google/callback
   GITHUB_REDIRECT_URI=http://20.121.45.123/api/v1/auth/github/callback
   ```

3. **Use self-signed certificate or HTTP only**:
   - Option A: Generate self-signed certificate (browsers will warn)
   - Option B: Comment out HTTPS in nginx config (not secure)

4. **Access your API**:
   ```bash
   # HTTP access
   curl http://20.121.45.123/health
   
   # HTTPS with self-signed cert
   curl -k https://20.121.45.123/health
   ```

### Getting a Free Domain

If you need a domain, consider these free options:
- **FreeDNS** (freedns.afraid.org) - Free subdomain
- **DuckDNS** (duckdns.org) - Free dynamic DNS
- **No-IP** (noip.com) - Free hostname
- **Freenom** (.tk, .ml, .ga domains) - Free domains

Then point your domain's A record to your VM's public IP.

---

## Quick Deployment (With Domain)

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd fit-up/server
```

### 2. Configure Environment

```bash
# Copy the environment template
cp .env.production.template .env.production

# Edit the environment file with your actual values
nano .env.production
```

### 3. Run Deployment Script

```bash
# Make the script executable
chmod +x scripts/deploy-complete.sh

# Run the deployment script
sudo ./scripts/deploy-complete.sh
```

## Manual Deployment

If you prefer to deploy manually, follow these steps:

### 1. Install Docker and Docker Compose

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Setup Environment

```bash
# Copy environment template
cp docker/.env.azure.template docker/.env.azure

# Edit with your values
nano docker/.env.azure
```

### 3. Setup SSL Certificates

#### Option A: Let's Encrypt (Requires Domain)

```bash
# Install certbot
sudo apt install certbot

# Get certificate (requires domain pointing to this server)
sudo certbot certonly --standalone -d api.yourdomain.com
sudo cp /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/api.yourdomain.com/privkey.pem nginx/ssl/key.pem
```

#### Option B: Self-Signed Certificate (No Domain Required)

```bash
# Create SSL directory
mkdir -p nginx/ssl

# Generate self-signed certificate
# Replace YOUR_IP with your actual IP address or use localhost
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/key.pem \
  -out nginx/ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=FitUp/CN=YOUR_IP"

# Example for IP 20.121.45.123:
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/key.pem \
  -out nginx/ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=FitUp/CN=20.121.45.123"

# Set permissions
sudo chmod 644 nginx/ssl/cert.pem
sudo chmod 600 nginx/ssl/key.pem
```

⚠️ **Note**: Self-signed certificates will show security warnings in browsers. For production, always use a real domain with Let's Encrypt.

#### Option C: HTTP Only (Testing/Development)

If you want to skip HTTPS entirely for testing:

```bash
# Edit nginx config to disable HTTPS
nano nginx/nginx.production.conf

# Comment out the HTTPS server block (lines with 'listen 443 ssl')
# Keep only the HTTP server block on port 80
```

### 4. Deploy Services

```bash
# Build and start services
docker-compose -f docker-compose.production.yml --env-file .env.production build
docker-compose -f docker-compose.production.yml --env-file .env.production up -d
```

### 5. Verify Deployment

```bash
# Wait for services to start
sleep 30

# Test with domain
curl https://api.yourdomain.com/health

# OR test with IP address
curl http://YOUR_IP/health
curl -k https://YOUR_IP/health  # -k flag ignores SSL warnings

# Expected response:
# {"status":"healthy","service":"fit-up-api","timestamp":"..."}
```

## Configuration

### Environment Variables

The following environment variables must be configured in `.env.production`:

#### Required Variables
- `PUBLIC_HOST`: Your domain name (e.g., api.yourapp.com) **OR** your VM's public IP (e.g., 20.121.45.123)
- `POSTGRES_PASSWORD`: Secure password for PostgreSQL
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Strong secret key (minimum 32 characters)
- `REDIS_PASSWORD`: Secure password for Redis

#### OAuth Configuration
- `GOOGLE_CLIENT_ID` & `GOOGLE_CLIENT_SECRET`: Google OAuth credentials
- `GITHUB_CLIENT_ID` & `GITHUB_CLIENT_SECRET`: GitHub OAuth credentials
- Redirect URIs should point to your domain **OR** IP address
  - With domain: `https://api.yourdomain.com/api/v1/auth/google/callback`
  - With IP: `http://YOUR_IP/api/v1/auth/google/callback`

#### CORS Configuration
- `CORS_ORIGINS`: Comma-separated list of allowed origins
  - With domain: `https://yourdomain.com,https://app.yourdomain.com`
  - With IP: `http://YOUR_IP,http://localhost:19006`

#### Example Configuration (With IP Address)

```bash
PUBLIC_HOST=20.121.45.123
FRONTEND_URL=http://20.121.45.123
POSTGRES_PASSWORD=your-secure-password
DATABASE_URL=postgres://fitup:your-secure-password@postgres:5432/fitup?sslmode=disable
JWT_SECRET=your-32-character-secret-key-here
REDIS_PASSWORD=your-redis-password
GOOGLE_REDIRECT_URI=http://20.121.45.123/api/v1/auth/google/callback
GITHUB_REDIRECT_URI=http://20.121.45.123/api/v1/auth/github/callback
CORS_ORIGINS=http://20.121.45.123,http://localhost:19006
```

#### Example Configuration (With Domain)

```bash
PUBLIC_HOST=api.yourdomain.com
FRONTEND_URL=https://yourdomain.com
POSTGRES_PASSWORD=your-secure-password
DATABASE_URL=postgres://fitup:your-secure-password@postgres:5432/fitup?sslmode=disable
JWT_SECRET=your-32-character-secret-key-here
REDIS_PASSWORD=your-redis-password
GOOGLE_REDIRECT_URI=https://api.yourdomain.com/api/v1/auth/google/callback
GITHUB_REDIRECT_URI=https://api.yourdomain.com/api/v1/auth/github/callback
CORS_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

### Network Configuration

The deployment creates a custom Docker network (`fitup-network`) with subnet `172.20.0.0/16`.

### Port Configuration

- **Port 80**: HTTP (redirects to HTTPS)
- **Port 443**: HTTPS
- Internal services communicate through the Docker network

## Services

The deployment includes the following services:

1. **Nginx** - Reverse proxy and SSL termination
2. **PostgreSQL** - Primary database
3. **Redis** - Caching and session storage
4. **FitUp API Server** - Monolithic Go application serving:
   - Authentication & OAuth
   - Workout & Exercise Management
   - Real-time Messaging (WebSocket)
   - Food Tracking & Nutrition
   - Schema Validation

## Monitoring and Maintenance

### View Service Status

```bash
docker-compose -f docker-compose.production.yml ps
```

### View Logs

```bash
# All services
docker-compose -f docker-compose.production.yml logs -f

# Specific service
docker-compose -f docker-compose.production.yml logs -f api-server
docker-compose -f docker-compose.production.yml logs -f nginx
docker-compose -f docker-compose.production.yml logs -f postgres
```

### Update Services

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose -f docker-compose.production.yml --env-file .env.production build
docker-compose -f docker-compose.production.yml --env-file .env.production up -d
```

### Backup Data

```bash
# Backup PostgreSQL
docker exec fitup-postgres pg_dump -U fitup fitup > backup_$(date +%Y%m%d_%H%M%S).sql

# Backup Docker volumes
docker run --rm -v fitup_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup_$(date +%Y%m%d_%H%M%S).tar.gz -C /data .
```

## Security Considerations

1. **SSL/TLS**: Always use HTTPS with real certificates in production
   - ⚠️ Self-signed certificates are only for testing
   - ⚠️ HTTP-only deployments are insecure for production
   - ✅ Use Let's Encrypt with a real domain for production

2. **Firewall**: Only expose ports 22 (SSH), 80 (HTTP), and 443 (HTTPS)
   ```bash
   sudo ufw status
   sudo ufw allow 22/tcp
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```

3. **Environment Variables**: Keep sensitive data in environment files, never in code
   - Never commit `.env.production` to git
   - Use strong, random passwords (32+ characters)
   - Generate JWT secret: `openssl rand -base64 32`

4. **Database**: Use strong passwords and consider using Azure Database for PostgreSQL for production

5. **Updates**: Regularly update the base images and system packages
   ```bash
   sudo apt update && sudo apt upgrade -y
   docker-compose pull
   docker-compose up -d --build
   ```

6. **OAuth Security**:
   - With domain: Use HTTPS redirect URIs only
   - With IP: OAuth providers may reject IP-based redirects
   - Configure authorized domains in OAuth provider settings

## Troubleshooting

### Service Won't Start

1. Check logs for the specific service
2. Verify environment variables are set correctly
3. Ensure all dependencies are healthy

### Database Connection Issues

1. Verify DATABASE_URL is correct
2. Check PostgreSQL container is running and healthy
3. Ensure network connectivity between services

### SSL Certificate Issues

1. Verify certificates are in the correct location (`nginx/ssl/`)
2. Check certificate file permissions
3. Ensure certificates are valid and not expired

### Performance Issues

1. Monitor resource usage with `docker stats`
2. Check log files for errors
3. Consider scaling services horizontally

## Support

For issues and questions:
1. Check the logs first
2. Review the environment configuration
3. Ensure all prerequisites are met
4. Submit an issue with relevant logs and configuration (without secrets)

## Production Recommendations

1. **Get a Domain Name**: Required for:
   - Let's Encrypt SSL certificates
   - OAuth provider approval
   - Professional appearance
   - SEO and branding
   
2. **Use managed database**: Consider Azure Database for PostgreSQL
   
3. **Use managed Redis**: Consider Azure Cache for Redis
   
4. **Use Azure Container Instances**: For better scaling and management
   
5. **Set up monitoring**: Use Azure Monitor or similar
   
6. **Implement backup strategy**: Regular automated backups
   
7. **Use Azure Key Vault**: For managing secrets
   
8. **Set up CI/CD**: Automate deployments with Azure DevOps or GitHub Actions

## Temporary/Testing Deployment (Without Domain)

If you're just testing and don't have a domain yet:

### Quick HTTP Deployment

1. **Use IP address in environment**:
   ```bash
   PUBLIC_HOST=YOUR_VM_IP
   FRONTEND_URL=http://YOUR_VM_IP
   CORS_ORIGINS=http://YOUR_VM_IP,http://localhost:19006
   ```

2. **Skip SSL or use self-signed certificate**

3. **Test API**:
   ```bash
   curl http://YOUR_VM_IP/health
   ```

4. **Mobile app configuration**:
   ```typescript
   // In your React Native app
   const API_URL = 'http://YOUR_VM_IP/api/v1';
   ```

### Important Notes for IP-Based Deployment

- ⚠️ **OAuth limitations**: Many OAuth providers (Google, GitHub) may not accept IP-based redirect URIs
- ⚠️ **No SSL**: Communication is unencrypted (not secure)
- ⚠️ **Browser warnings**: Self-signed certificates show security warnings
- ⚠️ **Not production-ready**: Use domains for production deployments
- ✅ **Good for**: Development, testing, proof-of-concept

### Migration Path

When you get a domain:

1. **Update environment variables** with domain name
2. **Get Let's Encrypt certificate**
3. **Update OAuth redirect URIs**
4. **Update CORS origins**
5. **Redeploy**:
   ```bash
   docker-compose -f docker-compose.production.yml down
   docker-compose -f docker-compose.production.yml --env-file .env.production up -d
   ```