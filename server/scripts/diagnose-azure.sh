#!/bin/bash

# Azure VM Diagnostic Script for Lornian Backend

echo "🔍 Azure VM Diagnostic Report"
echo "============================="
echo ""

# System Information
echo "📊 System Information:"
echo "  OS: $(uname -a)"
echo "  User: $(whoami)"
echo "  Working Directory: $(pwd)"
echo ""

# Check required commands
echo "🛠️  Command Availability:"
commands=("docker" "docker-compose" "openssl" "getent" "netstat" "ss" "curl")
for cmd in "${commands[@]}"; do
    if command -v "$cmd" >/dev/null 2>&1; then
        echo "  ✅ $cmd: $(which $cmd)"
    else
        echo "  ❌ $cmd: Not found"
    fi
done
echo ""

# Docker status
echo "🐳 Docker Status:"
if command -v docker >/dev/null 2>&1; then
    echo "  Docker Version: $(docker --version)"
    echo "  Docker Service Status:"
    if docker info >/dev/null 2>&1; then
        echo "    ✅ Docker daemon is running"
        echo "  Running Containers:"
        docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || echo "    No containers running"
    else
        echo "    ❌ Docker daemon is not accessible"
    fi
else
    echo "  ❌ Docker not installed"
fi
echo ""

# Network and ports
echo "🌐 Network Status:"
echo "  External IP: $(curl -s ifconfig.me 2>/dev/null || echo "Unable to determine")"

if command -v netstat >/dev/null 2>&1; then
    echo "  Ports in use:"
    netstat -tlnp 2>/dev/null | grep -E ':(80|443|8080)' || echo "    No services on ports 80, 443, 8080"
elif command -v ss >/dev/null 2>&1; then
    echo "  Ports in use:"
    ss -tlnp 2>/dev/null | grep -E ':(80|443|8080)' || echo "    No services on ports 80, 443, 8080"
else
    echo "  ⚠️  Cannot check port status (netstat/ss not available)"
fi
echo ""

# DNS Resolution
echo "🏷️  DNS Resolution:"
domains=("lornian.com" "api.lornian.com")
for domain in "${domains[@]}"; do
    echo "  $domain:"
    if command -v getent >/dev/null 2>&1; then
        result=$(getent hosts "$domain" 2>/dev/null)
        if [ -n "$result" ]; then
            echo "    ✅ $result"
        else
            echo "    ❌ No DNS record found"
        fi
    else
        echo "    ⚠️  Cannot check DNS (getent not available)"
    fi
done
echo ""

# File Permissions
echo "📁 File Permissions:"
files=("scripts/setup-ssl.sh" "scripts/deploy-production.sh" "nginx/ssl")
for file in "${files[@]}"; do
    if [ -e "$file" ]; then
        echo "  ✅ $file: $(ls -ld "$file" | cut -d' ' -f1,3,4)"
    else
        echo "  ❌ $file: Not found"
    fi
done
echo ""

# SSL Certificates
echo "🔒 SSL Certificate Status:"
if [ -f "nginx/ssl/cert.pem" ] && [ -f "nginx/ssl/key.pem" ]; then
    echo "  ✅ SSL certificates exist"
    echo "  Certificate info:"
    openssl x509 -in nginx/ssl/cert.pem -text -noout | grep -E "(Subject:|DNS:|Not After)" | sed 's/^/    /'
else
    echo "  ❌ SSL certificates not found"
fi
echo ""

# Recent logs
echo "📝 Recent System Issues:"
if command -v journalctl >/dev/null 2>&1; then
    echo "  Recent Docker errors:"
    journalctl -u docker --since "1 hour ago" --no-pager -q | tail -5 | sed 's/^/    /' || echo "    No recent Docker errors"
else
    echo "  ⚠️  Cannot check system logs (journalctl not available)"
fi
echo ""

echo "🎯 Recommendations:"
echo "  1. If Docker is not running: sudo systemctl start docker"
echo "  2. If user not in docker group: sudo usermod -aG docker \$USER && logout"
echo "  3. If ports are busy: sudo systemctl stop nginx apache2 (if any)"
echo "  4. If DNS doesn't resolve: Check domain configuration"
echo "  5. For SSL issues: Run ./scripts/setup-ssl.sh lornian.com"
echo ""
echo "✅ Diagnostic complete!"
