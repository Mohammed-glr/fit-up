package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tdmdh/fit-up-server/internal/message/services"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

type MessageHandler struct {
	authMiddleware  *middleware.AuthMiddleware
	service         services.MessageServiceManager
	realtimeService *services.RealtimeService
}

func NewMessageHandler(
	service services.MessageServiceManager,
	authMiddleware *middleware.AuthMiddleware,
) *MessageHandler {
	return &MessageHandler{
		authMiddleware:  authMiddleware,
		service:         service,
		realtimeService: service.Realtime(),
	}
}

func SetupMessageRoutes(
	r chi.Router,
	messageHandler *MessageHandler,
	conversationHandler *ConversationHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireJWTAuth())

		r.Route("/conversations", func(r chi.Router) {
			r.Post("/", conversationHandler.CreateConversation)
			r.Get("/", conversationHandler.ListConversations)

			r.Route("/{conversation_id}", func(r chi.Router) {
				r.Get("/", conversationHandler.GetConversation)
				r.Get("/unread-count", conversationHandler.GetUnreadCount)
				r.Get("/messages", messageHandler.GetMessages)
				r.Post("/messages/read-all", messageHandler.MarkAllAsRead)
			})
		})

		r.Route("/messages", func(r chi.Router) {
			r.Post("/", messageHandler.SendMessage)

			r.Route("/{message_id}", func(r chi.Router) {
				r.Put("/", messageHandler.UpdateMessage)
				r.Delete("/", messageHandler.DeleteMessage)
				r.Post("/read", messageHandler.MarkMessageAsRead)
			})
		})
	})
}

func SetupMessageRoutesAlt(
	r chi.Router,
	messageHandler *MessageHandler,
	conversationHandler *ConversationHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	r.Route("/api/v1/messages", func(r chi.Router) {
		r.Use(authMiddleware.RequireJWTAuth())

		r.Post("/conversations", conversationHandler.CreateConversation)
		r.Get("/conversations", conversationHandler.ListConversations)
		r.Get("/conversations/{conversation_id}", conversationHandler.GetConversation)
		r.Get("/conversations/{conversation_id}/unread-count", conversationHandler.GetUnreadCount)

		r.Post("/send", messageHandler.SendMessage)
		r.Get("/conversations/{conversation_id}/messages", messageHandler.GetMessages)
		r.Put("/{message_id}", messageHandler.UpdateMessage)
		r.Delete("/{message_id}", messageHandler.DeleteMessage)
		r.Post("/{message_id}/read", messageHandler.MarkMessageAsRead)
		r.Post("/conversations/{conversation_id}/read-all", messageHandler.MarkAllAsRead)
	})
}
