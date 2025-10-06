package repository

import (
	"context"
)

func (s *Store) MarkMessageAsRead(ctx context.Context, messageID int64, userID string) error {
	q := `
		INSERT INTO message_read_status (message_id, user_id, read_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (message_id, user_id) DO NOTHING
	`
	_, err := s.db.Exec(ctx, q, messageID, userID)
	return err
}

func (s *Store) MarkAllAsRead(ctx context.Context, conversationID int, userID string) error {
	q := `
		INSERT INTO message_read_status (message_id, user_id, read_at)
		SELECT m.message_id, $1, NOW()
		FROM messages m
		WHERE m.conversation_id = $2 AND m.sender_id != $1
		ON CONFLICT (message_id, user_id) DO NOTHING
	`
	_, err := s.db.Exec(ctx, q, userID, conversationID)
	return err
}

func (s *Store) CountUnreadMessages(ctx context.Context, conversationID int, userID string) (int, error) {
	q := `
		SELECT COUNT(*)
		FROM messages m
		LEFT JOIN message_read_status rs ON m.message_id = rs.message_id AND rs.user_id = $2
		WHERE m.conversation_id = $1 
		  AND m.sender_id != $2
		  AND m.is_deleted = false
		  AND rs.read_at IS NULL
	`
	var count int
	err := s.db.QueryRow(ctx, q, conversationID, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
