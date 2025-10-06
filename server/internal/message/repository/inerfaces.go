package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/types"
)

type ConversationRepo interface {
	CreateConversation(ctx context.Context, coachID, clientID string) (*types.Conversation, error)
	GetConversationByID(ctx context.Context, conversationID int) (*types.Conversation, error)
	GetConversationByParticipants(ctx context.Context, coachID, clientID string) (*types.Conversation, error)
	ListConversationsByUser(ctx context.Context, userID string, includeArchived bool) ([]types.ConversationOverview, error)

	IsParticipant(ctx context.Context, conversationID int, userID string) (bool, error)
}

type MessageRepo interface {
	CreateMessage(ctx context.Context, conversationID int, senderID, messageText string, replyToMessageID *int64) (*types.Message, error)
	GetMessageByID(ctx context.Context, messageID int64) (*types.Message, error)

	UpdateMessage(ctx context.Context, messageID int64, messageText string) error
	DeleteMessage(ctx context.Context, messageID int64) error
}

type MessageReadStatusRepo interface {
	MarkMessageAsRead(ctx context.Context, messageID int64, userID string) error
	MarkAllAsRead(ctx context.Context, conversationID int, userID string) error
	CountUnreadMessages(ctx context.Context, conversationID int, userID string) (int, error)
}

type MessageAttachmentRepo interface {
	CreateAttachment(ctx context.Context, messageID int64, attachmentType types.AttachmentType, fileName, fileURL string) (*types.MessageAttachment, error)
	ListAttachmentsByMessage(ctx context.Context, messageID int64) ([]types.MessageAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID int64) error
}

type MessageStore interface {
	Conversations() ConversationRepo
	Messages() MessageRepo
	ReadStatus() MessageReadStatusRepo
	Attachments() MessageAttachmentRepo
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
