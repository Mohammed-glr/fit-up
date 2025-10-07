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

	if coachID > clientID {
		coachID, clientID = clientID, coachID
	}

	return s.repo.GetConversationByParticipants(ctx, coachID, clientID)
}

func (s *conversationService) ListConversationsByUser(ctx context.Context, userID string, includeArchived bool) ([]types.ConversationOverview, error) {
	if userID == "" {
		return nil, types.ErrInvalidUserID
	}

	userRole := middleware.GetUserRoleFromContext(ctx)
	if userRole != "coach" && userRole != "client" {
		return nil, types.ErrUnauthorized
	}

	return s.repo.ListConversationsByUser(ctx, userID, includeArchived)
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
