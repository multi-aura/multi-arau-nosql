package main

import (
	"fmt"
	"log"
	"multiaura/cmd/websocket/server"
	"net/http"
)

func main() {
	// Tạo phòng chat mới và quản lý nó qua server
	chatRoom := server.NewRoom()

	// Chạy room chat trong một goroutine
	go chatRoom.Run()

	// Xử lý WebSocket qua hàm ServeWs từ server
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(chatRoom, w, r)
	})

	// Khởi động server
	fmt.Println("Server websocket started at :3002")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
