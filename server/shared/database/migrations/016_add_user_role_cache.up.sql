CREATE TABLE IF NOT EXISTS user_roles_cache (
    auth_user_id TEXT PRIMARY KEY,
    role TEXT NOT NULL DEFAULT 'user',
    last_synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_role CHECK (role IN ('user', 'coach', 'admin'))
);

CREATE INDEX idx_user_roles_cache_role ON user_roles_cache(role);
CREATE INDEX idx_user_roles_cache_synced ON user_roles_cache(last_synced_at);

CREATE OR REPLACE FUNCTION update_role_sync_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_synced_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_role_sync
BEFORE UPDATE ON user_roles_cache
FOR EACH ROW
EXECUTE FUNCTION update_role_sync_timestamp();