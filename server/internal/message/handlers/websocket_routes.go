package handlers

import (
	"github.com/go-chi/chi/v5"
)

func SetupWebSocketRoutes(r chi.Router, wsHandler *WebSocketHandler) {
	r.Route("/ws", func(r chi.Router) {
		r.HandleFunc("/", wsHandler.HandleWebSocketUpgrade().ServeHTTP)
		r.Get("/health", wsHandler.HandleWebSocketHealthCheck)
		r.Get("/metrics", wsHandler.HandleWebSocketMetrics)
		r.Post("/subscribe", wsHandler.HandleSubscribe)
		r.Post("/unsubscribe", wsHandler.HandleUnsubscribe)
	})
}

func SetupWebSocketRoutesV1(r chi.Router, wsHandler *WebSocketHandler) {
	r.Route("/api/v1/ws", func(r chi.Router) {
		r.HandleFunc("/", wsHandler.HandleWebSocketUpgrade().ServeHTTP)
		r.Get("/health", wsHandler.HandleWebSocketHealthCheck)
		r.Get("/metrics", wsHandler.HandleWebSocketMetrics)
		r.Post("/subscribe", wsHandler.HandleSubscribe)
		r.Post("/unsubscribe", wsHandler.HandleUnsubscribe)
	})
}

func SetupRealtimeRoutes(r chi.Router, wsHandler *WebSocketHandler) {
	r.Route("/api/v1/realtime", func(r chi.Router) {
		r.HandleFunc("/connect", wsHandler.HandleWebSocketUpgrade().ServeHTTP)
		r.Get("/status", wsHandler.HandleWebSocketHealthCheck)
		r.Get("/metrics", wsHandler.HandleWebSocketMetrics)
	})
}
