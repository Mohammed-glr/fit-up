package repository

import (
	"context"
	"fmt"

	"github.com/tdmdh/fit-up-server/internal/message/types"
)

func (s *Store) CreateConversation(ctx context.Context, coachID, clientID string) (*types.Conversation, error) {
	q := `
		INSERT INTO conversations (coach_id, client_id)
		VALUES ($1, $2)
		RETURNING conversation_id, coach_id, client_id, is_archived, created_at, updated_at, last_message_at
	`

	var conv types.Conversation
	err := s.db.QueryRow(ctx, q, coachID, clientID).Scan(
		&conv.ConversationID,
		&conv.CoachID,
		&conv.ClientID,
		&conv.IsArchived,
		&conv.CreatedAt,
		&conv.UpdatedAt,
		&conv.LastMessageAt,
	)
	if err != nil {
		return nil, err
	}

	return &conv, nil
}

func (s *Store) GetConversationByID(ctx context.Context, conversationID int) (*types.Conversation, error) {
	q := `
		SELECT conversation_id, coach_id, client_id, is_archived, created_at, updated_at, last_message_at
		FROM conversations
		WHERE conversation_id = $1
	`
	var conv types.Conversation
	err := s.db.QueryRow(ctx, q, conversationID).Scan(
		&conv.ConversationID,
		&conv.CoachID,
		&conv.ClientID,
		&conv.IsArchived,
		&conv.CreatedAt,
		&conv.UpdatedAt,
		&conv.LastMessageAt,
	)
	if err != nil {
		return nil, err
	}

	return &conv, nil
}

func (s *Store) GetConversationByParticipants(ctx context.Context, coachID, clientID string) (*types.Conversation, error) {
	q := `
		SELECT conversation_id, coach_id, client_id, is_archived, created_at, updated_at, last_message_at
		FROM conversations
		WHERE coach_id = $1 AND client_id = $2
	`

	var conv types.Conversation
	err := s.db.QueryRow(ctx, q, coachID, clientID).Scan(
		&conv.ConversationID,
		&conv.CoachID,
		&conv.ClientID,
		&conv.IsArchived,
		&conv.CreatedAt,
		&conv.UpdatedAt,
		&conv.LastMessageAt,
	)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (s *Store) ListConversationsByUser(ctx context.Context, userID string, includeArchived bool, limit, offset int) ([]types.ConversationOverview, int, error) {
	baseQuery := `
		SELECT 
			c.conversation_id,
			c.coach_id,
			c.client_id,
			coach.name AS coach_name,
			coach.avatar_url AS coach_image,
			client.name AS client_name,
			client.avatar_url AS client_image,
			c.created_at,
			c.last_message_at,
			c.is_archived,
			lm.message_text AS last_message_text,
			lm.sender_id AS last_message_sender_id,
			lm.sent_at AS last_message_sent_at,
			COALESCE(msg_count.total, 0) AS total_messages
		FROM conversations c
		LEFT JOIN users coach ON c.coach_id = coach.user_id
		LEFT JOIN users client ON c.client_id = client.user_id
		LEFT JOIN LATERAL (
			SELECT message_text, sender_id, sent_at
			FROM messages
			WHERE conversation_id = c.conversation_id
			ORDER BY sent_at DESC
			LIMIT 1
		) lm ON true
		LEFT JOIN LATERAL (
			SELECT COUNT(*) as total
			FROM messages
			WHERE conversation_id = c.conversation_id
		) msg_count ON true
		WHERE (c.coach_id = $1 OR c.client_id = $1)
	`

	countQuery := `
		SELECT COUNT(*)
		FROM conversations c
		WHERE (c.coach_id = $1 OR c.client_id = $1)
	`

	if !includeArchived {
		baseQuery += " AND c.is_archived = false"
		countQuery += " AND c.is_archived = false"
	}

	baseQuery += " ORDER BY c.last_message_at DESC NULLS LAST, c.created_at DESC"

	var total int
	if err := s.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	params := []interface{}{userID}
	limitPlaceholder := len(params) + 1
	offsetPlaceholder := limitPlaceholder + 1
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPlaceholder, offsetPlaceholder)
	params = append(params, limit, offset)

	rows, err := s.db.Query(ctx, baseQuery, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var conversations []types.ConversationOverview
	for rows.Next() {
		var conv types.ConversationOverview
		if err := rows.Scan(
			&conv.ConversationID,
			&conv.CoachID,
			&conv.ClientID,
			&conv.CoachName,
			&conv.CoachImage,
			&conv.ClientName,
			&conv.ClientImage,
			&conv.CreatedAt,
			&conv.LastMessageAt,
			&conv.IsArchived,
			&conv.LastMessageText,
			&conv.LastMessageSenderID,
			&conv.LastMessageSentAt,
			&conv.TotalMessages,
		); err != nil {
			return nil, 0, err
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

func (s *Store) IsParticipant(ctx context.Context, conversationID int, userID string) (bool, error) {
	q := `
		SELECT EXISTS(
			SELECT 1 FROM conversations 
			WHERE conversation_id = $1 AND (coach_id = $2 OR client_id = $2)
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, q, conversationID, userID).Scan(&exists)
	return exists, err
}
