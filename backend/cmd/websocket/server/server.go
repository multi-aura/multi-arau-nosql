package server

import (
	"log"
	"multiaura/cmd/websocket/client"
	"multiaura/cmd/websocket/room"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép tất cả kết nối
	},
}

// ServeWs xử lý kết nối WebSocket và đăng ký client vào room
func ServeWs(chatRoom *room.Room, w http.ResponseWriter, r *http.Request) {
	// Lấy các query parameters từ URL
	userID := r.URL.Query().Get("user_id")
	conversationID := r.URL.Query().Get("conversation_id")

	// Kiểm tra nếu thiếu user_id hoặc conversation_id
	if userID == "" || conversationID == "" {
		http.Error(w, "Missing user_id or conversation_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	// Tạo client mới
	chatClient := client.NewClient(conn)

	// Đăng ký client vào room
	chatRoom.Register <- chatClient

	// Chạy các goroutine để xử lý việc đọc và ghi tin nhắn
	go chatClient.ReadPump(func(message []byte) {
		chatRoom.Broadcast <- message
	})

	go chatClient.WritePump()
}

// NewRoom tạo và trả về một room mới từ server.go
func NewRoom() *room.Room {
	return room.NewRoom()
}
