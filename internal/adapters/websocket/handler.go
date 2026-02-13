package websocket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/albin6/api/config"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, restrict to your frontend domain
		return true
	},
}

// ServeWs handles WebSocket requests from clients
func ServeWs(hub *Hub, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from query parameter
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// Validate token and extract user ID
		claims, err := utils.ValidateToken(token, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, err := strconv.ParseUint(claims.Sub, 10, 32)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("websocket upgrade error: %v", err)
			return
		}

		// Create new client
		client := &Client{
			hub:    hub,
			conn:   conn,
			userID: uint(userID),
			send:   make(chan *domain.Notification, 256),
		}

		// Register client with hub
		client.hub.register <- client

		// Start read and write pumps in separate goroutines
		go client.writePump()
		go client.readPump()
	}
}

// Note: Import domain is needed, add it to imports
// "github.com/albin6/api/internal/core/domain"
