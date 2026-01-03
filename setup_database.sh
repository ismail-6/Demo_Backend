#!/bin/bash

echo "🗄️  LearnHub Database Setup Script"
echo "=================================="
echo ""

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ Error: DATABASE_URL environment variable is not set"
    echo ""
    echo "Usage:"
    echo "  export DATABASE_URL='postgresql://user:password@host:port/database'"
    echo "  ./setup_database.sh"
    exit 1
fi

echo "📋 Database URL: ${DATABASE_URL%%@*}@***" # Hide credentials in output
echo ""

# Step 1: Create schema
echo "📐 Step 1: Creating database schema..."
psql "$DATABASE_URL" -f schema.sql

if [ $? -ne 0 ]; then
    echo "❌ Schema creation failed!"
    exit 1
fi
echo "✅ Schema created successfully"
echo ""

# Step 2: Seed data
echo "🌱 Step 2: Seeding database with initial data..."
psql "$DATABASE_URL" -f seed_with_relations.sql

if [ $? -ne 0 ]; then
    echo "❌ Data seeding failed!"
    exit 1
fi
echo "✅ Data seeded successfully"
echo ""

echo "=================================="
echo "✅ Database setup completed!"
echo ""
echo "Summary:"
echo "  - Schema created with all tables and relationships"
echo "  - 3 sample users added"
echo "  - 5 chapters added"
echo "  - 5 videos added (one per chapter)"
echo "  - 15 quiz questions added (3 per chapter)"
echo "  - Sample progress data added"
echo ""
echo "You can now start your application!"
