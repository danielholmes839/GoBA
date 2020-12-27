package ws

import (
	"container/list"
	"encoding/json"
	"sync"
	"time"
)

var (
	// UpdateEvent "Name" value
	UpdateEvent = "UPDATE"
)

// EventHandler interface
type EventHandler interface {
	Handle(*ClientEvent)
}

// ServerEvent type
type ServerEvent struct {
	Subscription string          `json:"subscription"`
	Name         string          `json:"event"`
	Timestamp    int64           `json:"timestamp"`
	Data         json.RawMessage `json:"data"`
}

// NewServerEvent func
func NewServerEvent(subscription string, name string, data []byte) *ServerEvent {
	return &ServerEvent{Subscription: subscription, Name: name, Timestamp: time.Now().Unix(), Data: data}
}

// Serialize func
func (serverEvent *ServerEvent) Serialize() []byte {
	data, _ := json.Marshal(serverEvent)
	return data
}

// ClientEvent type
type ClientEvent struct {
	Client    *Client         `json:"-"`
	Category  string          `json:"category"`
	Name      string          `json:"event"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// NewClientEvent func
func NewClientEvent(client *Client, data []byte) (*ClientEvent, error) {
	event := &ClientEvent{}
	err := json.Unmarshal(data, event)
	event.Client = client
	return event, err
}

// GetClient func
func (event *ClientEvent) GetClient() *Client {
	return event.Client
}

// GetCategory func
func (event *ClientEvent) GetCategory() string {
	return event.Category
}

// GetName func
func (event *ClientEvent) GetName() string {
	return event.Name
}

// GetTimestamp func
func (event *ClientEvent) GetTimestamp() int64 {
	return event.Timestamp
}

// GetData func
func (event *ClientEvent) GetData() []byte {
	return event.Data
}

// EventQueue - used to store events for the game
type EventQueue struct {
	queue *list.List
	lock  *sync.Mutex
}

// Events returns a channel that will be used to empty the queue. While reading no more events can be added
func (e *EventQueue) Read() <-chan *ClientEvent {
	e.lock.Lock()
	c := make(chan *ClientEvent)
	go func() {
		for e.queue.Len() > 0 {
			event := e.queue.Front()
			e.queue.Remove(event)
			c <- event.Value.(*ClientEvent)
		}
		e.lock.Unlock()
		close(c)
	}()
	return c
}

// Push an event to the queue.
func (e *EventQueue) Push(event *ClientEvent) {
	e.lock.Lock()
	e.queue.PushBack(event)
	e.lock.Unlock()
}

// NewEventQueue func
func NewEventQueue() *EventQueue {
	return &EventQueue{queue: list.New(), lock: &sync.Mutex{}}
}