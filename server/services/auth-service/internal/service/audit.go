package service

import (
	"context"
	"time"

	"github.com/tdmdh/lornian-backend/services/auth-service/internal/interfaces"
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
)

// AuditLogger handles authentication audit logging
type AuditLogger struct {
	store interfaces.AuditStore
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(store interfaces.AuditStore) *AuditLogger {
	return &AuditLogger{
		store: store,
	}
}

// LogAuthEvent logs an authentication event
func (a *AuditLogger) LogAuthEvent(ctx context.Context, event *types.AuthAuditEvent) error {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	return a.store.CreateAuditLog(ctx, event)
}

// LogLogin logs a login attempt
func (a *AuditLogger) LogLogin(ctx context.Context, userID, email, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "login",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	if details == nil {
		event.Details = map[string]interface{}{
			"email": email,
		}
	} else {
		event.Details["email"] = email
	}

	return a.LogAuthEvent(ctx, event)
}

// LogLogout logs a logout event
func (a *AuditLogger) LogLogout(ctx context.Context, userID, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "logout",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return a.LogAuthEvent(ctx, event)
}

// LogTokenRefresh logs a token refresh event
func (a *AuditLogger) LogTokenRefresh(ctx context.Context, userID, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "token_refresh",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return a.LogAuthEvent(ctx, event)
}

// LogPasswordChange logs a password change event
func (a *AuditLogger) LogPasswordChange(ctx context.Context, userID, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "password_change",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return a.LogAuthEvent(ctx, event)
}

// LogPasswordReset logs a password reset event
func (a *AuditLogger) LogPasswordReset(ctx context.Context, userID, email, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "password_reset",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	if details == nil {
		event.Details = map[string]interface{}{
			"email": email,
		}
	} else {
		event.Details["email"] = email
	}

	return a.LogAuthEvent(ctx, event)
}

// LogRegistration logs a user registration event
func (a *AuditLogger) LogRegistration(ctx context.Context, userID, email, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "registration",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	if details == nil {
		event.Details = map[string]interface{}{
			"email": email,
		}
	} else {
		event.Details["email"] = email
	}

	return a.LogAuthEvent(ctx, event)
}

// LogEmailVerification logs an email verification event
func (a *AuditLogger) LogEmailVerification(ctx context.Context, userID, email, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    "email_verification",
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	if details == nil {
		event.Details = map[string]interface{}{
			"email": email,
		}
	} else {
		event.Details["email"] = email
	}

	return a.LogAuthEvent(ctx, event)
}

// LogSecurityEvent logs a general security event
func (a *AuditLogger) LogSecurityEvent(ctx context.Context, userID, action, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	event := &types.AuthAuditEvent{
		UserID:    userID,
		Action:    action,
		Success:   success,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return a.LogAuthEvent(ctx, event)
}

// LogFailedLoginAttempt logs a failed login attempt with additional context
func (a *AuditLogger) LogFailedLoginAttempt(ctx context.Context, email, ipAddress, userAgent, reason string) error {
	details := map[string]interface{}{
		"email":  email,
		"reason": reason,
	}

	event := &types.AuthAuditEvent{
		UserID:    "", // No user ID for failed login
		Action:    "failed_login",
		Success:   false,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return a.LogAuthEvent(ctx, event)
}

// GetAuditLogs retrieves audit logs for a user
func (a *AuditLogger) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]*types.AuthAuditEvent, error) {
	return a.store.GetAuditLogsByUserID(ctx, userID, limit, offset)
}

// GetSecurityEvents retrieves security events within a time range
func (a *AuditLogger) GetSecurityEvents(ctx context.Context, startTime, endTime time.Time, limit, offset int) ([]*types.AuthAuditEvent, error) {
	return a.store.GetAuditLogsByTimeRange(ctx, startTime, endTime, limit, offset)
}

// DetectSuspiciousActivity detects potentially suspicious authentication patterns
func (a *AuditLogger) DetectSuspiciousActivity(ctx context.Context, ipAddress string, timeWindow time.Duration) (*types.SuspiciousActivityReport, error) {
	endTime := time.Now()
	startTime := endTime.Add(-timeWindow)

	events, err := a.store.GetAuditLogsByIPAddress(ctx, ipAddress, startTime, endTime)
	if err != nil {
		return nil, err
	}

	report := &types.SuspiciousActivityReport{
		IPAddress:    ipAddress,
		TimeWindow:   timeWindow,
		TotalEvents:  len(events),
		FailedLogins: 0,
		UserAgents:   make(map[string]int),
		Actions:      make(map[string]int),
	}

	for _, event := range events {
		report.Actions[event.Action]++

		if event.UserAgent != "" {
			report.UserAgents[event.UserAgent]++
		}

		if event.Action == "login" && !event.Success {
			report.FailedLogins++
		}

		if event.Action == "failed_login" {
			report.FailedLogins++
		}
	}

	// Determine if activity is suspicious
	report.IsSuspicious = report.FailedLogins > 5 || len(report.UserAgents) > 3

	return report, nil
}
