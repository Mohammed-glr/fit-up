-- Create message_attachments table
CREATE TABLE message_attachments (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id TEXT NOT NULL,
    file_name TEXT NOT NULL,
    original_name TEXT NOT NULL,
    file_type VARCHAR(20) NOT NULL CHECK (file_type IN ('image', 'document', 'audio', 'video', 'archive', 'other')),
    file_size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    url TEXT NOT NULL,
    thumbnail_url TEXT DEFAULT NULL,
    width INTEGER DEFAULT NULL,
    height INTEGER DEFAULT NULL,
    duration INTEGER DEFAULT NULL, 
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    CONSTRAINT fk_attachment_message FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE,
    CONSTRAINT chk_file_size_positive CHECK (file_size > 0)
);

CREATE INDEX idx_attachments_message_id ON message_attachments(message_id);
CREATE INDEX idx_attachments_file_type ON message_attachments(file_type);
CREATE INDEX idx_attachments_uploaded_at ON message_attachments(uploaded_at DESC);