@echo off
REM Fit-Up Server Startup Script for Windows
echo 🚀 Starting Fit-Up Server Services...

REM Check if Docker is running
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker is not running. Please start Docker first.
    exit /b 1
)

REM Check for docker-compose
where docker-compose >nul 2>&1
if %errorlevel% equ 0 (
    set COMPOSE_CMD=docker-compose
) else (
    docker compose version >nul 2>&1
    if %errorlevel% equ 0 (
        set COMPOSE_CMD=docker compose
    ) else (
        echo ❌ docker-compose not found. Please install docker-compose.
        exit /b 1
    )
)

REM Set default environment variables if not set
if "%GOOGLE_CLIENT_ID%"=="" set GOOGLE_CLIENT_ID=your-google-client-id
if "%GOOGLE_CLIENT_SECRET%"=="" set GOOGLE_CLIENT_SECRET=your-google-client-secret
if "%GITHUB_CLIENT_ID%"=="" set GITHUB_CLIENT_ID=your-github-client-id
if "%GITHUB_CLIENT_SECRET%"=="" set GITHUB_CLIENT_SECRET=your-github-client-secret

echo 📦 Building and starting services...

REM Build and start all services
%COMPOSE_CMD% up --build -d

REM Wait for services to be ready
echo ⏳ Waiting for services to be ready...
timeout /t 10 /nobreak >nul

echo 🔍 Checking service health...

REM Check API Gateway
curl -f -s "http://localhost:8080/health" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ API Gateway is healthy
) else (
    echo ❌ API Gateway is not responding
)

REM Check Auth Service
curl -f -s "http://localhost:8081/health" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Auth Service is healthy
) else (
    echo ❌ Auth Service is not responding
)

REM Check Message Service
curl -f -s "http://localhost:8082/health" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Message Service is healthy
) else (
    echo ❌ Message Service is not responding
)

REM Check Schema Service
curl -f -s "http://localhost:8083/health" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Schema Service is healthy
) else (
    echo ❌ Schema Service is not responding
)

echo.
echo 🎉 Services are starting up!
echo.
echo 📍 Service URLs:
echo    API Gateway:     http://localhost:8080
echo    Auth Service:    http://localhost:8081
echo    Message Service: http://localhost:8082
echo    Schema Service:  http://localhost:8083
echo    PostgreSQL:      localhost:5432
echo.
echo 🔧 Useful commands:
echo    View logs:       %COMPOSE_CMD% logs -f [service-name]
echo    Stop services:   %COMPOSE_CMD% down
echo    Restart:         %COMPOSE_CMD% restart [service-name]
echo.
echo 📖 Test the API:
echo    Health check:    curl http://localhost:8080/health
echo    Auth endpoints:  curl http://localhost:8080/auth/login