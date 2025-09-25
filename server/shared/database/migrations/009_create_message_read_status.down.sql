DROP VIEW IF EXISTS unread_message_counts;
DROP INDEX IF EXISTS idx_read_status_read_at;
DROP INDEX IF EXISTS idx_read_status_user_id;
DROP TABLE IF EXISTS message_read_status;