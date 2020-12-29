package ws

// EventHandler interface
type EventHandler interface {
	Handle(*ClientEvent)
}
