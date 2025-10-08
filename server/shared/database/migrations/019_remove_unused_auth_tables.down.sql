-- Rollback: Recreate unused authentication and management tables
-- Note: This rollback recreates the tables but data will be lost

-- Recreate sessions table
CREATE TABLE sessions (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    expires TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    access_token_jti TEXT,
    refresh_token_id TEXT,
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_session FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_access_token_jti ON sessions(access_token_jti) WHERE access_token_jti IS NOT NULL;
CREATE INDEX idx_sessions_refresh_token_id ON sessions(refresh_token_id) WHERE refresh_token_id IS NOT NULL;
CREATE INDEX idx_sessions_is_active ON sessions(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_sessions_last_activity_at ON sessions(last_activity_at);

-- Recreate verification_tokens table
CREATE TABLE verification_tokens (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE (email, token)
);

CREATE INDEX idx_verification_tokens_token ON verification_tokens(token);

-- Recreate two_factor_tokens table
CREATE TABLE two_factor_tokens (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE (email, token)
);

CREATE INDEX idx_two_factor_tokens_token ON two_factor_tokens(token);

-- Recreate two_factor_confirmations table
CREATE TABLE two_factor_confirmations (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL UNIQUE,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Recreate login_activities table
CREATE TABLE login_activities (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    location TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_login_activity FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Recreate subroles table
CREATE TABLE subroles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS subrole_id INT REFERENCES subroles(id);

-- Recreate JWT blacklist table
CREATE TABLE jwt_blacklist (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    jti TEXT NOT NULL UNIQUE,
    token_hash TEXT NOT NULL,
    user_id TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    blacklisted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reason TEXT DEFAULT 'logout',
    CONSTRAINT fk_jwt_blacklist_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_jwt_blacklist_jti ON jwt_blacklist(jti);
CREATE INDEX idx_jwt_blacklist_token_hash ON jwt_blacklist(token_hash);
CREATE INDEX idx_jwt_blacklist_user_id ON jwt_blacklist(user_id);
CREATE INDEX idx_jwt_blacklist_expires_at ON jwt_blacklist(expires_at);

-- Recreate auth_audit_log table
CREATE TABLE auth_audit_log (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT,
    action TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_auth_audit_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_auth_audit_log_user_id ON auth_audit_log(user_id);
CREATE INDEX idx_auth_audit_log_action ON auth_audit_log(action);
CREATE INDEX idx_auth_audit_log_created_at ON auth_audit_log(created_at);
CREATE INDEX idx_auth_audit_log_ip_address ON auth_audit_log(ip_address);

-- Recreate rate_limits table
CREATE TABLE rate_limits (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    attempts INTEGER DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE,
    blocked_until TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (identifier, endpoint, window_start)
);

CREATE INDEX idx_rate_limits_identifier_endpoint ON rate_limits(identifier, endpoint);
CREATE INDEX idx_rate_limits_window_end ON rate_limits(window_end);
CREATE INDEX idx_rate_limits_is_blocked ON rate_limits(is_blocked) WHERE is_blocked = TRUE;

CREATE TRIGGER set_timestamp_rate_limits
    BEFORE UPDATE ON rate_limits
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Recreate token_usage_stats table
CREATE TABLE token_usage_stats (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT,
    token_type TEXT NOT NULL,
    action TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_token_stats_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_token_usage_stats_user_id ON token_usage_stats(user_id);
CREATE INDEX idx_token_usage_stats_token_type ON token_usage_stats(token_type);
CREATE INDEX idx_token_usage_stats_action ON token_usage_stats(action);
CREATE INDEX idx_token_usage_stats_created_at ON token_usage_stats(created_at);

-- Recreate OAuth tables
CREATE TABLE oauth_states (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    state TEXT NOT NULL UNIQUE,
    provider TEXT NOT NULL,
    redirect_url TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_oauth_states_expires_at ON oauth_states(expires_at);
CREATE INDEX idx_oauth_states_provider ON oauth_states(provider);

CREATE TABLE accounts (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    provider TEXT NOT NULL,
    provider_account_id TEXT NOT NULL,
    refresh_token TEXT,
    access_token TEXT,
    expires_at INTEGER,
    token_type TEXT,
    scope TEXT,
    id_token TEXT,
    session_state TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_accounts_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(provider, provider_account_id)
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_provider ON accounts(provider);
CREATE INDEX idx_accounts_provider_account_id ON accounts(provider_account_id);

-- Recreate cleanup function
CREATE OR REPLACE FUNCTION cleanup_expired_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER := 0;
    temp_count INTEGER := 0;
BEGIN
    DELETE FROM jwt_refresh_tokens
    WHERE expires_at < NOW() OR is_revoked = TRUE;
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM jwt_blacklist
    WHERE expires_at < NOW();
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM sessions
    WHERE expires < NOW() OR is_active = FALSE;
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM rate_limits
    WHERE window_end < NOW() - INTERVAL '24 hours';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM auth_audit_log
    WHERE created_at < NOW() - INTERVAL '90 days';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    DELETE FROM token_usage_stats
    WHERE created_at < NOW() - INTERVAL '30 days';
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
