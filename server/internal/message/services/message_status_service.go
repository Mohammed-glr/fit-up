package services

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/types"
)

type messageStatusService struct {
	repo repository.MessageReadStatusRepo
}

func NewMessageReadStatusService(repo repository.MessageReadStatusRepo) *messageStatusService {
	return &messageStatusService{
		repo: repo,
	}
}


func (s *messageStatusService) MarkMessageAsRead(ctx context.Context, messageID int64, userID string) error {
	if messageID <= 0 {
		return types.ErrInvalidMessageID
	}
	if userID == "" {
		return types.ErrInvalidUserID
	}

	return s.repo.MarkMessageAsRead(ctx, messageID, userID)
}

func (s *messageStatusService) MarkAllAsRead(ctx context.Context, conversationID int, userID string) error {
	if conversationID <= 0 {
		return types.ErrInvalidConversationID
	}
	if userID == "" {
		return types.ErrInvalidUserID
	}

	return s.repo.MarkAllAsRead(ctx, conversationID, userID)
}


func (s *messageStatusService) CountUnreadMessages(ctx context.Context, conversationID int, userID string) (int, error) {
	if conversationID <= 0 {
		return 0, types.ErrInvalidConversationID
	}

	if userID == "" {
		return 0, types.ErrInvalidUserID
	}

	return s.repo.CountUnreadMessages(ctx, conversationID, userID)
}
