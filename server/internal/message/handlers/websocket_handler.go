package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	authService "github.com/tdmdh/fit-up-server/internal/auth/services"
	authTypes "github.com/tdmdh/fit-up-server/internal/auth/types"
	"github.com/tdmdh/fit-up-server/internal/message/services"
	"github.com/tdmdh/fit-up-server/internal/message/types"
	"github.com/tdmdh/fit-up-server/shared/config"
	"github.com/tdmdh/fit-up-server/shared/middleware"
	"golang.org/x/net/websocket"
)

type WebSocketHandler struct {
	realtimeService *services.RealtimeService
	authMiddleware  *middleware.AuthMiddleware
}

func NewWebSocketHandler(
	realtimeService *services.RealtimeService,
	authMiddleware *middleware.AuthMiddleware,
) *WebSocketHandler {
	return &WebSocketHandler{
		realtimeService: realtimeService,
		authMiddleware:  authMiddleware,
	}
}
func (wsh *WebSocketHandler) HandleWebSocketUpgrade() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		ctx := context.Background()

		token := wsh.extractToken(ws.Request())
		if token == "" {
			ws.Close()
			return
		}

		claims, err := wsh.validateToken(token)
		if err != nil {
			wsh.sendError(ws, "Authentication failed")
			ws.Close()
			return
		}

		userID := claims.UserID

		if err := wsh.realtimeService.HandleConnection(ctx, userID, ws); err != nil {
			wsh.sendError(ws, "Connection failed")
			ws.Close()
			return
		}

	})
}

func (wsh *WebSocketHandler) extractToken(r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

func (wsh *WebSocketHandler) sendError(ws *websocket.Conn, errorMsg string) {
	errMsg := types.WebSocketMessage{
		Type:      types.WSTypeError,
		Error:     &errorMsg,
		Timestamp: time.Now(),
	}
	websocket.JSON.Send(ws, errMsg)
}

func (wsh *WebSocketHandler) validateToken(tokenString string) (*authTypes.TokenClaims, error) {
	cfg := config.NewConfig()
	secret := []byte(cfg.JWTSecret)

	if len(secret) == 0 {
		return nil, fmt.Errorf("JWT secret not configured")
	}

	claims, err := authService.ValidateJWT(tokenString, nil, secret)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func (wsh *WebSocketHandler) HandleWebSocketHealthCheck(w http.ResponseWriter, r *http.Request) {
	activeConnections := 0
	if wsh.realtimeService != nil && wsh.realtimeService.Hub != nil {
		activeConnections = wsh.realtimeService.Hub.GetActiveConnections()
	}

	response := map[string]interface{}{
		"status":              "healthy",
		"active_connections":  activeConnections,
		"websocket_available": true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (wsh *WebSocketHandler) HandleWebSocketMetrics(w http.ResponseWriter, r *http.Request) {
	if wsh.realtimeService == nil || wsh.realtimeService.Hub == nil {
		http.Error(w, "WebSocket service not available", http.StatusServiceUnavailable)
		return
	}

	metrics := map[string]interface{}{
		"active_connections": wsh.realtimeService.Hub.GetActiveConnections(),
		"connected_users":    len(wsh.realtimeService.Hub.GetConnectedUsers()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

func (wsh *WebSocketHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	var req struct {
		ConversationID int `json:"conversation_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	channel := fmt.Sprintf("conversation:%d", req.ConversationID)
	wsh.realtimeService.Hub.Subscribe(userID, channel)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "subscribed"})
}

func (wsh *WebSocketHandler) HandleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	var req struct {
		ConversationID int `json:"conversation_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	channel := fmt.Sprintf("conversation:%d", req.ConversationID)
	wsh.realtimeService.Hub.Unsubscribe(userID, channel)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "unsubscribed"})
}
