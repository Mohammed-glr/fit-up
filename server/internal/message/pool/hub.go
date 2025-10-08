package pool

import (
	"context"
	"log"
	"sync"

	"github.com/tdmdh/fit-up-server/internal/message/types"
	"golang.org/x/net/websocket"
)

type Hub struct {
	connections   map[string]*websocket.Conn
	register      chan *types.Connection
	unregister    chan *types.Connection
	subscriptions map[string]map[string]bool
	mutex         sync.RWMutex
	done          chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		connections: 	make(map[string]*websocket.Conn),
		register: 		make(chan *types.Connection, 256),
		unregister: 	make(chan *types.Connection, 256),
		subscriptions:  make(map[string]map[string]bool),
		mutex: 			sync.RWMutex{},
		done: 			make(chan struct{}),
	}
}


func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.shutdown()
			return
		case <-h.done:
			return
		case conn := <-h.register:
			h.handleRegister(conn)
		case conn := <-h.unregister:
			h.handleUnregister(conn)
		}
	}
}

func (h *Hub) Stop() {
	close(h.done)
}

func (h *Hub) shutdown() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for userID, conn := range h.connections {
		if conn != nil {
			conn.Close()
		}
		delete(h.connections, userID)
	}
	h.subscriptions = make(map[string]map[string]bool) //
}

func (h *Hub) handleUnregister(conn *types.Connection) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.connections[conn.UserID]; exists {
		delete(h.connections, conn.UserID)
		delete(h.subscriptions, conn.UserID)
		log.Printf("User %s disconnected", conn.UserID)
	}
}

func (h *Hub) handleRegister(conn *types.Connection) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.connections[conn.UserID]; exists {
		log.Printf("User %s reconnected, replacing old connection", conn.UserID)
	}

	h.connections[conn.UserID] = conn.Conn
	log.Printf("User %s connected", conn.UserID)
}


func (h *Hub) Connect(userID string, wsConn *websocket.Conn) {
	if userID == "" || wsConn == nil {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.connections[userID] = wsConn


	if h.subscriptions[userID] == nil {
		h.subscriptions[userID] = make(map[string]bool)
	}

	log.Printf("User %s connected", userID)
}

func (h *Hub) Disconnect(userID string) {
	if userID == "" {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if conn, exists := h.connections[userID]; exists {
		if conn != nil {
			conn.Close()
		}
		delete(h.connections, userID)
	}

	delete(h.connections, userID)

	log.Printf("User %s disconnected from WebSocket hub", userID)
}

func (h *Hub) Subscribe(userID, channel string) {
	if userID == "" || channel == "" {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.subscriptions[userID] == nil {
		h.subscriptions[userID] = make(map[string]bool)
	}

	h.subscriptions[userID][channel] = true
	log.Printf("User %s subscribed to channel %s", userID, channel)
}

func (h *Hub) Unsubscribe(userID, channel string) {
	if userID == "" || channel == "" {
		return
	}
	
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if subs, exists := h.subscriptions[userID]; exists {
		if _, subExists := subs[channel]; subExists {
			delete(subs, channel)
			log.Printf("User %s unsubscribed from channel %s", userID, channel)
		}
		if len(subs) == 0 {
			delete(h.subscriptions, userID)
		}	
	}
}

func (h *Hub) SendMessage(userID, message string) error {
	if userID == "" || message == "" {
		return nil
	}

	h.mutex.RLock()
	conn, exists := h.connections[userID]
	h.mutex.RUnlock()

	if !exists || conn == nil {
		log.Printf("No active connection for user %s", userID)
		return nil
	}

	if _, err := conn.Write([]byte(message)); err != nil {
		log.Printf("Error sending message to user %s: %v", userID, err)
		return err
	}
	
	return nil
}


func (h *Hub) GetConnection(userID string) (*websocket.Conn, bool) {
	if userID == "" {
		return nil, false
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	conn, exists := h.connections[userID]
	return conn, exists
}

func (h *Hub) IsConnected(userID string) bool {
	if userID == "" {
		return false
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, exists := h.connections[userID]
	return exists
}

func (h *Hub) GetActiveConnections() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.connections)
}

func (h *Hub) GetConnectedUsers() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]string, 0, len(h.connections))
	for userID := range h.connections {
		users = append(users, userID)
	}
	return users
}
