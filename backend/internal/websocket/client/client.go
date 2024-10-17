package client

import (
	"log"
	"time"

	"github.com/gofiber/websocket/v2" // Sử dụng Fiber WebSocket
)

const (
	writeWait = 10 * time.Second
)

type Client struct {
	Conn   *websocket.Conn // Kết nối WebSocket
	Send   chan []byte     // Kênh để gửi tin nhắn tới client
	UserID string          // Định danh người dùng
}

// NewClient tạo một đối tượng client mới
func NewClient(conn *websocket.Conn, userID string) *Client {
	return &Client{
		Conn:   conn,
		Send:   make(chan []byte, 500), // Đặt một bộ đệm tin nhắn
		UserID: userID,
	}
}

// ReadPump sẽ lắng nghe các tin nhắn từ WebSocket client
func (c *Client) ReadPump(broadcast func([]byte), unregister chan<- *Client) {
	defer func() {
		unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetPongHandler(func(appData string) error {
		log.Printf("Received Pong from client %s", c.UserID)
		return nil
	})

	ticker := time.NewTicker(5 * time.Second) // Thời gian giữa các gói ping, có thể điều chỉnh
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending Ping message to client %s: %v", c.UserID, err)
				return
			}
		default:
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message from client %s: %v", c.UserID, err)
				return
			}

			broadcast(message)
		}
	}
}

// WritePump sẽ lắng nghe và gửi các tin nhắn tới client qua WebSocket
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Printf("Error: Client channel closed for user %s", c.UserID)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if c.Conn != nil {
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := c.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error writing message to client %s: %v", c.UserID, err)
					return
				}
			} else {
				log.Printf("Error: Connection for client %s is nil", c.UserID)
			}
		}
	}
}
