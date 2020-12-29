package ws

import (
	"encoding/json"
	"time"
)

// ServerEvent type
type ServerEvent struct {
	Subscription string          `json:"subscription"`
	Name         string          `json:"name"`
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