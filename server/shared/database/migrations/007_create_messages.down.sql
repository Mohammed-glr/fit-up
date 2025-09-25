DROP TRIGGER IF EXISTS set_timestamp_messages ON messages;
DROP INDEX IF EXISTS idx_messages_not_deleted;
DROP INDEX IF EXISTS idx_messages_reply_to;
DROP INDEX IF EXISTS idx_messages_sender;
DROP INDEX IF EXISTS idx_messages_conversation_created;
DROP TABLE IF EXISTS messages;