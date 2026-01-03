#!/bin/bash

echo "Learning App Backend - Go REST API"
echo "==================================="
echo ""
echo "Make sure you have Go 1.21+ installed!"
echo ""
echo "To install dependencies:"
echo "  go mod download"
echo ""
echo "To run the server:"
echo "  go run main.go"
echo ""
echo "Server will start on: http://localhost:8080"
echo "Health check: http://localhost:8080/health"
echo ""
echo "Starting server..."
echo ""

go run main.go
