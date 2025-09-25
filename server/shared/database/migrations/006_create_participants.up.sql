CREATE TABLE participants (
    conversation_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'moderator', 'member', 'readonly')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'left', 'kicked', 'banned')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    is_muted BOOLEAN DEFAULT FALSE,
    muted_until TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    PRIMARY KEY (conversation_id, user_id),
    CONSTRAINT fk_participant_conversation FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_participant_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_participants_user_id ON participants(user_id);
CREATE INDEX idx_participants_conversation_active ON participants(conversation_id) WHERE status = 'active';