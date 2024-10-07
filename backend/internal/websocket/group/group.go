package group

import (
	"multiaura/internal/websocket/client"
	"sync"
)

type Group struct {
	Clients    map[*client.Client]bool
	Broadcast  chan []byte
	Register   chan *client.Client
	Unregister chan *client.Client
	mutex      sync.Mutex
}

func NewGroup() *Group {
	return &Group{
		Clients:    make(map[*client.Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *client.Client),
		Unregister: make(chan *client.Client),
	}
}

func (g *Group) Run() {
	for {
		select {
		case client := <-g.Register:
			g.mutex.Lock()
			g.Clients[client] = true
			g.mutex.Unlock()
		case client := <-g.Unregister:
			g.mutex.Lock()
			if _, ok := g.Clients[client]; ok {
				delete(g.Clients, client)
				close(client.Send)
			}
			g.mutex.Unlock()
		case message := <-g.Broadcast:
			g.mutex.Lock()
			for client := range g.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(g.Clients, client)
				}
			}
			g.mutex.Unlock()
		}
	}
}

// BroadcastMessage gửi tin nhắn tới tất cả các client trong nhóm
func (g *Group) BroadcastMessage(message []byte) {
	g.Broadcast <- message
}
