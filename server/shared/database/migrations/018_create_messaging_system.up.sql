
CREATE TABLE IF NOT EXISTS conversations (
    conversation_id SERIAL PRIMARY KEY,
    coach_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_message_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_archived BOOLEAN DEFAULT FALSE,
    
    CONSTRAINT unique_coach_client_conversation UNIQUE (coach_id, client_id),
    
    CONSTRAINT check_different_users CHECK (coach_id != client_id)
);

CREATE INDEX idx_conversations_coach_id ON conversations(coach_id) WHERE is_archived = FALSE;
CREATE INDEX idx_conversations_client_id ON conversations(client_id) WHERE is_archived = FALSE;
CREATE INDEX idx_conversations_last_message ON conversations(last_message_at DESC);

COMMENT ON TABLE conversations IS 'Chat conversations between coaches and their assigned clients';
COMMENT ON COLUMN conversations.is_archived IS 'Soft delete - hides conversation from active list';

CREATE TABLE IF NOT EXISTS messages (
    message_id BIGSERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(conversation_id) ON DELETE CASCADE,
    sender_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_text TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP WITH TIME ZONE,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    reply_to_message_id BIGINT REFERENCES messages(message_id) ON DELETE SET NULL,
    
    CONSTRAINT check_message_not_empty CHECK (LENGTH(TRIM(message_text)) > 0 OR is_deleted = TRUE)
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id, sent_at DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_reply_to ON messages(reply_to_message_id) WHERE reply_to_message_id IS NOT NULL;
CREATE INDEX idx_messages_active ON messages(conversation_id) WHERE is_deleted = FALSE;

COMMENT ON TABLE messages IS 'Individual messages within coach-client conversations';
COMMENT ON COLUMN messages.is_deleted IS 'Soft delete - message still exists but content hidden';
COMMENT ON COLUMN messages.reply_to_message_id IS 'References another message for threaded replies';

CREATE TABLE IF NOT EXISTS message_read_status (
    read_status_id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES messages(message_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_user_message_read UNIQUE (message_id, user_id)
);

CREATE INDEX idx_message_read_status_message ON message_read_status(message_id);
CREATE INDEX idx_message_read_status_user ON message_read_status(user_id);

COMMENT ON TABLE message_read_status IS 'Tracks when users have read specific messages';

CREATE TABLE IF NOT EXISTS message_attachments (
    attachment_id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES messages(message_id) ON DELETE CASCADE,
    attachment_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_size INTEGER,
    mime_type VARCHAR(100),
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    metadata JSONB DEFAULT '{}'::jsonb,
    
    CONSTRAINT check_attachment_type CHECK (attachment_type IN (
        'image', 'document', 'workout_plan'
    ))
);

CREATE INDEX idx_message_attachments_message ON message_attachments(message_id);
CREATE INDEX idx_message_attachments_type ON message_attachments(attachment_type);

COMMENT ON TABLE message_attachments IS 'Files and media attached to messages';
COMMENT ON COLUMN message_attachments.metadata IS 'Additional data like dimensions for images, duration for videos, etc.';


CREATE OR REPLACE FUNCTION update_conversation_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE conversations 
    SET 
        last_message_at = NEW.sent_at,
        updated_at = NEW.sent_at
    WHERE conversation_id = NEW.conversation_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_conversation_timestamp
    AFTER INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_conversation_timestamp();


CREATE OR REPLACE VIEW conversation_overview AS
SELECT 
    c.conversation_id,
    c.coach_id,
    c.client_id,
    u_coach.name AS coach_name,
    u_coach.image AS coach_image,
    u_client.name AS client_name,
    u_client.image AS client_image,
    c.created_at,
    c.last_message_at,
    c.is_archived,
    m.message_text AS last_message_text,
    m.sender_id AS last_message_sender_id,
    m.sent_at AS last_message_sent_at,
    (SELECT COUNT(*) FROM messages WHERE conversation_id = c.conversation_id AND is_deleted = FALSE) AS total_messages
FROM conversations c
LEFT JOIN users u_coach ON c.coach_id = u_coach.id
LEFT JOIN users u_client ON c.client_id = u_client.id
LEFT JOIN LATERAL (
    SELECT message_text, sender_id, sent_at
    FROM messages
    WHERE conversation_id = c.conversation_id 
      AND is_deleted = FALSE
    ORDER BY sent_at DESC
    LIMIT 1
) m ON true
ORDER BY c.last_message_at DESC;

COMMENT ON VIEW conversation_overview IS 'Convenient view showing conversation details with last message';
