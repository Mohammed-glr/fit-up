package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/types"
)

func (s *Store) CreateMessage(ctx context.Context, conversationID int, senderID, messageText string, replyToMessageID *int64) (*types.Message, error) {
	q := `
		INSERT INTO messages (conversation_id, sender_id, message_text, reply_to_message_id)
		VALUES ($1, $2, $3, $4)
		RETURNING message_id, conversation_id, sender_id, message_text, reply_to_message_id, sent_at, edited_at
	`

	var msg types.Message
	err := s.db.QueryRow(ctx, q, conversationID, senderID, messageText, replyToMessageID).Scan(
		&msg.MessageID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.MessageText,
		&msg.ReplyToMessageID,
		&msg.SentAt,
		&msg.EditedAt,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *Store) GetMessageByID(ctx context.Context, messageID int64) (*types.Message, error) {
	q := `
		SELECT message_id, conversation_id, sender_id, message_text, reply_to_message_id, sent_at, edited_at, is_deleted, deleted_at
		FROM messages
		WHERE message_id = $1
	`

	var msg types.Message
	err := s.db.QueryRow(ctx, q, messageID).Scan(
		&msg.MessageID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.MessageText,
		&msg.ReplyToMessageID,
		&msg.SentAt,
		&msg.EditedAt,
		&msg.IsDeleted,
		&msg.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *Store) ListMessages(ctx context.Context, conversationID int, limit, offset int) ([]types.MessageWithDetails, error) {
	q := `
		SELECT 
			m.message_id, m.conversation_id, m.sender_id, m.message_text, m.reply_to_message_id, 
			m.sent_at, m.edited_at, m.is_deleted, m.deleted_at,
			u.name AS sender_name, u.avatar_url AS sender_image
		FROM messages m
		JOIN users u ON m.sender_id = u.user_id
		WHERE m.conversation_id = $1 AND m.is_deleted = false
		ORDER BY m.sent_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(ctx, q, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []types.MessageWithDetails
	for rows.Next() {
		var msg types.MessageWithDetails
		err := rows.Scan(
			&msg.MessageID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.MessageText,
			&msg.ReplyToMessageID,
			&msg.SentAt,
			&msg.EditedAt,
			&msg.IsDeleted,
			&msg.DeletedAt,
			&msg.SenderName,
			&msg.SenderImage,
		)
		if err != nil {
			return nil, err
		}
		// Attachments will be populated separately by service layer
		messages = append(messages, msg)
	}
	return messages, nil
}

func (s *Store) UpdateMessage(ctx context.Context, messageID int64, newText string) error {
	q := `
		UPDATE messages
		SET message_text = $1, edited_at = NOW()
		WHERE message_id = $2 AND is_deleted = false
	`
	_, err := s.db.Exec(ctx, q, newText, messageID)
	return err
}

func (s *Store) DeleteMessage(ctx context.Context, messageID int64) error {
	q := `
		UPDATE messages
		SET is_deleted = true, deleted_at = NOW(), message_text = '[Message deleted]'
		WHERE message_id = $1
	`
	_, err := s.db.Exec(ctx, q, messageID)
	return err
}
