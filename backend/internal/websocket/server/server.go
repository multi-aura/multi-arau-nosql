package server

import (
	"log"
	"multiaura/internal/websocket/client"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép tất cả các kết nối
	},
}

func ServeWs(w http.ResponseWriter, r *http.Request, broadcast func([]byte)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	userID := r.URL.Query().Get("userID") // Lấy ID người dùng từ request
	chatClient := client.NewClient(conn, userID)

	go chatClient.ReadPump(broadcast)
	go chatClient.WritePump()
}
