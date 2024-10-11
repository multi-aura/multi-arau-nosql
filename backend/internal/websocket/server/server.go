package server

import (
	"log"
	"multiaura/internal/websocket/client"
	"multiaura/internal/websocket/group"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(w http.ResponseWriter, r *http.Request, group *group.Group) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	userID := r.URL.Query().Get("userID")
	chatClient := client.NewClient(conn, userID)

	group.Register <- chatClient // Đăng ký client vào group

	// Đọc và xử lý tin nhắn
	go chatClient.ReadPump(func(message []byte) {
		group.BroadcastMessage(message)
	}, group.Unregister)

	// Gửi tin nhắn tới client
	go chatClient.WritePump()
}
