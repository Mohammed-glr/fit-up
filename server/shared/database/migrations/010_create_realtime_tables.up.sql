CREATE TABLE typing_indicators (
    conversation_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    is_typing BOOLEAN DEFAULT FALSE,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id),
    CONSTRAINT fk_typing_conversation FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_typing_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE user_presence (
    user_id TEXT PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (status IN ('online', 'away', 'busy', 'offline')),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    device_info JSONB DEFAULT '{}',
    CONSTRAINT fk_presence_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE conversation_settings (
    conversation_id TEXT PRIMARY KEY,
    is_private BOOLEAN DEFAULT FALSE,
    allow_invites BOOLEAN DEFAULT TRUE,
    mute_notifications BOOLEAN DEFAULT FALSE,
    retention_days INTEGER DEFAULT NULL,
    max_participants INTEGER DEFAULT NULL,
    allow_message_deletion BOOLEAN DEFAULT TRUE,
    allow_message_editing BOOLEAN DEFAULT TRUE,
    custom_settings JSONB DEFAULT '{}',
    CONSTRAINT fk_settings_conversation FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE INDEX idx_typing_indicators_conversation ON typing_indicators(conversation_id) WHERE is_typing = TRUE;
CREATE INDEX idx_user_presence_status ON user_presence(status);
CREATE INDEX idx_user_presence_last_seen ON user_presence(last_seen DESC);