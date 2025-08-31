#!/bin/bash

# SecureWallet Development Environment (Docker)
echo "🚀 Starting SecureWallet Development Environment (Docker)..."

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from env.example..."
    cp env.example .env
    echo "✅ .env file created. You can edit it if needed."
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Stop any existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.dev.yml down

# Start all services
echo "🐳 Starting all development services..."
docker-compose -f docker-compose.dev.yml up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 15

# Check backend status
echo "🔍 Checking backend status..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Backend is running at http://localhost:8080"
else
    echo "❌ Backend is not responding. Check logs with: docker-compose -f docker-compose.dev.yml logs backend"
fi

# Check frontend status
echo "🔍 Checking frontend status..."
if curl -s http://localhost:3000 > /dev/null; then
    echo "✅ Frontend is running at http://localhost:3000"
else
    echo "❌ Frontend is not responding. Check logs with: docker-compose -f docker-compose.dev.yml logs frontend"
fi

echo ""
echo "🎉 SecureWallet Development Environment is running!"
echo "🎨 Frontend: http://localhost:3000"
echo "🐍 Backend: http://localhost:8080"
echo "📚 API Docs: http://localhost:8080/docs"
echo ""
echo "📊 View logs:"
echo "  Frontend: docker-compose -f docker-compose.dev.yml logs -f frontend"
echo "  Backend:  docker-compose -f docker-compose.dev.yml logs -f backend"
echo ""
echo "🛑 Stop services: docker-compose -f docker-compose.dev.yml down"