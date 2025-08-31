#!/bin/bash

# Start SecureWallet Go Backend

echo "Starting SecureWallet Go Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file with default values..."
    cat > .env << EOF
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=securewallet
DB_PASSWORD=securewalletpass
DB_NAME=securewallet

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Application Configuration
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug

# JWT Configuration
JWT_SECRET=expert_secret_key_789
JWT_EXPIRE_MINUTES=525600

# CORS Configuration
CORS_ALLOW_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
EOF
fi

# Install dependencies
echo "Installing Go dependencies..."
go mod tidy

# Run the application
echo "Starting Go server on port 8080..."
go run main.go
