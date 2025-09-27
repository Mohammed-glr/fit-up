# FitUp Backend - Azure Deployment Guide

This guide walks you through deploying the FitUp backend services to an Azure Virtual Machine using Docker Compose.

## Prerequisites

1. **Azure VM** with Ubuntu 20.04 or later
2. **SSH access** to your Azure VM
3. **Domain name** pointed to your Azure VM's public IP
4. **OAuth credentials** for Google and GitHub authentication

## Quick Deployment

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd fit-up/server
```

### 2. Configure Environment

```bash
# Copy the environment template
cp docker/.env.azure.template docker/.env.azure

# Edit the environment file with your actual values
nano docker/.env.azure
```

### 3. Run Deployment Script

```bash
# Make the script executable
chmod +x scripts/deploy-azure.sh

# Run the deployment script
sudo ./scripts/deploy-azure.sh
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

```bash
# Create SSL directory
mkdir -p nginx/ssl

# Option 1: Use your own certificates
# Copy your certificates to nginx/ssl/cert.pem and nginx/ssl/key.pem

# Option 2: Use Let's Encrypt (recommended)
sudo apt install certbot
sudo certbot certonly --standalone -d your-domain.com
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/key.pem
```

### 4. Deploy Services

```bash
# Build and start services
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure build
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure up -d
```

## Configuration

### Environment Variables

The following environment variables must be configured in `docker/.env.azure`:

#### Required Variables
- `PUBLIC_HOST`: Your domain name (e.g., api.yourapp.com)
- `POSTGRES_PASSWORD`: Secure password for PostgreSQL
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Strong secret key (minimum 32 characters)
- `REDIS_PASSWORD`: Secure password for Redis

#### OAuth Configuration
- `GOOGLE_CLIENT_ID` & `GOOGLE_CLIENT_SECRET`: Google OAuth credentials
- `GITHUB_CLIENT_ID` & `GITHUB_CLIENT_SECRET`: GitHub OAuth credentials
- Redirect URIs should point to your domain

#### CORS Configuration
- `CORS_ORIGINS`: Comma-separated list of allowed origins

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
4. **API Gateway** - Main API entry point
5. **Auth Service** - Authentication and authorization
6. **Message Service** - Messaging functionality
7. **Schema Service** - Data schema management

## Monitoring and Maintenance

### View Service Status

```bash
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure ps
```

### View Logs

```bash
# All services
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure logs -f

# Specific service
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure logs -f nginx
```

### Update Services

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure build
docker-compose -f docker/docker-compose.azure.yml --env-file docker/.env.azure up -d
```

### Backup Data

```bash
# Backup PostgreSQL
docker exec fitup-postgres pg_dump -U fitup fitup > backup_$(date +%Y%m%d_%H%M%S).sql

# Backup Docker volumes
docker run --rm -v fitup_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup_$(date +%Y%m%d_%H%M%S).tar.gz -C /data .
```

## Security Considerations

1. **SSL/TLS**: Always use HTTPS in production
2. **Firewall**: Only expose ports 22 (SSH), 80 (HTTP), and 443 (HTTPS)
3. **Environment Variables**: Keep sensitive data in environment files, never in code
4. **Database**: Use strong passwords and consider using Azure Database for PostgreSQL
5. **Updates**: Regularly update the base images and system packages

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

1. **Use managed database**: Consider Azure Database for PostgreSQL
2. **Use managed Redis**: Consider Azure Cache for Redis
3. **Use Azure Container Instances**: For better scaling and management
4. **Set up monitoring**: Use Azure Monitor or similar
5. **Implement backup strategy**: Regular automated backups
6. **Use Azure Key Vault**: For managing secrets
7. **Set up CI/CD**: Automate deployments with Azure DevOps or GitHub Actions