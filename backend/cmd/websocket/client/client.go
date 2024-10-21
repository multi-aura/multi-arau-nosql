package client

import (
	"github.com/gorilla/websocket"
)

// Client đại diện cho mỗi kết nối WebSocket
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

// NewClient tạo client mới
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

// ReadPump nhận tin nhắn từ client và gửi nó tới room
func (c *Client) ReadPump(broadcast func([]byte)) {
	defer c.Conn.Close()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast(message)
	}
}

// WritePump gửi tin nhắn từ room tới client qua WebSocket
func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
