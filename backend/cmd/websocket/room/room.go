package room

import "multiaura/cmd/websocket/client"

// Room đại diện cho phòng chat
type Room struct {
	Clients    map[*client.Client]bool
	Broadcast  chan []byte
	Register   chan *client.Client
	Unregister chan *client.Client
}

// NewRoom tạo một room mới
func NewRoom() *Room {
	return &Room{
		Clients:    make(map[*client.Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
	}
}

// Run chạy vòng lặp để quản lý việc đăng ký, hủy đăng ký, và phát tin nhắn
// Room struct và các phương thức của nó
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
		case message := <-r.Broadcast: // Khi nhận được tin nhắn broadcast
			for client := range r.Clients { // Gửi đến tất cả các client
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}
