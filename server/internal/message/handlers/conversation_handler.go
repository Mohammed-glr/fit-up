package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/message/services"
	"github.com/tdmdh/fit-up-server/internal/message/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type ConversationHandler struct {
	authMiddleware  *middleware.AuthMiddleware
	service         services.MessageServiceManager
	realtimeService *services.RealtimeService
}

func NewConversationHandler(
	service services.MessageServiceManager,
	authMiddleware *middleware.AuthMiddleware,
) *ConversationHandler {
	return &ConversationHandler{
		authMiddleware:  authMiddleware,
		service:         service,
		realtimeService: service.Realtime(),
	}
}

func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	var req types.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.CoachID != userID && req.ClientID != userID {
		respondError(w, http.StatusForbidden, "You must be a participant in the conversation")
		return
	}

	existingConv, err := h.service.Conversations().GetConversationByParticipants(ctx, req.CoachID, req.ClientID)
	if err == nil && existingConv != nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"conversation": existingConv,
			"message":      "Conversation already exists",
		})
		return
	}

	conversation, err := h.service.Conversations().CreateConversation(ctx, &req)
	if err != nil {
		log.Printf("Error creating conversation: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create conversation")
		return
	}

	if h.realtimeService != nil {
		h.realtimeService.SubscribeToConversation(ctx, req.CoachID, conversation.ConversationID)
		h.realtimeService.SubscribeToConversation(ctx, req.ClientID, conversation.ConversationID)
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"conversation": conversation,
	})
}

func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	conversationIDStr := chi.URLParam(r, "conversation_id")
	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	isParticipant, err := h.service.Conversations().IsParticipant(ctx, conversationID, userID)
	if err != nil {
		log.Printf("Error checking participant: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to verify permissions")
		return
	}

	if !isParticipant {
		respondError(w, http.StatusForbidden, "You are not a participant in this conversation")
		return
	}

	conversation, err := h.service.Conversations().GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Printf("Error fetching conversation: %v", err)
		respondError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	if h.realtimeService != nil {
		if err := h.realtimeService.SubscribeToConversation(ctx, userID, conversationID); err != nil {
			log.Printf("Failed to subscribe to conversation: %v", err)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"conversation": conversation,
	})
}

func (h *ConversationHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	includeArchived := r.URL.Query().Get("include_archived") == "true"

	conversations, err := h.service.Conversations().ListConversationsByUser(ctx, userID, includeArchived)
	if err != nil {
		log.Printf("Error listing conversations: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list conversations")
		return
	}

	if h.realtimeService != nil {
		for _, conv := range conversations {
			if !conv.IsArchived {
				h.realtimeService.SubscribeToConversation(ctx, userID, conv.ConversationID)
			}
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"conversations": conversations,
		"total":         len(conversations),
	})
}

func (h *ConversationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	conversationIDStr := chi.URLParam(r, "conversation_id")
	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	isParticipant, err := h.service.Conversations().IsParticipant(ctx, conversationID, userID)
	if err != nil || !isParticipant {
		respondError(w, http.StatusForbidden, "You are not a participant in this conversation")
		return
	}
	count, err := h.service.ReadStatus().CountUnreadMessages(ctx, conversationID, userID)
	if err != nil {
		log.Printf("Error counting unread messages: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to count unread messages")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"conversation_id": conversationID,
		"unread_count":    count,
	})
}
