package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tdmdh/fit-up-server/services/auth-service/internal/interfaces"
)

// CleanupService handles cleanup of expired tokens and old audit logs
type CleanupService struct {
	store interfaces.UserStore
}

// NewCleanupService creates a new cleanup service
func NewCleanupService(store interfaces.UserStore) *CleanupService {
	return &CleanupService{
		store: store,
	}
}

// CleanupExpiredTokens removes expired refresh tokens and blacklist entries
func (c *CleanupService) CleanupExpiredTokens(ctx context.Context) error {
	log.Println("Starting token cleanup job...")

	err := c.store.CleanupExpiredRefreshTokens(ctx)
	if err != nil {
		log.Printf("Failed to cleanup expired refresh tokens: %v", err)
		return fmt.Errorf("failed to cleanup expired refresh tokens: %w", err)
	}

	err = c.store.CleanupExpiredTokens(ctx)
	if err != nil {
		log.Printf("Failed to cleanup expired blacklist tokens: %v", err)
		return fmt.Errorf("failed to cleanup expired blacklist tokens: %w", err)
	}

	log.Println("Token cleanup job completed successfully")
	return nil
}

// CleanupOldAuditLogs removes audit logs older than specified duration
func (c *CleanupService) CleanupOldAuditLogs(ctx context.Context, retentionPeriod time.Duration) error {
	log.Printf("Starting audit log cleanup job (retention: %v)...", retentionPeriod)

	cutoffTime := time.Now().Add(-retentionPeriod)

	if auditStore, ok := c.store.(interfaces.AuditStore); ok {
		err := auditStore.CleanupOldAuditLogs(ctx, cutoffTime)
		if err != nil {
			log.Printf("Failed to cleanup old audit logs: %v", err)
			return fmt.Errorf("failed to cleanup old audit logs: %w", err)
		}

		log.Printf("Audit log cleanup job completed successfully (cutoff: %v)", cutoffTime)
	} else {
		log.Println("Store does not implement AuditStore interface, skipping audit log cleanup")
	}

	return nil
}

// RunPeriodicCleanup runs cleanup jobs periodically
func (c *CleanupService) RunPeriodicCleanup(ctx context.Context, tokenCleanupInterval, auditCleanupInterval time.Duration, auditRetentionPeriod time.Duration) {
	tokenTicker := time.NewTicker(tokenCleanupInterval)
	auditTicker := time.NewTicker(auditCleanupInterval)

	defer tokenTicker.Stop()
	defer auditTicker.Stop()

	log.Printf("Starting periodic cleanup service (token cleanup: %v, audit cleanup: %v, audit retention: %v)",
		tokenCleanupInterval, auditCleanupInterval, auditRetentionPeriod)

	// Run initial cleanup
	if err := c.CleanupExpiredTokens(ctx); err != nil {
		log.Printf("Initial token cleanup failed: %v", err)
	}

	if err := c.CleanupOldAuditLogs(ctx, auditRetentionPeriod); err != nil {
		log.Printf("Initial audit log cleanup failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup service shutting down...")
			return

		case <-tokenTicker.C:
			if err := c.CleanupExpiredTokens(ctx); err != nil {
				log.Printf("Periodic token cleanup failed: %v", err)
			}

		case <-auditTicker.C:
			if err := c.CleanupOldAuditLogs(ctx, auditRetentionPeriod); err != nil {
				log.Printf("Periodic audit log cleanup failed: %v", err)
			}
		}
	}
}

// GetCleanupStats returns statistics about cleanup operations
func (c *CleanupService) GetCleanupStats(ctx context.Context) (*CleanupStats, error) {
	return &CleanupStats{
		LastTokenCleanup:      time.Now(), // Would store this in DB or memory
		LastAuditCleanup:      time.Now(), // Would store this in DB or memory
		TokensCleanedToday:    0,          // Would query from audit logs
		AuditLogsCleanedToday: 0,          // Would query from audit logs
	}, nil
}

// CleanupStats represents cleanup operation statistics
type CleanupStats struct {
	LastTokenCleanup      time.Time `json:"last_token_cleanup"`
	LastAuditCleanup      time.Time `json:"last_audit_cleanup"`
	TokensCleanedToday    int       `json:"tokens_cleaned_today"`
	AuditLogsCleanedToday int       `json:"audit_logs_cleaned_today"`
}

var (
	// Clean up expired tokens every hour
	DefaultTokenCleanupInterval = time.Hour

	// Clean up old audit logs once per day
	DefaultAuditCleanupInterval = 24 * time.Hour

	// Keep audit logs for 90 days
	DefaultAuditRetentionPeriod = 90 * 24 * time.Hour
)
