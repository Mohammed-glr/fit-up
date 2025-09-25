-- Add JWT and session management tables

-- JWT Refresh Tokens table
CREATE TABLE jwt_refresh_tokens (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE, -- Store hashed refresh token for security
    access_token_jti TEXT, -- JWT ID claim for linking with access token
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    user_agent TEXT,
    ip_address INET,
    CONSTRAINT fk_jwt_refresh_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- JWT Token Blacklist (for logout and token revocation)
CREATE TABLE jwt_blacklist (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    jti TEXT NOT NULL UNIQUE, -- JWT ID claim
    token_hash TEXT NOT NULL, -- Hashed token for security
    user_id TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- When the token would have expired
    blacklisted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reason TEXT DEFAULT 'logout', -- logout, security_breach, admin_revoke, etc.
    CONSTRAINT fk_jwt_blacklist_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Enhanced Sessions table (extend existing sessions with JWT support)
-- Add columns to existing sessions table
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS access_token_jti TEXT;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS refresh_token_id TEXT;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS ip_address INET;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS user_agent TEXT;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- User Authentication Audit Log
CREATE TABLE auth_audit_log (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT,
    action TEXT NOT NULL, -- login, logout, token_refresh, password_change, etc.
    success BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    details JSONB, -- Additional context (error codes, token info, etc.)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_auth_audit_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Rate Limiting table (for JWT endpoints and auth operations)
CREATE TABLE rate_limits (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier TEXT NOT NULL, -- IP address, user_id, or composite key
    endpoint TEXT NOT NULL, -- /login, /refresh-token, /validate-token, etc.
    attempts INTEGER DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE,
    blocked_until TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (identifier, endpoint, window_start)
);

-- Token Usage Statistics (for monitoring and analytics)
CREATE TABLE token_usage_stats (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT,
    token_type TEXT NOT NULL, -- access, refresh
    action TEXT NOT NULL, -- generate, validate, refresh, revoke
    success BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_token_stats_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Add indexes for performance
CREATE INDEX idx_jwt_refresh_tokens_user_id ON jwt_refresh_tokens(user_id);
CREATE INDEX idx_jwt_refresh_tokens_token_hash ON jwt_refresh_tokens(token_hash);
CREATE INDEX idx_jwt_refresh_tokens_expires_at ON jwt_refresh_tokens(expires_at);
CREATE INDEX idx_jwt_refresh_tokens_is_revoked ON jwt_refresh_tokens(is_revoked) WHERE is_revoked = FALSE;

CREATE INDEX idx_jwt_blacklist_jti ON jwt_blacklist(jti);
CREATE INDEX idx_jwt_blacklist_token_hash ON jwt_blacklist(token_hash);
CREATE INDEX idx_jwt_blacklist_user_id ON jwt_blacklist(user_id);
CREATE INDEX idx_jwt_blacklist_expires_at ON jwt_blacklist(expires_at);

CREATE INDEX idx_sessions_access_token_jti ON sessions(access_token_jti) WHERE access_token_jti IS NOT NULL;
CREATE INDEX idx_sessions_refresh_token_id ON sessions(refresh_token_id) WHERE refresh_token_id IS NOT NULL;
CREATE INDEX idx_sessions_is_active ON sessions(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_sessions_last_activity_at ON sessions(last_activity_at);

CREATE INDEX idx_auth_audit_log_user_id ON auth_audit_log(user_id);
CREATE INDEX idx_auth_audit_log_action ON auth_audit_log(action);
CREATE INDEX idx_auth_audit_log_created_at ON auth_audit_log(created_at);
CREATE INDEX idx_auth_audit_log_ip_address ON auth_audit_log(ip_address);

CREATE INDEX idx_rate_limits_identifier_endpoint ON rate_limits(identifier, endpoint);
CREATE INDEX idx_rate_limits_window_end ON rate_limits(window_end);
CREATE INDEX idx_rate_limits_is_blocked ON rate_limits(is_blocked) WHERE is_blocked = TRUE;

CREATE INDEX idx_token_usage_stats_user_id ON token_usage_stats(user_id);
CREATE INDEX idx_token_usage_stats_token_type ON token_usage_stats(token_type);
CREATE INDEX idx_token_usage_stats_action ON token_usage_stats(action);
CREATE INDEX idx_token_usage_stats_created_at ON token_usage_stats(created_at);

-- Add triggers for updated_at timestamps
CREATE TRIGGER set_timestamp_jwt_refresh_tokens
    BEFORE UPDATE ON jwt_refresh_tokens
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_rate_limits
    BEFORE UPDATE ON rate_limits
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Add a function to clean up expired tokens
CREATE OR REPLACE FUNCTION cleanup_expired_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER := 0;
    temp_count INTEGER := 0;
BEGIN
    -- Clean up expired refresh tokens
    DELETE FROM jwt_refresh_tokens
    WHERE expires_at < NOW() OR is_revoked = TRUE;
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    -- Clean up expired blacklisted tokens
    DELETE FROM jwt_blacklist
    WHERE expires_at < NOW();
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    -- Clean up expired sessions
    DELETE FROM sessions
    WHERE expires < NOW() OR is_active = FALSE;
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    -- Clean up old rate limit records (older than 24 hours)
    DELETE FROM rate_limits
    WHERE window_end < NOW() - INTERVAL '24 hours';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    -- Clean up old audit logs (older than 90 days)
    DELETE FROM auth_audit_log
    WHERE created_at < NOW() - INTERVAL '90 days';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    -- Clean up old token usage stats (older than 30 days)
    DELETE FROM token_usage_stats
    WHERE created_at < NOW() - INTERVAL '30 days';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE jwt_refresh_tokens IS 'Stores refresh tokens for JWT authentication';
COMMENT ON TABLE jwt_blacklist IS 'Blacklisted JWT tokens for logout and revocation';
COMMENT ON TABLE auth_audit_log IS 'Audit log for authentication events';
COMMENT ON TABLE rate_limits IS 'Rate limiting data for API endpoints';
COMMENT ON TABLE token_usage_stats IS 'Token usage statistics for monitoring';

COMMENT ON COLUMN jwt_refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token for security';
COMMENT ON COLUMN jwt_refresh_tokens.access_token_jti IS 'JWT ID claim linking to access token';
COMMENT ON COLUMN jwt_blacklist.jti IS 'JWT ID claim from the blacklisted token';
COMMENT ON COLUMN jwt_blacklist.token_hash IS 'SHA-256 hash of the blacklisted token';
COMMENT ON COLUMN auth_audit_log.details IS 'JSON details about the authentication event';
COMMENT ON COLUMN rate_limits.identifier IS 'IP address, user ID, or composite identifier';
COMMENT ON COLUMN rate_limits.endpoint IS 'API endpoint being rate limited';
