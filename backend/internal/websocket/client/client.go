package client

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
}

func NewClient(conn *websocket.Conn, UserID string) *Client {
	return &Client{
		Conn:   conn,
		Send:   make(chan []byte, 500),
		UserID: UserID,
	}
}

func (c *Client) ReadPump(broadcast func([]byte)) {
	defer c.Conn.Close() // Đóng kết nối khi kết thúc

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast(message) // Gửi tin nhắn tới server
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}
