package server

import (
	"log"
	"multiaura/internal/websocket/client"
	"multiaura/internal/websocket/group"

	"github.com/gofiber/websocket/v2"
)

// ServeWs xử lý kết nối WebSocket
func ServeWs(c *websocket.Conn, chatGroup *group.Group) {
	defer c.Close()

	// Lấy userID từ query string (hoặc từ token, tùy thuộc vào cách bạn triển khai)
	userID := c.Query("userID")
	if userID == "" {
		log.Println("UserID is missing")
		return
	}

	// Tạo client mới với userID
	chatClient := client.NewClient(c, userID)

	// Đăng ký client vào group
	chatGroup.AddClient(chatClient)
	log.Printf("User %s connected successfully via WebSocket", userID)

	// Đọc và xử lý tin nhắn từ client
	go chatClient.ReadPump(func(message []byte) {
		chatGroup.BroadcastMessage(message)
	}, chatGroup.Unregister)

	// Gửi tin nhắn tới client
	go chatClient.WritePump()
}
