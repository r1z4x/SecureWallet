#!/bin/bash

# Start SecureWallet Go Backend

echo "Starting SecureWallet Go Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Install dependencies
echo "Installing Go dependencies..."
go mod tidy

# Run the application
echo "Starting Go server on port 8080..."
go run main.go
