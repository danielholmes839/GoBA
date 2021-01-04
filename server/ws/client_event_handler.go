package ws

// ClientEventHandler interface
type ClientEventHandler interface {
	Handle(*ClientEvent)
}
