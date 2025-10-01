package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdmdh/fit-up-server/shared/config"
)

// ConnectDB establishes a connection to PostgreSQL using pgxpool
func ConnectDB(ctx context.Context, databaseURL string, dbConfig config.DatabaseConfig) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Apply connection pool settings
	poolConfig.MaxConns = int32(dbConfig.MaxConnections)
	poolConfig.MinConns = int32(dbConfig.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(dbConfig.MaxConnLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = time.Duration(dbConfig.MaxConnIdleTime) * time.Minute
	poolConfig.HealthCheckPeriod = time.Duration(dbConfig.HealthCheckPeriod) * time.Minute
	poolConfig.ConnConfig.ConnectTimeout = time.Duration(dbConfig.ConnectTimeout) * time.Second

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("✅ Database connected successfully")
	log.Printf("   - Max connections: %d", dbConfig.MaxConnections)
	log.Printf("   - Min connections: %d", dbConfig.MinConnections)

	return pool, nil
}

// Close gracefully closes the database connection pool
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("Database connection closed")
	}
}
