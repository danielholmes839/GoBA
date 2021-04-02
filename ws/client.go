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
	name          string
	conn          *websocket.Conn
	closed        bool
	closedWg      *sync.WaitGroup
	writing       *sync.Mutex
	write         chan []byte
	subscriptions map[*Subscription]bool
}

// NewClient func
func NewClient(conn *websocket.Conn, name string) *Client {
	client := &Client{
		id:   uuid.New(),
		name: name,
		conn: conn,

		closed:        false,
		closedWg:      &sync.WaitGroup{},
		writing:       &sync.Mutex{},
		write:         make(chan []byte),
		subscriptions: make(map[*Subscription]bool),
	}

	// Add to the wait group
	client.closedWg.Add(1)
	fmt.Printf("client:%s has connected\n", client.id)
	return client
}

// GetID func
func (client *Client) GetID() uuid.UUID {
	return client.id
}

// GetName func
func (client *Client) GetName() string {
	return client.name
}

// ReceiveMessages - read incoming messages, break when the connection is closed
func (client *Client) ReceiveMessages(handler ClientEventHandler) {
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
	client.writing.Lock()
	defer client.writing.Unlock()

	// Write the message
	client.conn.WriteMessage(websocket.TextMessage, data)
}

// WriteMessage - write a message to this client (ServerEvent category will always be "personal")
func (client *Client) WriteMessage(eventName string, data []byte) {
	client.writing.Lock()
	defer client.writing.Unlock()

	// Write the message
	client.conn.WriteMessage(websocket.TextMessage, NewServerEvent("personal", eventName, data).Serialize())
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
	client.closedWg.Wait()
}

// Close the client
func (client *Client) Close() {
	client.writing.Lock()				// Don't close while writing
	defer client.writing.Unlock()

	if client.closed {
		return
	}
	client.closed = true

	for subscription := range client.subscriptions {
		client.Unsubscribe(subscription)
	}

	close(client.write)
	client.conn.Close()
	client.closedWg.Done()

	fmt.Printf("client:%s has disconnected\n", client.id)
}
