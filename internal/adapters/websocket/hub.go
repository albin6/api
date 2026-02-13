package websocket

import (
	"fmt"
	"sync"

	"github.com/albin6/api/internal/core/domain"
)

type Hub struct {
	// Registered clients mapped by user ID
	clients map[uint]*Client

	// Inbound notifications to be broadcast
	broadcast chan *NotificationMessage

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe access to clients map
	mu sync.RWMutex
}

type NotificationMessage struct {
	UserID       uint
	Notification *domain.Notification
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		broadcast:  make(chan *NotificationMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// If user already has a connection, close the old one
			if existingClient, exists := h.clients[client.userID]; exists {
				close(existingClient.send)
			}
			h.clients[client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			// Only close and delete if this client is still the active one
			// This prevents double-close if a new client registered in the meantime
			if existingClient, exists := h.clients[client.userID]; exists && existingClient == client {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			client, exists := h.clients[message.UserID]
			fmt.Printf("[Hub] Processing broadcast for user %d. Client exists: %v\n", message.UserID, exists)
			if exists {
				select {
				case client.send <- message.Notification:
					fmt.Printf("[Hub] Notification sent to client channel for user %d\n", message.UserID)
				default:
					fmt.Printf("[Hub] Client channel full for user %d, disconnecting\n", message.UserID)
					// Client's send buffer is full, close connection
					h.mu.RUnlock()
					h.unregister <- client
					h.mu.RLock()
				}
			} else {
				fmt.Printf("[Hub] No active client found for user %d. Current clients: %d\n", message.UserID, len(h.clients))
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToUser sends a notification to a specific user
func (h *Hub) BroadcastToUser(userID uint, notification *domain.Notification) {
	fmt.Printf("[Hub] Broadcasting to user %d\n", userID)
	h.broadcast <- &NotificationMessage{
		UserID:       userID,
		Notification: notification,
	}
}

// GetConnectedUserCount returns the number of currently connected users
func (h *Hub) GetConnectedUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
