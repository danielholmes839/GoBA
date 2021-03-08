package ws

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler interface {
	Handle(*ClientEvent)
}

// Client struct
type Client struct {
	*sync.Mutex
	id            uuid.UUID
	username      string
	conn          *websocket.Conn
	running       *sync.WaitGroup
	subscriptions map[*Subscription]bool
}

// NewClient func
func NewClient(conn *websocket.Conn, username string) *Client {
	client := &Client{
		Mutex:         &sync.Mutex{},
		id:            uuid.New(),
		username:      username,
		conn:          conn,
		running:       &sync.WaitGroup{},
		subscriptions: make(map[*Subscription]bool),
	}

	// Add to the wait group
	client.running.Add(1)
	fmt.Printf("client:%s has connected\n", client.id)
	return client
}

// GetID func
func (client *Client) GetID() uuid.UUID {
	return client.id
}

// GetName func
func (client *Client) GetUsername() string {
	return client.username
}

// ReceiveMessages - read incoming messages, break when the connection is closed
func (client *Client) SetHandler(handler Handler) {
	for {
		// Read message
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			client.Close()
			return
		}

		// Create an event
		event, _ := NewClientEvent(client, message)
		handler.Handle(event)
	}
}

// Write bytes to client
func (client *Client) Write(data []byte) {
	client.Lock()
	defer client.Unlock()

	// Write the message
	client.conn.WriteMessage(websocket.TextMessage, data)
}

// WriteMessage - write a message to this client (ServerEvent category will always be "personal")
func (client *Client) WriteMessage(eventName string, data []byte) {
	client.Lock()
	defer client.Unlock()

	// Write the message
	client.conn.WriteMessage(websocket.TextMessage, NewServerEvent("direct", eventName, data).Serialize())
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
	client.running.Wait()
}

// Close func
func (client *Client) Close() {
	client.Lock()
	defer client.Unlock()

	if client.running == nil {
		return
	}

	for subscription := range client.subscriptions {
		client.Unsubscribe(subscription)
	}

	client.conn.Close()
	client.running.Done()
	client.running = nil

	fmt.Printf("client:%s has disconnected\n", client.id)
}
