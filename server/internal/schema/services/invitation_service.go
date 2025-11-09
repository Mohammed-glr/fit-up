package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/repository"
)

type invitationService struct {
	repo repository.CoachInvitationRepo
	// emailService *email.Service // TODO: Add email service when implemented
}
func NewInvitationService(repo repository.CoachInvitationRepo) InvitationService {
	return &invitationService{
		repo: repo,
	}
}

type CreateInvitationRequest struct {
	CoachID       string  `json:"coach_id"`
	Email         string  `json:"email"`
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	CustomMessage *string `json:"custom_message"`
}

type InvitationResponse struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	FirstName       *string    `json:"first_name"`
	LastName        *string    `json:"last_name"`
	Status          string     `json:"status"`
	CustomMessage   *string    `json:"custom_message"`
	ExpiresAt       time.Time  `json:"expires_at"`
	CreatedAt       time.Time  `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	InvitationToken string     `json:"invitation_token,omitempty"`
}

func generateInvitationToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (s *invitationService) CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*InvitationResponse, error) {
	existing, err := s.repo.GetInvitationByCoachAndEmail(ctx, req.CoachID, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing invitation: %w", err)
	}

	if existing != nil {
		if existing.ExpiresAt.After(time.Now()) {
			return &InvitationResponse{
				ID:            existing.ID,
				Email:         existing.Email,
				FirstName:     existing.FirstName,
				LastName:      existing.LastName,
				Status:        existing.Status,
				CustomMessage: existing.CustomMessage,
				ExpiresAt:     existing.ExpiresAt,
				CreatedAt:     existing.CreatedAt,
				AcceptedAt:    existing.AcceptedAt,
			}, nil
		}

		err = s.repo.UpdateInvitationStatus(ctx, existing.ID, "expired")
		if err != nil {
			log.Printf("[CreateInvitation] Failed to expire old invitation: %v", err)
		}
	}

	token, err := generateInvitationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}

	invitation := &repository.CoachInvitation{
		ID:              generateID(),
		CoachID:         req.CoachID,
		Email:           req.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		InvitationToken: token,
		Status:          "pending",
		CustomMessage:   req.CustomMessage,
		ExpiresAt:       time.Now().Add(7 * 24 * time.Hour), // 7 days expiration
	}

	err = s.repo.CreateInvitation(ctx, invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	// TODO: Send invitation email
	log.Printf("[CreateInvitation] Created invitation for %s (token: %s)", req.Email, token)

	return &InvitationResponse{
		ID:              invitation.ID,
		Email:           invitation.Email,
		FirstName:       invitation.FirstName,
		LastName:        invitation.LastName,
		Status:          invitation.Status,
		CustomMessage:   invitation.CustomMessage,
		ExpiresAt:       invitation.ExpiresAt,
		CreatedAt:       invitation.CreatedAt,
		InvitationToken: token, // Include token in response for testing
	}, nil
}

// GetInvitations retrieves all invitations for a coach
func (s *invitationService) GetInvitations(ctx context.Context, coachID string) ([]*InvitationResponse, error) {
	// Expire old invitations first
	_, err := s.repo.ExpireOldInvitations(ctx)
	if err != nil {
		log.Printf("[GetInvitations] Failed to expire old invitations: %v", err)
	}

	invitations, err := s.repo.GetInvitationsByCoachID(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	responses := make([]*InvitationResponse, len(invitations))
	for i, inv := range invitations {
		responses[i] = &InvitationResponse{
			ID:            inv.ID,
			Email:         inv.Email,
			FirstName:     inv.FirstName,
			LastName:      inv.LastName,
			Status:        inv.Status,
			CustomMessage: inv.CustomMessage,
			ExpiresAt:     inv.ExpiresAt,
			CreatedAt:     inv.CreatedAt,
			AcceptedAt:    inv.AcceptedAt,
		}
	}

	return responses, nil
}

// ResendInvitation resends an invitation (creates a new token)
func (s *invitationService) ResendInvitation(ctx context.Context, invitationID string) (*InvitationResponse, error) {
	// Get the invitation
	invitations, err := s.repo.GetInvitationsByCoachID(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	var invitation *repository.CoachInvitation
	for _, inv := range invitations {
		if inv.ID == invitationID {
			invitation = inv
			break
		}
	}

	if invitation == nil {
		return nil, fmt.Errorf("invitation not found")
	}

	if invitation.Status != "pending" {
		return nil, fmt.Errorf("cannot resend non-pending invitation")
	}

	// Generate new token
	token, err := generateInvitationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}

	// Create new invitation (mark old as cancelled)
	err = s.repo.UpdateInvitationStatus(ctx, invitationID, "cancelled")
	if err != nil {
		return nil, fmt.Errorf("failed to cancel old invitation: %w", err)
	}

	newInvitation := &repository.CoachInvitation{
		ID:              generateID(),
		CoachID:         invitation.CoachID,
		Email:           invitation.Email,
		FirstName:       invitation.FirstName,
		LastName:        invitation.LastName,
		InvitationToken: token,
		Status:          "pending",
		CustomMessage:   invitation.CustomMessage,
		ExpiresAt:       time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.repo.CreateInvitation(ctx, newInvitation)
	if err != nil {
		return nil, fmt.Errorf("failed to create new invitation: %w", err)
	}

	// TODO: Send invitation email
	log.Printf("[ResendInvitation] Resent invitation for %s (token: %s)", invitation.Email, token)

	return &InvitationResponse{
		ID:              newInvitation.ID,
		Email:           newInvitation.Email,
		FirstName:       newInvitation.FirstName,
		LastName:        newInvitation.LastName,
		Status:          newInvitation.Status,
		CustomMessage:   newInvitation.CustomMessage,
		ExpiresAt:       newInvitation.ExpiresAt,
		CreatedAt:       newInvitation.CreatedAt,
		InvitationToken: token,
	}, nil
}

// CancelInvitation cancels an invitation
func (s *invitationService) CancelInvitation(ctx context.Context, invitationID string) error {
	err := s.repo.DeleteInvitation(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("failed to cancel invitation: %w", err)
	}

	return nil
}

// AcceptInvitation accepts an invitation and assigns the client to the coach
func (s *invitationService) AcceptInvitation(ctx context.Context, token, userID string) error {
	// Get invitation by token
	invitation, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get invitation: %w", err)
	}

	// Validate invitation
	if invitation.Status != "pending" {
		return fmt.Errorf("invitation is not pending")
	}

	if invitation.ExpiresAt.Before(time.Now()) {
		// Mark as expired
		s.repo.UpdateInvitationStatus(ctx, invitation.ID, "expired")
		return fmt.Errorf("invitation has expired")
	}

	// Accept the invitation
	err = s.repo.AcceptInvitation(ctx, invitation.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to accept invitation: %w", err)
	}

	// TODO: Assign client to coach using coach service
	log.Printf("[AcceptInvitation] User %s accepted invitation from coach %s", userID, invitation.CoachID)

	return nil
}
