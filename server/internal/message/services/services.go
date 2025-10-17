package services

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/types"
)

type ConversationService interface {
	CreateConversation(ctx context.Context, req *types.CreateConversationRequest) (*types.Conversation, error)
	GetConversationByID(ctx context.Context, conversationID int) (*types.Conversation, error)
	GetConversationByParticipants(ctx context.Context, coachID, clientID string) (*types.Conversation, error)
	ListConversationsByUser(ctx context.Context, userID string, includeArchived bool) ([]types.ConversationOverview, error)

	IsParticipant(ctx context.Context, conversationID int, userID string) (bool, error)
}

type MessageService interface {
	CreateMessage(ctx context.Context, conversationID int, senderID, messageText string, replyToMessageID *int64) (*types.Message, error)
	GetMessageByID(ctx context.Context, messageID int64) (*types.Message, error)
	ListMessages(ctx context.Context, conversationID int, userID string, limit, offset int) (*types.MessagesResponse, error)

	UpdateMessage(ctx context.Context, messageID int64, messageText string) error
	DeleteMessage(ctx context.Context, messageID int64) error
}

type MessageReadStatusService interface {
	MarkMessageAsRead(ctx context.Context, messageID int64, userID string) error
	MarkAllAsRead(ctx context.Context, conversationID int, userID string) error
	CountUnreadMessages(ctx context.Context, conversationID int, userID string) (int, error)
}

type MessageAttachmentService interface {
	CreateAttachment(ctx context.Context, messageID int64, attachmentType types.AttachmentType, fileName, fileURL string) (*types.MessageAttachment, error)
	ListAttachmentsByMessage(ctx context.Context, messageID int64) ([]types.MessageAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID int64) error
}

type MessageServiceManager interface {
	Conversations() ConversationService
	Messages() MessageService
	ReadStatus() MessageReadStatusService
	Attachments() MessageAttachmentService
	Realtime() *RealtimeService
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	repo repository.MessageStore

	conversationService      ConversationService
	realtimeService          *RealtimeService
	messageService           MessageService
	messageReadStatusService MessageReadStatusService
	messageAttachmentService MessageAttachmentService
}

func NewMessagesService(repo repository.MessageStore) *Service {
	return &Service{
		repo:                     repo,
		conversationService:      NewConversationService(repo),
		messageService:           NewMessageService(repo),
		messageReadStatusService: NewMessageReadStatusService(repo.ReadStatus()),
		messageAttachmentService: NewMessageAttachmentService(repo),
		realtimeService:          nil, // Will be set later with SetRealtimeService
	}
}

func (s *Service) SetRealtimeService(realtimeService *RealtimeService) {
	s.realtimeService = realtimeService
}

func (s *Service) Conversations() ConversationService {
	return s.conversationService
}

func (s *Service) Messages() MessageService {
	return s.messageService
}

func (s *Service) ReadStatus() MessageReadStatusService {
	return s.messageReadStatusService
}

func (s *Service) Attachments() MessageAttachmentService {
	return s.messageAttachmentService
}

func (s *Service) Realtime() *RealtimeService {
	return s.realtimeService
}

func (s *Service) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return s.repo.WithTransaction(ctx, fn)
}
