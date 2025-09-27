# Fit-Up Server

A complete microservices backend implementation for the Fit-Up React Native application with authentication, messaging, and schema validation services.

## 🏗️ Architecture

```
React Native App → API Gateway → Microservices
      ↓               ↓              ↓
   Client Side → Port 8080 → Auth(8081), Message(8082), Schema(8083)
```

### Services Overview

| Service | Port | Description | Status |
|---------|------|-------------|---------|
| **API Gateway** | 8080 | Routes requests, handles CORS, load balancing | ✅ Complete |
| **Auth Service** | 8081 | User authentication, JWT tokens, OAuth2 | ✅ Complete |
| **Message Service** | 8082 | User messaging and conversations | ✅ Basic Implementation |
| **Schema Service** | 8083 | Data validation and schema management | ✅ Basic Implementation |
| **PostgreSQL** | 5432 | Primary database for all services | ✅ Complete |

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose**: For containerized deployment
- **Go 1.21+**: For local development
- **PostgreSQL 15+**: Database (included in Docker setup)

### Option 1: Docker (Recommended)

1. **Clone and navigate to server directory**:
   ```bash
   cd fit-up/server
   ```

2. **Start all services**:
   ```bash
   # Linux/Mac
   chmod +x start.sh
   ./start.sh
   
   # Windows
   start.bat
   ```

3. **Verify services are running**:
   ```bash
   curl http://localhost:8080/health  # API Gateway
   curl http://localhost:8081/health  # Auth Service
   curl http://localhost:8082/health  # Message Service
   curl http://localhost:8083/health  # Schema Service
   ```

### Option 2: Local Development

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL**:
   ```bash
   docker run --name fitup-postgres -e POSTGRES_DB=fitup -e POSTGRES_USER=fitup -e POSTGRES_PASSWORD=fitup_password -p 5432:5432 -d postgres:15-alpine
   ```

4. **Run services individually**:
   ```bash
   # Terminal 1 - Auth Service
   go run services/auth-service/cmd/main.go
   
   # Terminal 2 - Message Service
   go run services/message-service/cmd/main.go
   
   # Terminal 3 - Schema Service
   go run services/schema-service/cmd/main.go
   
   # Terminal 4 - API Gateway
   go run services/api-gateway/cmd/main.go
   ```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the server root:

```env
# Database
DATABASE_URL=postgres://fitup:fitup_password@localhost:5432/fitup?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXP=86400
REFRESH_TOKEN_EXP=2592000

# OAuth2 Configuration (Optional)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://localhost:19006

# Frontend URL
FRONTEND_URL=http://localhost:3000
```

### Service Configuration

Each service can be configured via environment variables:

- **PORT**: Service port (defaults: Gateway=8080, Auth=8081, Message=8082, Schema=8083)
- **DATABASE_URL**: PostgreSQL connection string
- **Service URLs**: For API Gateway routing

## 📚 API Documentation

### Authentication Endpoints

**Base URL**: `http://localhost:8080/auth`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/login` | User login with email/username and password | ❌ |
| POST | `/register` | User registration | ❌ |
| POST | `/logout` | User logout (revokes refresh tokens) | ✅ |
| POST | `/refresh-token` | Refresh access token | ❌ |
| POST | `/validate-token` | Validate JWT token | ❌ |
| POST | `/forgot-password` | Request password reset | ❌ |
| POST | `/reset-password` | Reset password with token | ❌ |
| POST | `/change-password` | Change password (authenticated users) | ✅ |
| GET | `/{username}` | Get user profile by username | ❌ |

### OAuth2 Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/oauth/{provider}` | Get OAuth authorization URL |
| GET | `/oauth/callback/{provider}` | Handle OAuth callback |
| POST | `/link/{provider}` | Link OAuth account to user |
| DELETE | `/unlink/{provider}` | Unlink OAuth account |
| GET | `/linked-accounts` | Get user's linked accounts |

**Supported Providers**: `google`, `github`, `facebook`

### Example API Calls

**Login**:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user@example.com",
    "password": "password123"
  }'
```

**Register**:
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "user@example.com",
    "password": "password123",
    "name": "New User"
  }'
```

**Authenticated Request**:
```bash
curl -X POST http://localhost:8080/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "current_password": "oldpass",
    "new_password": "newpass"
  }'
```

## 🔒 Security Features

### ✅ Implemented Security Measures

- **JWT Authentication**: Secure token-based authentication
- **Refresh Tokens**: Long-lived tokens for seamless user experience
- **Password Hashing**: bcrypt with proper salt rounds
- **Rate Limiting**: Protection against brute force attacks
- **CORS Configuration**: Controlled cross-origin requests
- **SQL Injection Protection**: Parameterized queries
- **Input Validation**: Request payload validation
- **Token Blacklisting**: Revoked token management
- **OAuth2 Integration**: Social login support

### 🛡️ Security Best Practices

- **Environment Variables**: Sensitive data not hardcoded
- **Database Connection**: Secure PostgreSQL configuration
- **Error Handling**: No sensitive information in error responses
- **Audit Logging**: Track authentication events
- **Health Checks**: Service monitoring and health endpoints

## 🧪 Testing

### Health Checks

```bash
# Check all services
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

### Authentication Flow Test

```bash
# 1. Register a user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","name":"Test User"}'

# 2. Login
LOGIN_RESPONSE=$(curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"test@example.com","password":"password123"}')

# 3. Extract token and test authenticated endpoint
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')
curl -X POST http://localhost:8080/auth/validate-token \
  -H "Authorization: Bearer $TOKEN"
```

## 🚨 Troubleshooting

### Common Issues

**1. Services not starting**:
```bash
# Check Docker status
docker ps

# View service logs
docker-compose logs auth-service
docker-compose logs api-gateway
```

**2. Database connection issues**:
```bash
# Check PostgreSQL is running
docker-compose logs postgres

# Test database connection
docker exec -it fit-up-server-postgres-1 psql -U fitup -d fitup -c "SELECT version();"
```

**3. CORS errors**:
- Update `CORS_ORIGINS` environment variable
- Ensure frontend URL is included in allowed origins

**4. JWT token issues**:
- Verify `JWT_SECRET` is set and consistent
- Check token expiration settings

### Useful Commands

```bash
# View all container logs
docker-compose logs -f

# Restart specific service
docker-compose restart auth-service

# Rebuild and restart
docker-compose up --build -d

# Stop all services
docker-compose down

# Clean up everything
docker-compose down -v --remove-orphans
```

## 📋 Development Status

### ✅ Completed Features

- **API Gateway**: Full proxy implementation with CORS
- **Auth Service**: Complete JWT authentication system
- **Database Layer**: PostgreSQL integration with connection pooling
- **Security**: Rate limiting, password hashing, token management
- **OAuth2**: Social login framework (needs provider configuration)
- **Docker Setup**: Fully containerized deployment
- **Health Monitoring**: Health check endpoints for all services

### 🚧 In Progress / Future Enhancements

- **Message Service**: Full messaging system implementation
- **Schema Service**: Advanced JSON schema validation
- **WebSocket Support**: Real-time messaging
- **File Upload**: User profile images and attachments
- **Email Service**: Account verification and notifications
- **Admin Dashboard**: Service monitoring and user management
- **API Documentation**: OpenAPI/Swagger documentation
- **Testing Suite**: Comprehensive unit and integration tests

### 🎯 Next Steps for Production

1. **SSL/TLS Configuration**: HTTPS termination
2. **Database Migrations**: Automated schema management
3. **Monitoring & Logging**: Prometheus, Grafana, ELK stack
4. **Load Balancing**: Multiple service instances
5. **Secret Management**: Vault or similar for sensitive data
6. **CI/CD Pipeline**: Automated testing and deployment
7. **Backup Strategy**: Database backup and recovery procedures

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/new-feature`
5. Submit a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Ready for development! 🚀**

The server implementation is complete and ready for React Native integration. All authentication flows are working, and the microservices architecture provides a solid foundation for scaling.