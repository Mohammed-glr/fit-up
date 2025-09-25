#!/bin/bash

# Migration Runner Script for NeonDB
echo "🚀 Running Database Migrations to NeonDB..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL environment variable is not set"
    echo "Please set it to your NeonDB connection string:"
    echo "export DATABASE_URL='postgresql://neondb_owner:your_password@ep-falling-sunset-a2d4ksjo-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require'"
    exit 1
fi

echo "📍 Using DATABASE_URL: ${DATABASE_URL:0:30}..."

# Install golang-migrate if not already installed
if ! command -v migrate &> /dev/null; then
    echo "📦 Installing golang-migrate..."
    
    # For Linux (GitHub Codespaces)
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
        sudo mv migrate /usr/local/bin/
        echo "✅ golang-migrate installed"
    else
        echo "❌ Please install golang-migrate manually for your OS"
        echo "Visit: https://github.com/golang-migrate/migrate/releases"
        exit 1
    fi
fi

# Set migration directory
MIGRATION_DIR="/workspaces/lornian-backend/shared/database/migrations"

if [ ! -d "$MIGRATION_DIR" ]; then
    echo "❌ Migration directory not found: $MIGRATION_DIR"
    exit 1
fi

echo "📁 Migration directory: $MIGRATION_DIR"

# List available migrations
echo "📋 Available migrations:"
ls -la $MIGRATION_DIR/*.sql

# Run migrations up
echo "⬆️  Running migrations up..."
migrate -path $MIGRATION_DIR -database "$DATABASE_URL" -verbose up

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed successfully!"
    
    # Show current migration version
    echo "📊 Current migration version:"
    migrate -path $MIGRATION_DIR -database "$DATABASE_URL" version
    
    # Verify tables were created
    echo "🔍 Verifying database tables..."
    psql "$DATABASE_URL" -c "\dt"
    
else
    echo "❌ Migration failed!"
    exit 1
fi

echo "🎉 Database migration to NeonDB complete!"
