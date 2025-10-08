package handlers

import (
	"github.com/tdmdh/fit-up-server/internal/message/services"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type MessageHandler struct {
	authMiddleware *middleware.AuthMiddleware
	service        services.MessageServiceManager
}

func NewMessageHandler(
	service services.MessageServiceManager,
	authMiddleware *middleware.AuthMiddleware,
) *MessageHandler {
	return &MessageHandler{
		authMiddleware: authMiddleware,
		service:        service,
	}
}



