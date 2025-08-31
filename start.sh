#!/bin/bash

# SecureWallet - Digital Banking Platform (Vulnerable) Startup Script
# This script helps you start the vulnerable application safely

echo "🚨 SecureWallet - Digital Banking Platform (Vulnerable) Startup Script"
echo ""

echo "⚠️  WARNING: This application is intentionally vulnerable!"
echo "   Only use in controlled, isolated environments."
echo "   Never deploy on public networks."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install it first."
    exit 1
fi

echo "✅ Docker and docker-compose are available"
echo ""

# Ask for confirmation
read -p "Do you want to start the vulnerable application? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Startup cancelled."
    exit 0
fi

echo ""
echo "🚀 Starting SecureWallet - Digital Banking Platform (Vulnerable)..."
echo ""

# Start the application
docker-compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo ""
    echo "✅ Application started successfully!"
    echo ""
    echo "🌐 Access Points:"
    echo "   • API Documentation: http://localhost:8000/docs"
    echo "   • Application: http://localhost:8000"
    echo "   • Frontend: http://localhost:3000"
    echo "   • API Info: http://localhost:8000/api/info"
    echo ""
    echo "🔑 Default Credentials:"
    echo "   • Admin: admin@vulnerable-app.com / admin123"
    echo ""
    echo "🔧 Vulnerability Levels:"
    echo "   • Basic: http://localhost:8000/api/vulnerabilities/info"
    echo "   • Change level in .env file: VULNERABILITY_LEVEL=basic|medium|hard|expert"
    echo ""
    echo "📚 Testing Examples:"
    echo "   • SQL Injection: curl 'http://localhost:8000/api/vulnerabilities/sql-injection/basic/user-search?username=admin'"
    echo "   • XSS: curl 'http://localhost:8000/api/vulnerabilities/xss/basic/reflected?user_input=<script>alert(1)</script>'"
else
    echo "❌ Failed to start application. Check docker-compose logs for details."
    exit 1
fi
