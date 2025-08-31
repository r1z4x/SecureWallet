#!/bin/bash

echo "🚀 Initializing SecureWallet Full Application..."

# Stop any existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.dev.yml down

# Remove old volumes if needed
echo "🧹 Cleaning up old data..."
docker volume prune -f

# Start all services
echo "🐳 Starting all services..."
docker-compose -f docker-compose.dev.yml up -d

# Wait for MySQL to be ready
echo "⏳ Waiting for MySQL to be ready..."
sleep 30

# Check MySQL health
echo "🔍 Checking MySQL health..."
until docker-compose -f docker-compose.dev.yml exec mysql mysqladmin ping -h localhost -u root -prootpassword --silent; do
    echo "⏳ MySQL is not ready yet..."
    sleep 5
done
echo "✅ MySQL is ready!"

# Wait for backend to be ready
echo "⏳ Waiting for backend to be ready..."
sleep 20

# Check backend health
echo "🔍 Checking backend health..."
until curl -s http://localhost:8080/health > /dev/null; do
    echo "⏳ Backend is not ready yet..."
    sleep 5
done
echo "✅ Backend is ready!"

# Wait for frontend to be ready
echo "⏳ Waiting for frontend to be ready..."
sleep 15

# Check frontend health
echo "🔍 Checking frontend health..."
until curl -s http://localhost:3000 > /dev/null; do
    echo "⏳ Frontend is not ready yet..."
    sleep 5
done
echo "✅ Frontend is ready!"

echo ""
echo "🎉 SecureWallet Application is running!"
echo "🎨 Frontend: http://localhost:3000"
echo "🐍 Backend: http://localhost:8080"
echo "📚 API Docs: http://localhost:8080/docs"
echo "📊 Health Check: http://localhost:8080/health"
echo ""
echo "🔍 Debug: ./debug-backend.sh"
echo "📋 Create users: python create-users.py"
