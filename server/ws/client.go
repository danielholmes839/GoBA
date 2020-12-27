package ws

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	id            uuid.UUID
	conn          *websocket.Conn
	close         sync.WaitGroup
	write         chan []byte
	subscriptions map[*Subscription]bool
}

// NewClient func
func NewClient(conn *websocket.Conn) *Client {
	client := &Client{
		id:            uuid.New(),
		conn:          conn,
		write:         make(chan []byte),
		subscriptions: make(map[*Subscription]bool),
	}
	client.close.Add(1)
	fmt.Printf("client:%s has connected\n", client.id)
	return client
}

// GetID func
func (client *Client) GetID() uuid.UUID {
	return client.id
}

// ReadMessages - read incoming messages, break when the connection is closed
func (client *Client) ReadMessages(eventHandler EventHandler) {
	for {
		// Read message
		_, message, err := client.conn.ReadMessage()

		// Client disconnects
		if err != nil {
			client.Close()
			break
		}

		// Get the event
		event, _ := NewClientEvent(client, message)
		eventHandler.Handle(event)
	}
}

// WriteMessages - write messages that are sent by subscriptions or
func (client *Client) WriteMessages() {
	for message := range client.write {
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}

// WriteMessage - write a message to this client (ServerEvent category will always be "personal")
func (client *Client) WriteMessage(eventName string, data []byte) {
	client.write <- NewServerEvent("personal", eventName, data).Serialize()
}

// Subscribe func
func (client *Client) Subscribe(subscription *Subscription) {
	subscription.subscribe <- client
	client.subscriptions[subscription] = true
}

// Unsubscribe func
func (client *Client) Unsubscribe(subscription *Subscription) {
	subscription.unsubscribe <- client
	delete(client.subscriptions, subscription)
}

// Wait - blocks until the client closes
func (client *Client) Wait() {
	client.close.Wait()
}

// Close the client
func (client *Client) Close() {
	// Unsubscribe
	for subscription := range client.subscriptions {
		client.Unsubscribe(subscription)
	}

	// Close channels, connections, wait group
	close(client.write)
	client.conn.Close()
	client.close.Done()

	fmt.Printf("client:%s has disconnected\n", client.id)
}
