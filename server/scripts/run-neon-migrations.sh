#!/bin/bash

# Simple Migration Runner using psql
echo "🚀 Running Database Migrations to NeonDB using psql..."

# Set NeonDB connection string with endpoint ID
NEON_DATABASE_URL="postgresql://neondb_owner:npg_TlWxS1Gu9ihR@ep-falling-sunset-a2d4ksjo-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require&options=endpoint%3Dep-falling-sunset-a2d4ksjo"

echo "📍 Connecting to NeonDB..."

# Test connection first
echo "🔍 Testing database connection..."
psql "$NEON_DATABASE_URL" -c "SELECT version();" 

if [ $? -ne 0 ]; then
    echo "❌ Failed to connect to NeonDB. Please check your connection string."
    exit 1
fi

echo "✅ Connection successful!"

# Create schema_migrations table if it doesn't exist
echo "📊 Creating schema_migrations table..."
psql "$NEON_DATABASE_URL" -c "
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL DEFAULT false
);
"

# Check current migration status
echo "📋 Checking current migration status..."
CURRENT_VERSION=$(psql "$NEON_DATABASE_URL" -t -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" 2>/dev/null | xargs)

if [ -z "$CURRENT_VERSION" ]; then
    echo "📍 No migrations found. Starting from scratch."
    CURRENT_VERSION=0
else
    echo "📍 Current migration version: $CURRENT_VERSION"
fi

# Migration directory
MIGRATION_DIR="/workspaces/lornian-backend/shared/database/migrations"

# Run migration 001 if not already applied
if [ "$CURRENT_VERSION" -lt 1 ]; then
    echo "⬆️  Running migration 001_create_users.up.sql..."
    psql "$NEON_DATABASE_URL" -f "$MIGRATION_DIR/001_create_users.up.sql"
    
    if [ $? -eq 0 ]; then
        psql "$NEON_DATABASE_URL" -c "INSERT INTO schema_migrations (version) VALUES (1) ON CONFLICT (version) DO NOTHING;"
        echo "✅ Migration 001 completed successfully!"
    else
        echo "❌ Migration 001 failed!"
        exit 1
    fi
else
    echo "⏭️  Migration 001 already applied, skipping..."
fi

# Run migration 002 if not already applied
if [ "$CURRENT_VERSION" -lt 2 ]; then
    echo "⬆️  Running migration 002_add_jwt_management.up.sql..."
    psql "$NEON_DATABASE_URL" -f "$MIGRATION_DIR/002_add_jwt_management.up.sql"
    
    if [ $? -eq 0 ]; then
        psql "$NEON_DATABASE_URL" -c "INSERT INTO schema_migrations (version) VALUES (2) ON CONFLICT (version) DO NOTHING;"
        echo "✅ Migration 002 completed successfully!"
    else
        echo "❌ Migration 002 failed!"
        exit 1
    fi
else
    echo "⏭️  Migration 002 already applied, skipping..."
fi

# Show final migration status
echo "📊 Final migration status:"
psql "$NEON_DATABASE_URL" -c "SELECT * FROM schema_migrations ORDER BY version;"

# List all tables
echo "🗂️  Database tables:"
psql "$NEON_DATABASE_URL" -c "\dt"

# Show table details for key tables
echo "👥 Users table structure:"
psql "$NEON_DATABASE_URL" -c "\d users"

echo "🔑 JWT refresh tokens table structure:"
psql "$NEON_DATABASE_URL" -c "\d jwt_refresh_tokens"

echo "🚫 JWT blacklist table structure:"
psql "$NEON_DATABASE_URL" -c "\d jwt_blacklist"

echo "🎉 All migrations completed successfully on NeonDB!"
