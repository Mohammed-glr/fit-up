-- Drop JWT and session management tables and columns

-- Drop functions first
DROP FUNCTION IF EXISTS cleanup_expired_tokens();

-- Drop tables
DROP TABLE IF EXISTS token_usage_stats;
DROP TABLE IF EXISTS rate_limits;
DROP TABLE IF EXISTS auth_audit_log;
DROP TABLE IF EXISTS jwt_blacklist;
DROP TABLE IF EXISTS jwt_refresh_tokens;

-- Remove added columns from sessions table
ALTER TABLE sessions DROP COLUMN IF EXISTS access_token_jti;
ALTER TABLE sessions DROP COLUMN IF EXISTS refresh_token_id;
ALTER TABLE sessions DROP COLUMN IF EXISTS ip_address;
ALTER TABLE sessions DROP COLUMN IF EXISTS user_agent;
ALTER TABLE sessions DROP COLUMN IF EXISTS is_active;
ALTER TABLE sessions DROP COLUMN IF EXISTS last_activity_at;
