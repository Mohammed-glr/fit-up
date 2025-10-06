
DROP VIEW IF EXISTS conversation_overview;

DROP TRIGGER IF EXISTS trigger_update_conversation_timestamp ON messages;

DROP FUNCTION IF EXISTS update_conversation_timestamp();

DROP TABLE IF EXISTS message_attachments CASCADE;
DROP TABLE IF EXISTS message_read_status CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS conversations CASCADE;
