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
	write         chan []byte
	writing       *sync.Mutex
	closed        *sync.WaitGroup
	subscriptions map[*Subscription]bool
}

// NewClient func
func NewClient(conn *websocket.Conn, name string) *Client {
	client := &Client{
		id:            uuid.New(),
		name:          name,
		conn:          conn,
		write:         make(chan []byte),
		writing:       &sync.Mutex{},
		closed:        &sync.WaitGroup{},
		subscriptions: make(map[*Subscription]bool),
	}

	// Add to the wait group
	client.closed.Add(1)
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

// WriteMessages func - The messages will come from subscriptions for non subscription message use WriteMessage func
func (client *Client) WriteMessages() {
	for message := range client.write {
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}

// ReceiveMessages - read incoming messages, break when the connection is closed
func (client *Client) ReceiveMessages(handler ClientEventHandler) {
	defer client.Close()
	for {
		// Read message
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}

		// Create an event
		event, _ := NewClientEvent(client, message)
		handler.Handle(event)
	}
}

// WriteMessage - write a message to this client (ServerEvent category will always be "personal")
func (client *Client) WriteMessage(eventName string, data []byte) {
	client.writing.Lock()
	client.conn.WriteMessage(websocket.TextMessage, NewServerEvent("personal", eventName, data).Serialize())
	client.writing.Unlock()
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
	client.closed.Wait()
}

// Close the client
func (client *Client) Close() {
	client.writing.Lock()
	defer client.writing.Unlock()

	for subscription := range client.subscriptions {
		client.Unsubscribe(subscription)
	}

	close(client.write)
	client.conn.Close()
	client.closed.Done()

	fmt.Printf("client:%s has disconnected\n", client.id)
}