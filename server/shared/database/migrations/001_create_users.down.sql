DROP TRIGGER IF EXISTS set_timestamp ON users;
DROP FUNCTION IF EXISTS trigger_set_timestamp;

DROP TABLE IF EXISTS subroles;
DROP TABLE IF EXISTS login_activities;
DROP TABLE IF EXISTS two_factor_confirmations;
DROP TABLE IF EXISTS two_factor_tokens;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_role;
