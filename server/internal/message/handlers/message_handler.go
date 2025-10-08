package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/message/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	var req types.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	isParticipant, err := h.service.Conversations().IsParticipant(ctx, req.ConversationID, userID)
	if err != nil {
		log.Printf("Error checking participant: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to verify permissions")
		return
	}

	if !isParticipant {
		respondError(w, http.StatusForbidden, "You are not a participant in this conversation")
		return
	}

	message, err := h.service.Messages().CreateMessage(
		ctx,
		req.ConversationID,
		userID,
		req.MessageText,
		req.ReplyToMessageID,
	)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create message")
		return
	}

	messageWithDetails := &types.MessageWithDetails{
		Message: *message,
		SenderName: userID,
		IsRead:     false,
	}

	if h.realtimeService != nil {
		if err := h.realtimeService.BroadcastNewMessage(ctx, req.ConversationID, messageWithDetails); err != nil {
			log.Printf("Failed to broadcast message: %v", err)
		}
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": messageWithDetails,
	})
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
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

	if h.realtimeService != nil {
		if err := h.realtimeService.SubscribeToConversation(ctx, userID, conversationID); err != nil {
			log.Printf("Failed to subscribe to conversation: %v", err)
		}
	}

	// TODO: Implement message fetching with pagination
	// For now, return empty array
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"messages": []types.MessageWithDetails{},
		"total":    0,
		"has_more": false,
	})
}

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	messageIDStr := chi.URLParam(r, "message_id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	var req types.UpdateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	message, err := h.service.Messages().GetMessageByID(ctx, messageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	if message.SenderID != userID {
		respondError(w, http.StatusForbidden, "You can only edit your own messages")
		return
	}

	if err := h.service.Messages().UpdateMessage(ctx, messageID, req.MessageText); err != nil {
		log.Printf("Error updating message: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update message")
		return
	}

	updatedMessage, _ := h.service.Messages().GetMessageByID(ctx, messageID)
	messageWithDetails := &types.MessageWithDetails{
		Message:    *updatedMessage,
		SenderName: userID,
	}

	if h.realtimeService != nil {
		if err := h.realtimeService.BroadcastMessageEdited(ctx, message.ConversationID, messageWithDetails); err != nil {
			log.Printf("Failed to broadcast message update: %v", err)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": messageWithDetails,
	})
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	messageIDStr := chi.URLParam(r, "message_id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	message, err := h.service.Messages().GetMessageByID(ctx, messageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	if message.SenderID != userID {
		respondError(w, http.StatusForbidden, "You can only delete your own messages")
		return
	}

	if err := h.service.Messages().DeleteMessage(ctx, messageID); err != nil {
		log.Printf("Error deleting message: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete message")
		return
	}

	if h.realtimeService != nil {
		if err := h.realtimeService.BroadcastMessageDeleted(ctx, message.ConversationID, messageID); err != nil {
			log.Printf("Failed to broadcast message deletion: %v", err)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Message deleted successfully",
	})
}

func (h *MessageHandler) MarkMessageAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetAuthIDFromContext(ctx)

	messageIDStr := chi.URLParam(r, "message_id")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	message, err := h.service.Messages().GetMessageByID(ctx, messageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		respondError(w, http.StatusNotFound, "Message not found")
		return
	}

	isParticipant, err := h.service.Conversations().IsParticipant(ctx, message.ConversationID, userID)
	if err != nil || !isParticipant {
		respondError(w, http.StatusForbidden, "You are not a participant in this conversation")
		return
	}

	if err := h.service.ReadStatus().MarkMessageAsRead(ctx, messageID, userID); err != nil {
		log.Printf("Error marking message as read: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to mark message as read")
		return
	}

	if h.realtimeService != nil {
		if err := h.realtimeService.BroadcastMessageRead(ctx, message.ConversationID, messageID, userID); err != nil {
			log.Printf("Failed to broadcast read status: %v", err)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Message marked as read",
	})
}

func (h *MessageHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.ReadStatus().MarkAllAsRead(ctx, conversationID, userID); err != nil {
		log.Printf("Error marking all as read: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to mark messages as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "All messages marked as read",
	})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error": message,
	})
}
