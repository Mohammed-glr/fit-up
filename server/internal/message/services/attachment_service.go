package services

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/types"
)

type attachmentService struct {
	repo repository.MessageAttachmentRepo
}

func NewMessageAttachmentService(repo repository.MessageStore) MessageAttachmentService {
	return &attachmentService{
		repo: repo.Attachments(),
	}
}


func (s *attachmentService) CreateAttachment(ctx context.Context, messageID int64, attachmentType types.AttachmentType, fileName, fileURL string) (*types.MessageAttachment, error) {
	if messageID <= 0 {
		return nil, types.ErrInvalidMessageID
	}
	if err := ValidateAttachmentType(attachmentType); err != nil {
		return nil, err
	}
	if fileName == "" {
		return nil, types.ErrInvalidFileName
	}
	if fileURL == "" {
		return nil, types.ErrInvalidFileURL
	}

	return s.repo.CreateAttachment(ctx, messageID, attachmentType, fileName, fileURL)
}

func ValidateAttachmentType(attachmentType types.AttachmentType) error {
	switch attachmentType {
	case types.AttachmentTypeImage, types.AttachmentTypeDocument:
		return nil
	default:
		return types.ErrInvalidAttachmentType
	}
}


func (s *attachmentService) ListAttachmentsByMessage(ctx context.Context, messageID int64) ([]types.MessageAttachment, error) {
	if messageID <= 0 {
		return nil, types.ErrInvalidMessageID
	}

	return s.repo.ListAttachmentsByMessage(ctx, messageID)
}

func (s *attachmentService) DeleteAttachment(ctx context.Context, attachmentID int64) error {
	if attachmentID <= 0 {
		return types.ErrInvalidMessageID
	}

	return s.repo.DeleteAttachment(ctx, attachmentID)
}
