DROP TABLE IF EXISTS user_roles_cache;
DROP FUNCTION IF EXISTS update_role_sync_timestamp;
DROP TRIGGER IF EXISTS trigger_update_role_sync ON user_roles_cache;
