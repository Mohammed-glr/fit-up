-- Create message_read_status table
CREATE TABLE message_read_status (
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id, user_id),
    CONSTRAINT fk_read_status_message FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE,
    CONSTRAINT fk_read_status_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_read_status_user_id ON message_read_status(user_id);
CREATE INDEX idx_read_status_read_at ON message_read_status(read_at DESC);

CREATE VIEW unread_message_counts AS
SELECT 
    p.user_id,
    p.conversation_id,
    COUNT(m.id) as unread_count,
    MAX(m.created_at) as latest_unread_at
FROM participants p
JOIN messages m ON m.conversation_id = p.conversation_id
LEFT JOIN message_read_status mrs ON mrs.message_id = m.id AND mrs.user_id = p.user_id
WHERE p.status = 'active' 
    AND m.deleted_at IS NULL 
    AND m.sender_id != p.user_id 
    AND mrs.message_id IS NULL
GROUP BY p.user_id, p.conversation_id;