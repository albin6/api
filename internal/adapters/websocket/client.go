package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/albin6/api/internal/core/domain"
	"github.com/gorilla/websocket"
)

const (
	
	writeWait = 10 * time.Second

	
	pongWait = 60 * time.Second

	
	pingPeriod = (pongWait * 9) / 10

	
	maxMessageSize = 512
)

type Client struct {
	hub *Hub

	
	conn *websocket.Conn

	
	userID uint

	
	send chan *domain.Notification
}


func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}
		
		
	}
}


func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case notification, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			
			data, err := json.Marshal(notification)
			if err != nil {
				log.Printf("error marshaling notification: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
