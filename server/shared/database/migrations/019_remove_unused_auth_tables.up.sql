-- Remove unused authentication and management tables

-- Drop unused OAuth tables
DROP TABLE IF EXISTS oauth_states CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;

-- Drop unused JWT and session management tables
DROP TABLE IF EXISTS jwt_blacklist CASCADE;
DROP TABLE IF EXISTS token_usage_stats CASCADE;
DROP TABLE IF EXISTS auth_audit_log CASCADE;
DROP TABLE IF EXISTS rate_limits CASCADE;

-- Drop unused two-factor authentication tables
DROP TABLE IF EXISTS two_factor_confirmations CASCADE;
DROP TABLE IF EXISTS two_factor_tokens CASCADE;

-- Drop unused verification tables
DROP TABLE IF EXISTS verification_tokens CASCADE;

-- Drop unused session columns
ALTER TABLE sessions DROP COLUMN IF EXISTS access_token_jti;
ALTER TABLE sessions DROP COLUMN IF EXISTS refresh_token_id;
ALTER TABLE sessions DROP COLUMN IF EXISTS ip_address;
ALTER TABLE sessions DROP COLUMN IF EXISTS user_agent;
ALTER TABLE sessions DROP COLUMN IF EXISTS is_active;
ALTER TABLE sessions DROP COLUMN IF EXISTS last_activity_at;

-- Drop sessions table entirely (not used in code)
DROP TABLE IF EXISTS sessions CASCADE;

-- Drop unused activity tracking
DROP TABLE IF EXISTS login_activities CASCADE;

-- Drop subroles and related column
ALTER TABLE users DROP COLUMN IF EXISTS subrole_id;
DROP TABLE IF EXISTS subroles CASCADE;

-- Drop unused cleanup function (references removed tables)
DROP FUNCTION IF EXISTS cleanup_expired_tokens();
