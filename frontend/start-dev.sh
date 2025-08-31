#!/bin/bash

# Frontend Development Start Script
echo "🚀 Starting SecureWallet Frontend Development Server..."

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Start development server
echo "🔥 Starting Vite development server..."
npm run dev
