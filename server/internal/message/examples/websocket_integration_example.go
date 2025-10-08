package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/tdmdh/fit-up-server/internal/message/handlers"
	"github.com/tdmdh/fit-up-server/internal/message/pool"
	"github.com/tdmdh/fit-up-server/internal/message/repository"
	"github.com/tdmdh/fit-up-server/internal/message/services"
	"github.com/tdmdh/fit-up-server/internal/message/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
)

// This is an example of how to integrate the WebSocket layer with your message service
// Add this integration to your main server initialization

func setupMessageServiceWithWebSockets() {
	// Initialize database connection (adjust based on your setup)
	// db := setupDatabase()
	// messageRepo := repository.NewMessageStore(db)

	// For demonstration, we'll show the structure:
	var messageRepo repository.MessageStore       // Your actual implementation
	var authMiddleware *middleware.AuthMiddleware // Your auth middleware

	// Initialize the WebSocket Hub
	hub := pool.NewHub()

	// Start the hub in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	// Initialize message services
	messageService := services.NewMessagesService(messageRepo)

	// Initialize realtime service
	realtimeService := services.NewRealtimeService(
		hub,
		messageService.Messages(),
		messageService.Conversations(),
		messageService.ReadStatus(),
	)

	// Set the realtime service in the message service manager
	messageService.SetRealtimeService(realtimeService)

	// Initialize handlers
	// messageHandler := handlers.NewMessageHandler(messageService, authMiddleware)
	wsHandler := handlers.NewWebSocketHandler(realtimeService, authMiddleware)

	// Setup router
	r := chi.NewRouter()

	// Add middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(chiMiddleware.RealIP)

	// Setup WebSocket routes
	handlers.SetupWebSocketRoutes(r, wsHandler)

	// Setup other message service routes
	// r.Route("/api/v1/messages", func(r chi.Router) {
	// 	r.Use(authMiddleware.RequireJWTAuth())
	// 	// Add your message routes here
	// })

	// Start HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Println("Starting message service with WebSocket support on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	// Shutdown the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Stop the hub
	hub.Stop()
	cancel()

	log.Println("Server stopped gracefully")
}

// Integration example for broadcasting messages when they're created
func exampleBroadcastNewMessage(
	ctx context.Context,
	messageService services.MessageService,
	realtimeService *services.RealtimeService,
	conversationID int,
	senderID string,
	messageText string,
) error {
	// Create the message
	message, err := messageService.CreateMessage(ctx, conversationID, senderID, messageText, nil)
	if err != nil {
		return err
	}

	// Build the full message details (you'd fetch this from your repository)
	messageWithDetails := &types.MessageWithDetails{
		Message: *message,
		// Add other fields like SenderName, Attachments, etc.
	}

	// Broadcast to all participants in the conversation
	if err := realtimeService.BroadcastNewMessage(ctx, conversationID, messageWithDetails); err != nil {
		log.Printf("Failed to broadcast message: %v", err)
		// Don't fail the message creation, just log the error
	}

	return nil
}

// Example: Subscribe user to conversation when they open it
func exampleSubscribeToConversation(
	ctx context.Context,
	realtimeService *services.RealtimeService,
	userID string,
	conversationID int,
) error {
	return realtimeService.SubscribeToConversation(ctx, userID, conversationID)
}
