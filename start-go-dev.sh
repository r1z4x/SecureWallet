#!/bin/bash

# Start SecureWallet Go Backend in Development Mode

echo "🚀 Starting SecureWallet Go Backend in Development Mode..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Check if air is installed for hot reload
if ! command -v air &> /dev/null; then
    echo "🔄 Installing Air for hot reload..."
    go install github.com/cosmtrek/air@v1.49.0
fi

echo "🐳 Starting Docker services..."
echo "📊 Services will be available at:"
echo "   - Backend API: http://localhost:8081"
echo "   - Frontend: http://localhost:3001"
echo "   - MySQL: localhost:3307"
echo "   - MongoDB: localhost:27018"
echo "   - Redis: localhost:6380"

# Start with Docker Compose
docker-compose -f docker-compose.dev.yml up -d --build
