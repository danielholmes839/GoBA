package ws

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	id            uuid.UUID
	conn          *websocket.Conn
	write         chan []byte
	close         chan bool
	subscriptions map[*Subscription]bool
}

// NewClient func
func NewClient(conn *websocket.Conn) *Client {
	client := &Client{
		id:            uuid.New(),
		conn:          conn,
		write:         make(chan []byte),
		close:         make(chan bool),
		subscriptions: make(map[*Subscription]bool),
	}
	fmt.Printf("client:%s has connected\n", client.id)
	return client
}

// Read func
func (client *Client) Read() {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			break
		}
		log.Println(string(message))
	}
}

// Write func
func (client *Client) Write() {
	for message := range client.write {
		client.conn.WriteMessage(websocket.BinaryMessage, message)
	}
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

// Close func
func (client *Client) Close() {
	for subscription := range client.subscriptions {
		client.Unsubscribe(subscription)
	}

	client.conn.Close()
	close(client.write)

	select {
	case client.close <- true:
	default:
	}
	fmt.Printf("client:%s has disconnected\n", client.id)
}

// Wait func
// block until the websocket disconnects
func (client *Client) Wait() {
	<-client.close
}
