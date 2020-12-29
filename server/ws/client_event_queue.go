package ws

import (
	"container/list"
	"sync"
)

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