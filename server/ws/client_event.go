package ws

import "encoding/json"

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
