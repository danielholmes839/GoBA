package ws

import (
	"errors"
	"fmt"
)

// Subscription type similar to https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go

// Subscription struct
type Subscription struct {
	name        string
	clients     map[*Client]bool
	stop        chan bool
	broadcast   chan []byte
	subscribe   chan *Client
	unsubscribe chan *Client
}

// NewSubscription func
func NewSubscription(name string) *Subscription {
	subscription := &Subscription{
		name:        name,
		clients:     make(map[*Client]bool),
		stop:        make(chan bool),
		broadcast:   make(chan []byte),
		subscribe:   make(chan *Client),
		unsubscribe: make(chan *Client)}

	go subscription.run()
	return subscription
}

// Run func
func (s *Subscription) run() {
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
				client.Write(message)
			}

		case <-s.stop:
			for client := range s.clients {
				client.Unsubscribe(s)
			}
			close(s.broadcast)
			close(s.stop)
			close(s.subscribe)
			close(s.unsubscribe)
			return
		}
	}
}

// Close the subscription
func (s *Subscription) Stop() error {
	select {
	case s.stop <- true:
		return nil
	default:
		return errors.New("Subscription is already closed")
	}
}

// Broadcast func
func (s *Subscription) Broadcast(event string, data []byte) {
	s.broadcast <- NewServerEvent(s.name, event, data).Serialize()
}
