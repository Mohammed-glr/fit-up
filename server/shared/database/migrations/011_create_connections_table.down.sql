
DROP TRIGGER IF EXISTS update_connections_updated_at ON connections;


DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_connections_last_ping;
DROP INDEX IF EXISTS idx_connections_connected_at;
DROP INDEX IF EXISTS idx_connections_conversation_ids;

DROP TABLE IF EXISTS connections;