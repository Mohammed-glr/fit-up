@echo off
REM Fit-Up Server Startup Script (Windows)
REM This script starts the monolithic Fit-Up API server

echo ================================================================================
echo                        Fit-Up API Server Startup
echo ================================================================================
echo.

REM Check if .env file exists
if not exist .env (
    echo WARNING: .env file not found!
    echo Creating .env from .env.example if it exists...
    if exist .env.example (
        copy .env.example .env
        echo .env created! Please configure it before running.
        echo.
    ) else (
        echo ERROR: No .env.example found either.
        echo Please create a .env file with the following variables:
        echo   - DATABASE_URL
        echo   - JWT_SECRET
        echo   - PORT (optional, defaults to 8080)
        echo.
        pause
        exit /b 1
    )
)

REM Load environment from .env (basic version)
echo Loading environment variables...
for /F "usebackq tokens=*" %%A in (".env") do (
    set %%A
)

REM Check required environment variables
if "%DATABASE_URL%"=="" (
    echo ERROR: DATABASE_URL not set in .env file
    pause
    exit /b 1
)

if "%JWT_SECRET%"=="" (
    echo ERROR: JWT_SECRET not set in .env file
    pause
    exit /b 1
)

if "%PORT%"=="" (
    echo PORT not set, using default: 8080
    set PORT=8080
)

echo.
echo Configuration:
echo   PORT: %PORT%
echo   DATABASE: %DATABASE_URL%
echo.
echo ================================================================================
echo.

REM Check if binary exists, if not build it
if not exist bin\fitup-server.exe (
    echo Building server...
    go build -o bin\fitup-server.exe cmd\main.go
    if errorlevel 1 (
        echo.
        echo ERROR: Build failed!
        pause
        exit /b 1
    )
    echo Build successful!
    echo.
)

REM Start the server
echo Starting Fit-Up API Server...
echo.
bin\fitup-server.exe

pause
