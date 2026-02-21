package websocket

import (
	"fmt"
	"sync"

	"github.com/albin6/api/internal/core/domain"
)

type Hub struct {
	
	clients map[string]*Client

	
	broadcast chan *NotificationMessage

	
	register chan *Client

	
	unregister chan *Client

	
	mu sync.RWMutex
}

type NotificationMessage struct {
	UserID       string
	Notification *domain.Notification
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
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
			
			if existingClient, exists := h.clients[client.userID]; exists {
				close(existingClient.send)
			}
			h.clients[client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			
			
			if existingClient, exists := h.clients[client.userID]; exists && existingClient == client {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			client, exists := h.clients[message.UserID]
			fmt.Printf("[Hub] Processing broadcast for user %s. Client exists: %v\n", message.UserID, exists)
			if exists {
				select {
				case client.send <- message.Notification:
					fmt.Printf("[Hub] Notification sent to client channel for user %s\n", message.UserID)
				default:
					fmt.Printf("[Hub] Client channel full for user %s, disconnecting\n", message.UserID)
					
					h.mu.RUnlock()
					h.unregister <- client
					h.mu.RLock()
				}
			} else {
				fmt.Printf("[Hub] No active client found for user %s. Current clients: %d\n", message.UserID, len(h.clients))
			}
			h.mu.RUnlock()
		}
	}
}


func (h *Hub) BroadcastToUser(userID string, notification *domain.Notification) {
	fmt.Printf("[Hub] Broadcasting to user %s\n", userID)
	h.broadcast <- &NotificationMessage{
		UserID:       userID,
		Notification: notification,
	}
}


func (h *Hub) GetConnectedUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
