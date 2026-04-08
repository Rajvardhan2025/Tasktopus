#!/bin/bash

echo "🚀 Project Management Backend - Quick Start"
echo "=========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go is installed: $(go version)"
echo ""

# Check if MongoDB is running
if ! command -v mongosh &> /dev/null && ! command -v mongo &> /dev/null; then
    echo "⚠️  MongoDB CLI not found. Make sure MongoDB is running."
    echo "   You can start MongoDB with Docker:"
    echo "   docker run -d -p 27017:27017 --name mongodb mongo:latest"
    echo ""
fi

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env file created. Please update with your MongoDB URI if needed."
    echo ""
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
echo "✅ Dependencies installed"
echo ""

# Build the application
echo "🔨 Building application..."
go build -o bin/app main.go
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi
echo ""

# Run the application
echo "🎯 Starting server..."
echo "   Server will be available at http://localhost:8080"
echo "   Health check: http://localhost:8080/health"
echo "   API base: http://localhost:8080/api"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

./bin/app
