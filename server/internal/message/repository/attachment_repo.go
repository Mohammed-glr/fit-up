package repository

import (
	"context"
	"fmt"

	"github.com/tdmdh/fit-up-server/internal/message/types"
)

func (s *Store) CreateAttachment(ctx context.Context, messageID int64, attachmentType types.AttachmentType, fileName, fileURL string) (*types.MessageAttachment, error) {
	q := `
		INSERT INTO message_attachments (message_id, attachment_type, file_name, file_url)
		VALUES ($1, $2, $3, $4)
		RETURNING attachment_id, message_id, attachment_type, file_name, file_url, file_size, mime_type, uploaded_at
	`
	var attachment types.MessageAttachment
	if err := s.db.QueryRow(ctx, q, messageID, attachmentType, fileName, fileURL).Scan(
		&attachment.AttachmentID,
		&attachment.MessageID,
		&attachment.AttachmentType,
		&attachment.FileName,
		&attachment.FileURL,
		&attachment.FileSize,
		&attachment.MimeType,
		&attachment.UploadedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}
	return &attachment, nil
}

func (s *Store) ListAttachmentsByMessage(ctx context.Context, messageID int64) ([]types.MessageAttachment, error) {
	q := `
		SELECT attachment_id, message_id, attachment_type, file_name, file_url, file_size, mime_type, uploaded_at, metadata
		FROM message_attachments
		WHERE message_id = $1
		ORDER BY uploaded_at ASC
	`
	rows, err := s.db.Query(ctx, q, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	defer rows.Close()

	var attachments []types.MessageAttachment
	for rows.Next() {
		var attachment types.MessageAttachment
		if err := rows.Scan(
			&attachment.AttachmentID,
			&attachment.MessageID,
			&attachment.AttachmentType,
			&attachment.FileName,
			&attachment.FileURL,
			&attachment.FileSize,
			&attachment.MimeType,
			&attachment.UploadedAt,
			&attachment.Metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

func (s *Store) DeleteAttachment(ctx context.Context, attachmentID int64) error {
	q := `DELETE FROM message_attachments WHERE attachment_id = $1`
	_, err := s.db.Exec(ctx, q, attachmentID)
	return err
}
