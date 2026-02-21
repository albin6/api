package websocket

import (
	"log"
	"net/http"

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
		
		
		return true
	},
}


func ServeWs(hub *Hub, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := utils.ValidateToken(token, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID := claims.Sub

		
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("websocket upgrade error: %v", err)
			return
		}

		
		client := &Client{
			hub:    hub,
			conn:   conn,
			userID: userID,
			send:   make(chan *domain.Notification, 256),
		}

		hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}
