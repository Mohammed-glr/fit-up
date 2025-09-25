CREATE TABLE IF NOT EXISTS connections (
    user_id VARCHAR(255) PRIMARY KEY,
    conversation_ids TEXT[] DEFAULT '{}',
    last_ping TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    connected_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_connections_last_ping ON connections(last_ping);
CREATE INDEX IF NOT EXISTS idx_connections_connected_at ON connections(connected_at);
CREATE INDEX IF NOT EXISTS idx_connections_conversation_ids ON connections USING GIN(conversation_ids);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_connections_updated_at 
    BEFORE UPDATE ON connections 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE connections IS 'Tracks active WebSocket connections for real-time messaging';
COMMENT ON COLUMN connections.user_id IS 'Unique identifier for the connected user';
COMMENT ON COLUMN connections.conversation_ids IS 'Array of conversation IDs the user is subscribed to';
COMMENT ON COLUMN connections.last_ping IS 'Timestamp of last ping/heartbeat from the client';
COMMENT ON COLUMN connections.connected_at IS 'Timestamp when the connection was established';
COMMENT ON COLUMN connections.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN connections.updated_at IS 'Record last update timestamp';