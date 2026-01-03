#!/bin/bash

echo "🚀 Deploying LearnHub Backend to Production..."
echo ""

# Set production environment variables
export DATABASE_URL="postgresql://production_db_em4b_user:aMDjUPiChzv5mfw2O70NyICuqf0IvqWc@dpg-d5ca45shg0os73e4ruu0-a.oregon-postgres.render.com/production_db_em4b"
export PORT=8080
export ENVIRONMENT=production

echo "✓ Environment: $ENVIRONMENT"
echo "✓ Database: PostgreSQL (Production)"
echo "✓ Port: $PORT"
echo ""

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "❌ Failed to download dependencies"
    exit 1
fi
echo "✓ Dependencies downloaded"
echo ""

# Build the application
echo "🔨 Building application..."
go build -o learnhub-server main.go
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✓ Build successful"
echo ""

# Run the server
echo "🌐 Starting server..."
./learnhub-server
