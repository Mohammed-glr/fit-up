package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: Step 1 - Define database connection utilities:
//   - NewConnection(config DatabaseConfig) (*sql.DB, error) - Create connection pool

type Postgres struct {
	Pool *pgxpool.Pool
}

func ConnectDB(ctx context.Context, dsn string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Connection pool configuration with production-ready settings
	config.MaxConns = 30                       // Maximum number of connections in the pool
	config.MinConns = 5                        // Minimum number of connections to maintain
	config.MaxConnLifetime = time.Hour         // Maximum lifetime of a connection
	config.MaxConnIdleTime = 30 * time.Minute  // Maximum time a connection can be idle
	config.HealthCheckPeriod = 5 * time.Minute // How often to check connection health

	// Connection timeout settings
	config.ConnConfig.ConnectTimeout = 10 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{
		Pool: pool,
	}, nil
}

// Close gracefully closes the connection pool
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// Ping checks if the database is reachable
func (p *Postgres) Ping(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

// Stats returns connection pool statistics
func (p *Postgres) Stats() *pgxpool.Stat {
	return p.Pool.Stat()
}

// Health returns detailed health information about the connection pool
func (p *Postgres) Health(ctx context.Context) map[string]interface{} {
	stats := p.Pool.Stat()

	health := map[string]interface{}{
		"total_connections":        stats.TotalConns(),
		"acquired_connections":     stats.AcquiredConns(),
		"idle_connections":         stats.IdleConns(),
		"constructing_connections": stats.ConstructingConns(),
		"max_connections":          stats.MaxConns(),
		"acquire_count":            stats.AcquireCount(),
		"acquire_duration":         stats.AcquireDuration().String(),
		"canceled_acquire_count":   stats.CanceledAcquireCount(),
		"empty_acquire_count":      stats.EmptyAcquireCount(),
		"new_connections_count":    stats.NewConnsCount(),
	}

	// Test connectivity
	if err := p.Pool.Ping(ctx); err != nil {
		health["ping_error"] = err.Error()
		health["status"] = "unhealthy"
	} else {
		health["status"] = "healthy"
	}

	return health
}

// DatabaseConfig represents database connection configuration
type DatabaseConfig struct {
	MaxConnections    int
	MinConnections    int
	MaxConnLifetime   int // in minutes
	MaxConnIdleTime   int // in minutes
	HealthCheckPeriod int // in minutes
	ConnectTimeout    int // in seconds
}

// ConnectDBWithConfig creates a new database connection with custom configuration
func ConnectDBWithConfig(ctx context.Context, dsn string, dbConfig DatabaseConfig) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Apply custom configuration
	config.MaxConns = int32(dbConfig.MaxConnections)
	config.MinConns = int32(dbConfig.MinConnections)
	config.MaxConnLifetime = time.Duration(dbConfig.MaxConnLifetime) * time.Minute
	config.MaxConnIdleTime = time.Duration(dbConfig.MaxConnIdleTime) * time.Minute
	config.HealthCheckPeriod = time.Duration(dbConfig.HealthCheckPeriod) * time.Minute
	config.ConnConfig.ConnectTimeout = time.Duration(dbConfig.ConnectTimeout) * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{
		Pool: pool,
	}, nil
}

//   - NewTransactionManager(db *sql.DB) TransactionManager - Transaction handling
//   - HealthCheck(db *sql.DB) error - Database connectivity check
//   - Close(db *sql.DB) error - Graceful connection cleanup
// TODO: Step 2 - Implement migration utilities:
//   - MigrationRunner interface (Up, Down, Status, Version)
//   - PostgresMigrationRunner struct with file-based migrations
//   - RunMigrations(db *sql.DB, migrationsPath string) error
//   - RollbackMigration(db *sql.DB, version int) error
//   - GetMigrationStatus(db *sql.DB) ([]MigrationInfo, error)
// TODO: Step 3 - Create query utilities:
//   - QueryBuilder for dynamic SQL generation
//   - PreparedStatementManager for statement caching
//   - BulkInsert(db *sql.DB, table string, data []interface{}) error
//   - Paginate(query string, offset, limit int) string
// TODO: Step 4 - Implement transaction utilities:
//   - WithTransaction(db *sql.DB, fn func(*sql.Tx) error) error
//   - TransactionContext for request-scoped transactions
//   - DeadlockRetry wrapper for handling deadlocks
// TODO: Step 5 - Add monitoring and metrics:
//   - ConnectionMetrics (active, idle, wait time)
//   - QueryMetrics (execution time, row counts)
//   - ErrorMetrics (connection errors, query failures)
// TODO: Step 6 - Implement database utilities:
//   - TableExists(db *sql.DB, tableName string) bool
//   - ColumnExists(db *sql.DB, table, column string) bool
//   - CreateIndex(db *sql.DB, table, column string) error
// TODO: Step 7 - Add connection pooling optimizations:
//   - Dynamic pool sizing based on load
//   - Connection lifecycle management
//   - Pool warming strategies

//Flow: Service initialization -> postgres.go -> PostgreSQL database
// Used by: All services requiring database connectivity and operations
