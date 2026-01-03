#!/bin/bash

echo "🌱 Seeding LearnHub Database..."
echo ""

# Run the Go seed script
go run cmd/seed/main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Seeding completed successfully!"
else
    echo ""
    echo "❌ Seeding failed!"
    exit 1
fi
