-- Drop trigger first
DROP TRIGGER IF EXISTS set_timestamp ON users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- Drop indexes
DROP INDEX IF EXISTS idx_two_factor_tokens_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP INDEX IF EXISTS idx_verification_tokens_token;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables in correct order (tables with foreign keys first)
DROP TABLE IF EXISTS login_activities CASCADE;
DROP TABLE IF EXISTS two_factor_confirmations CASCADE;
DROP TABLE IF EXISTS two_factor_tokens CASCADE;
DROP TABLE IF EXISTS password_reset_tokens CASCADE;
DROP TABLE IF EXISTS verification_tokens CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;

-- Drop users table (this will drop the subrole_id foreign key)
DROP TABLE IF EXISTS users CASCADE;

-- Now we can drop subroles
DROP TABLE IF EXISTS subroles CASCADE;

-- Drop enum type
DROP TYPE IF EXISTS user_role CASCADE;
