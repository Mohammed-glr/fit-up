package services

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/types"
)
type messageService struct {
	repo repository.MessageRepo
}

func NewMessageService(repo repository.MessageStore) MessageService {
	return &messageService{
		repo: repo.Messages(),
	}
}


func (s *messageService) CreateMessage(ctx context.Context, conversationID int, senderID, messageText string, replyToMessageID *int64) (*types.Message, error) {
	if err := ValidateMessageText(messageText); err != nil {
		return nil, err
	}

	return s.repo.CreateMessage(ctx, conversationID, senderID, messageText, replyToMessageID)
}


func ValidateMessageText(messageText string) error {
	if len(messageText) == 0 {
		return types.ErrMessageEmpty
	}
	if len(messageText) > 5000 {
		return types.ErrMessageTooLong
	}
	return nil
}

func (s *messageService) GetMessageByID(ctx context.Context, messageID int64) (*types.Message, error) {
	return s.repo.GetMessageByID(ctx, messageID)
}

func (s *messageService) UpdateMessage(ctx context.Context, messageID int64, messageText string) error {
	if err := ValidateMessageText(messageText); err != nil {
		return err
	}

	return s.repo.UpdateMessage(ctx, messageID, messageText)
}


func (s *messageService) DeleteMessage(ctx context.Context, messageID int64) error {
	return s.repo.DeleteMessage(ctx, messageID)
}