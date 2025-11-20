package services

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type conversationService struct {
	repo repository.ConversationRepo
}

func NewConversationService(repo repository.MessageStore) ConversationService {
	return &conversationService{
		repo: repo.Conversations(),
	}
}

func (s *conversationService) CreateConversation(ctx context.Context, req *types.CreateConversationRequest) (*types.Conversation, error) {
	if err := ValidateRequest(req); err != nil {
		return nil, err
	}

	existingConv, err := s.repo.GetConversationByParticipants(ctx, req.CoachID, req.ClientID)
	if err != nil {
		return nil, err
	}

	if existingConv != nil {
		return nil, types.ErrConversationExists
	}

	return s.repo.CreateConversation(ctx, req.CoachID, req.ClientID)
}

func (s *conversationService) GetConversationByID(ctx context.Context, conversationID int) (*types.Conversation, error) {
	return s.repo.GetConversationByID(ctx, conversationID)
}

func (s *conversationService) GetConversationByParticipants(ctx context.Context, coachID, clientID string) (*types.Conversation, error) {
	if err := ValidateParticipants(coachID, clientID); err != nil {
		return nil, err
	}

	return s.repo.GetConversationByParticipants(ctx, coachID, clientID)
}

func (s *conversationService) ListConversationsByUser(ctx context.Context, userID string, includeArchived bool, limit, offset int) (*types.ConversationsResponse, error) {
	if userID == "" {
		return nil, types.ErrInvalidUserID
	}

	// Allow both coaches and clients (users) to list their conversations
	userRole := middleware.GetUserRoleFromContext(ctx)
	if userRole != "coach" && userRole != "user" && userRole != "client" {
		return nil, types.ErrUnauthorized
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	conversations, total, err := s.repo.ListConversationsByUser(ctx, userID, includeArchived, limit, offset)
	if err != nil {
		return nil, err
	}

	hasMore := offset+len(conversations) < total

	return &types.ConversationsResponse{
		Conversations: conversations,
		Total:         total,
		HasMore:       hasMore,
	}, nil
}

func (s *conversationService) IsParticipant(ctx context.Context, conversationID int, userID string) (bool, error) {
	return s.repo.IsParticipant(ctx, conversationID, userID)
}

func ValidateParticipants(coachID, clientID string) error {
	if coachID == "" || clientID == "" {
		return types.ErrInvalidConversationParticipants
	}
	if coachID == clientID {
		return types.ErrInvalidConversationParticipants
	}
	return nil
}

func ValidateRequest(req *types.CreateConversationRequest) error {
	if req == nil {
		return types.ErrInvalidConversationParticipants
	}
	return ValidateParticipants(req.CoachID, req.ClientID)
}
