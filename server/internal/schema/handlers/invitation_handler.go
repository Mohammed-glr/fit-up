package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	service "github.com/tdmdh/fit-up-server/internal/schema/services"
)

// InvitationHandler handles invitation-related HTTP requests
type InvitationHandler struct {
	invitationService service.InvitationService
}

// NewInvitationHandler creates a new invitation handler
func NewInvitationHandler(invitationService service.InvitationService) *InvitationHandler {
	return &InvitationHandler{
		invitationService: invitationService,
	}
}

// CreateInvitation handles POST /api/v1/coach/invitations
func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get coach ID from context (set by auth middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Printf("[CreateInvitation] No user_id in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req service.CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[CreateInvitation] Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set coach ID from authenticated user
	req.CoachID = userID

	// Validate email
	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	invitation, err := h.invitationService.CreateInvitation(ctx, &req)
	if err != nil {
		log.Printf("[CreateInvitation] Failed to create invitation: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create invitation")
		return
	}

	log.Printf("[CreateInvitation] Successfully created invitation %s for %s", invitation.ID, invitation.Email)
	respondWithJSON(w, http.StatusCreated, invitation)
}

// GetInvitations handles GET /api/v1/coach/invitations
func (h *InvitationHandler) GetInvitations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get coach ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Printf("[GetInvitations] No user_id in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	invitations, err := h.invitationService.GetInvitations(ctx, userID)
	if err != nil {
		log.Printf("[GetInvitations] Failed to get invitations: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get invitations")
		return
	}

	log.Printf("[GetInvitations] Retrieved %d invitations for coach %s", len(invitations), userID)
	respondWithJSON(w, http.StatusOK, invitations)
}

// ResendInvitation handles POST /api/v1/coach/invitations/:id/resend
func (h *InvitationHandler) ResendInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get coach ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Printf("[ResendInvitation] No user_id in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	invitationID := chi.URLParam(r, "id")
	if invitationID == "" {
		respondWithError(w, http.StatusBadRequest, "Invitation ID is required")
		return
	}

	invitation, err := h.invitationService.ResendInvitation(ctx, invitationID)
	if err != nil {
		log.Printf("[ResendInvitation] Failed to resend invitation: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to resend invitation")
		return
	}

	log.Printf("[ResendInvitation] Successfully resent invitation %s by coach %s", invitationID, userID)
	respondWithJSON(w, http.StatusOK, invitation)
}

// CancelInvitation handles DELETE /api/v1/coach/invitations/:id
func (h *InvitationHandler) CancelInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get coach ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Printf("[CancelInvitation] No user_id in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	invitationID := chi.URLParam(r, "id")
	if invitationID == "" {
		respondWithError(w, http.StatusBadRequest, "Invitation ID is required")
		return
	}

	err := h.invitationService.CancelInvitation(ctx, invitationID)
	if err != nil {
		log.Printf("[CancelInvitation] Failed to cancel invitation: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel invitation")
		return
	}

	log.Printf("[CancelInvitation] Successfully cancelled invitation %s by coach %s", invitationID, userID)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Invitation cancelled successfully"})
}

// AcceptInvitation handles POST /api/v1/invitations/accept
func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Printf("[AcceptInvitation] No user_id in context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[AcceptInvitation] Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		respondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	err := h.invitationService.AcceptInvitation(ctx, req.Token, userID)
	if err != nil {
		log.Printf("[AcceptInvitation] Failed to accept invitation: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to accept invitation")
		return
	}

	log.Printf("[AcceptInvitation] User %s successfully accepted invitation", userID)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Invitation accepted successfully"})
}
