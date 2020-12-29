package ws

import "fmt"

// Subscription type similar to https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go

// Subscription struct
type Subscription struct {
	name        string
	clients     map[*Client]bool
	broadcast   chan []byte
	subscribe   chan *Client
	unsubscribe chan *Client
}

// Run func
func (s *Subscription) Run() {
	fmt.Printf("'%s' subscription opened\n", s.name)
	defer fmt.Printf("'%s' subscription closed\n", s.name)

	for {
		select {
		// subscribe a client
		case client := <-s.subscribe:
			s.clients[client] = true

		// unsubscribe a client
		case client := <-s.unsubscribe:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
			}

		// message all clients
		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				// Send message to client
				case client.write <- message:
				default:
					fmt.Printf("Subscription '%s' broadcast failed to client '%s'\n", s.name, client.id)
				}
			}
		}
	}
}

// Broadcast func
func (s *Subscription) Broadcast(event string, data []byte) {
	s.broadcast <- NewServerEvent(s.name, event, data).Serialize()
}

// NewSubscription func
func NewSubscription(name string) *Subscription {
	return &Subscription{
		name:        name,
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		subscribe:   make(chan *Client),
		unsubscribe: make(chan *Client)}
}
