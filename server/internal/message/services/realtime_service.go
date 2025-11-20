package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tdmdh/fit-up-server/internal/message/pool"
	"github.com/tdmdh/fit-up-server/internal/message/types"
	"golang.org/x/net/websocket"
)

type RealtimeService struct {
	Hub             *pool.Hub
	messageService  MessageService
	conversationSvc ConversationService
	readStatusSvc   MessageReadStatusService
}

func NewRealtimeService(
	hub *pool.Hub,
	messageService MessageService,
	conversationSvc ConversationService,
	readStatusSvc MessageReadStatusService,
) *RealtimeService {
	return &RealtimeService{
		Hub:             hub,
		messageService:  messageService,
		conversationSvc: conversationSvc,
		readStatusSvc:   readStatusSvc,
	}
}

func (rs *RealtimeService) HandleConnection(ctx context.Context, userID string, conn *websocket.Conn) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	if conn == nil {
		return fmt.Errorf("websocket connection cannot be nil")
	}

	rs.Hub.Connect(userID, conn)

	log.Printf("User %s connected to WebSocket", userID)

	rs.waitForDisconnect(userID)

	return nil
}

func (rs *RealtimeService) waitForDisconnect(userID string) {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !rs.Hub.IsConnected(userID) {
			return
		}
	}
}

func (rs *RealtimeService) BroadcastNewMessage(ctx context.Context, conversationID int, message *types.MessageWithDetails) error {
	wsMessage := types.WebSocketMessage{
		Type:           types.WSTypeNewMessage,
		ConversationID: conversationID,
		Message:        message,
		Timestamp:      time.Now(),
	}

	return rs.broadcastToConversation(conversationID, wsMessage)
}

func (rs *RealtimeService) BroadcastMessageEdited(ctx context.Context, conversationID int, message *types.MessageWithDetails) error {
	wsMessage := types.WebSocketMessage{
		Type:           types.WSTypeMessageEdited,
		ConversationID: conversationID,
		Message:        message,
		Timestamp:      time.Now(),
	}

	return rs.broadcastToConversation(conversationID, wsMessage)
}

func (rs *RealtimeService) BroadcastMessageDeleted(ctx context.Context, conversationID int, messageID int64) error {
	wsMessage := types.WebSocketMessage{
		Type:           types.WSTypeMessageDeleted,
		ConversationID: conversationID,
		MessageID:      &messageID,
		Timestamp:      time.Now(),
	}

	return rs.broadcastToConversation(conversationID, wsMessage)
}

func (rs *RealtimeService) BroadcastMessageRead(ctx context.Context, conversationID int, messageID int64, userID string) error {
	wsMessage := types.WebSocketMessage{
		Type:           types.WSTypeMessageRead,
		ConversationID: conversationID,
		MessageID:      &messageID,
		ReadBy:         &userID,
		Timestamp:      time.Now(),
	}

	return rs.broadcastToConversation(conversationID, wsMessage)
}

func (rs *RealtimeService) broadcastToConversation(conversationID int, message types.WebSocketMessage) error {
	channel := fmt.Sprintf("conversation:%d", conversationID)

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	rs.Hub.BroadcastToChannel(channel, string(messageBytes))
	return nil
}

func (rs *RealtimeService) SendToUser(userID string, message types.WebSocketMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return rs.Hub.SendMessage(userID, string(messageBytes))
}

func (rs *RealtimeService) sendErrorToUser(userID string, errorMsg string) {
	errMessage := types.WebSocketMessage{
		Type:      types.WSTypeError,
		Error:     &errorMsg,
		Timestamp: time.Now(),
	}

	if err := rs.SendToUser(userID, errMessage); err != nil {
		log.Printf("Failed to send error to user %s: %v", userID, err)
	}
}

func (rs *RealtimeService) SubscribeToConversation(ctx context.Context, userID string, conversationID int) error {
	isParticipant, err := rs.conversationSvc.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify participant: %w", err)
	}

	if !isParticipant {
		return fmt.Errorf("user is not a participant in this conversation")
	}

	channel := fmt.Sprintf("conversation:%d", conversationID)
	rs.Hub.Subscribe(userID, channel)

	return nil
}

func (rs *RealtimeService) UnsubscribeFromConversation(userID string, conversationID int) {
	channel := fmt.Sprintf("conversation:%d", conversationID)
	rs.Hub.Unsubscribe(userID, channel)
}

func (rs *RealtimeService) GetActiveConnections() int {
	return rs.Hub.GetActiveConnections()
}

func (rs *RealtimeService) GetConnectedUsers() []string {
	return rs.Hub.GetConnectedUsers()
}

func (rs *RealtimeService) IsUserConnected(userID string) bool {
	return rs.Hub.IsConnected(userID)
}

func (rs *RealtimeService) DisconnectUser(userID string) {
	rs.Hub.Disconnect(userID)
}

func (rs *RealtimeService) NotifyNewMessageToParticipants(ctx context.Context, message *types.MessageWithDetails, coachID, clientID string) error {
	channel := fmt.Sprintf("conversation:%d", message.ConversationID)

	rs.Hub.Subscribe(coachID, channel)
	rs.Hub.Subscribe(clientID, channel)
	return rs.BroadcastNewMessage(ctx, message.ConversationID, message)
}
