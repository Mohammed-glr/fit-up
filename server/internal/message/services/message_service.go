package services

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/types"
)

type messageService struct {
	repo           repository.MessageRepo
	attachmentRepo repository.MessageAttachmentRepo
}

func NewMessageService(repo repository.MessageStore) MessageService {
	return &messageService{
		repo:           repo.Messages(),
		attachmentRepo: repo.Attachments(),
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

func (s *messageService) ListMessages(ctx context.Context, conversationID int, userID string, limit, offset int) (*types.MessagesResponse, error) {
	if conversationID <= 0 {
		return nil, types.ErrInvalidConversationID
	}
	if userID == "" {
		return nil, types.ErrInvalidUserID
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	messages, total, err := s.repo.ListMessages(ctx, conversationID, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		attachments, err := s.attachmentRepo.ListAttachmentsByMessage(ctx, messages[i].MessageID)
		if err != nil {
			return nil, err
		}
		if len(attachments) > 0 {
			messages[i].Attachments = attachments
		}
	}

	hasMore := offset+len(messages) < total

	return &types.MessagesResponse{
		Messages: messages,
		Total:    total,
		HasMore:  hasMore,
	}, nil
}
