package pool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tdmdh/fit-up-server/internal/message/types"
	"golang.org/x/net/websocket"
)

const (
	writeTimeout     = 10 * time.Second
	pingInterval     = 30 * time.Second
	pongTimeout      = 60 * time.Second
	maxMessageBuffer = 256
)

type Hub struct {
	connections   map[string]*Connection
	register      chan *types.Connection
	unregister    chan string
	broadcast     chan *BroadcastMessage
	subscriptions map[string]map[string]bool
	mutex         sync.RWMutex
	done          chan struct{}
}

type Connection struct {
	conn         *websocket.Conn
	userID       string
	send         chan []byte
	hub          *Hub
	lastPong     time.Time
	mu           sync.Mutex
}

type BroadcastMessage struct {
	Channel string
	Message string
	UserIDs []string 
}

func NewHub() *Hub {
	return &Hub{
		connections:   make(map[string]*Connection),
		register:      make(chan *types.Connection, 256),
		unregister:    make(chan string, 256),
		broadcast:     make(chan *BroadcastMessage, 256),
		subscriptions: make(map[string]map[string]bool),
		mutex:         sync.RWMutex{},
		done:          make(chan struct{}),
	}
}

func (h *Hub) Run(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.shutdown()
			return
		case <-h.done:
			h.shutdown()
			return
		case conn := <-h.register:
			h.handleRegister(conn)
		case userID := <-h.unregister:
			h.handleUnregister(userID)
		case msg := <-h.broadcast:
			h.handleBroadcast(msg)
		case <-ticker.C:
			h.checkConnections()
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
			conn.close()
		}
		delete(h.connections, userID)
	}
	h.subscriptions = make(map[string]map[string]bool)
	log.Println("Hub shutdown complete")
}

func (h *Hub) handleRegister(connInfo *types.Connection) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if oldConn, exists := h.connections[connInfo.UserID]; exists {
		log.Printf("User %s reconnecting, closing old connection", connInfo.UserID)
		oldConn.close()
	}

	conn := &Connection{
		conn:     connInfo.Conn,
		userID:   connInfo.UserID,
		send:     make(chan []byte, maxMessageBuffer),
		hub:      h,
		lastPong: time.Now(),
	}

	h.connections[connInfo.UserID] = conn

	if h.subscriptions[connInfo.UserID] == nil {
		h.subscriptions[connInfo.UserID] = make(map[string]bool)
	}

	go conn.writePump()
	go conn.readPump()

	log.Printf("User %s connected", connInfo.UserID)
}

func (h *Hub) handleUnregister(userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if conn, exists := h.connections[userID]; exists {
		conn.close()
		delete(h.connections, userID)
		delete(h.subscriptions, userID)
		log.Printf("User %s disconnected", userID)
	}
}

func (h *Hub) handleBroadcast(msg *BroadcastMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if msg.Channel != "" {
		for userID, channels := range h.subscriptions {
			if channels[msg.Channel] {
				if len(msg.UserIDs) > 0 && !contains(msg.UserIDs, userID) {
					continue
				}
				if conn, exists := h.connections[userID]; exists {
					select {
					case conn.send <- []byte(msg.Message):
					default:
						log.Printf("Send buffer full for user %s, dropping message", userID)
					}
				}
			}
		}
	} else if len(msg.UserIDs) > 0 {
		for _, userID := range msg.UserIDs {
			if conn, exists := h.connections[userID]; exists {
				select {
				case conn.send <- []byte(msg.Message):
				default:
					log.Printf("Send buffer full for user %s, dropping message", userID)
				}
			}
		}
	}
}

func (h *Hub) checkConnections() {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	now := time.Now()
	for userID, conn := range h.connections {
		conn.mu.Lock()
		lastPong := conn.lastPong
		conn.mu.Unlock()

		if now.Sub(lastPong) > pongTimeout {
			log.Printf("User %s timed out, disconnecting", userID)
			go func(uid string) {
				h.unregister <- uid
			}(userID)
		}
	}
}

func (c *Connection) writePump() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if _, err := c.conn.Write(message); err != nil {
				log.Printf("Write error for user %s: %v", c.userID, err)
				c.hub.unregister <- c.userID
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := websocket.Message.Send(c.conn, "ping"); err != nil {
				log.Printf("Ping error for user %s: %v", c.userID, err)
				c.hub.unregister <- c.userID
				return
			}
		}
	}
}

func (c *Connection) readPump() {
	defer func() {
		c.hub.unregister <- c.userID
	}()

	for {
		var msg string
		err := websocket.Message.Receive(c.conn, &msg)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Read error for user %s: %v", c.userID, err)
			}
			return
		}

		c.mu.Lock()
		c.lastPong = time.Now()
		c.mu.Unlock()

		if msg == "pong" {
			continue
		}

		log.Printf("Received message from %s: %s", c.userID, msg)
	}
}

func (c *Connection) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.send:
	default:
		close(c.send)
	}
	c.conn.Close()
}

func (h *Hub) Connect(userID string, wsConn *websocket.Conn) {
	if userID == "" || wsConn == nil {
		return
	}

	h.register <- &types.Connection{
		UserID: userID,
		Conn:   wsConn,
	}
}

func (h *Hub) Disconnect(userID string) {
	if userID == "" {
		return
	}
	h.unregister <- userID
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
		delete(subs, channel)
		log.Printf("User %s unsubscribed from channel %s", userID, channel)
		if len(subs) == 0 {
			delete(h.subscriptions, userID)
		}
	}
}

func (h *Hub) SendMessage(userID, message string) error {
	if userID == "" || message == "" {
		return fmt.Errorf("userID and message cannot be empty")
	}

	h.mutex.RLock()
	conn, exists := h.connections[userID]
	h.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("no active connection for user %s", userID)
	}

	select {
	case conn.send <- []byte(message):
		return nil
	case <-time.After(writeTimeout):
		return fmt.Errorf("timeout sending message to user %s", userID)
	}
}

func (h *Hub) BroadcastToChannel(channel, message string) {
	if channel == "" || message == "" {
		return
	}

	h.broadcast <- &BroadcastMessage{
		Channel: channel,
		Message: message,
	}
}

func (h *Hub) SendToUsers(userIDs []string, message string) {
	if len(userIDs) == 0 || message == "" {
		return
	}

	h.broadcast <- &BroadcastMessage{
		UserIDs: userIDs,
		Message: message,
	}
}

func (h *Hub) BroadcastToAll(message string) {
	if message == "" {
		return
	}

	h.mutex.RLock()
	userIDs := make([]string, 0, len(h.connections))
	for userID := range h.connections {
		userIDs = append(userIDs, userID)
	}
	h.mutex.RUnlock()

	h.SendToUsers(userIDs, message)
}

func (h *Hub) GetConnection(userID string) (*websocket.Conn, bool) {
	if userID == "" {
		return nil, false
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if conn, exists := h.connections[userID]; exists {
		return conn.conn, true
	}
	return nil, false
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

func (h *Hub) GetChannelSubscribers(channel string) []string {
	if channel == "" {
		return nil
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var subscribers []string
	for userID, channels := range h.subscriptions {
		if channels[channel] {
			subscribers = append(subscribers, userID)
		}
	}
	return subscribers
}

func (h *Hub) GetUserSubscriptions(userID string) []string {
	if userID == "" {
		return nil
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	channels, exists := h.subscriptions[userID]
	if !exists {
		return nil
	}

	result := make([]string, 0, len(channels))
	for channel := range channels {
		result = append(result, channel)
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}