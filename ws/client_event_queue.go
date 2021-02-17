package ws

import (
	"container/list"
	"sync"
)

// ClientEventQueue - used to store events for the game
type ClientEventQueue struct {
	queue *list.List
	lock  *sync.Mutex
}

// Events returns a channel that will be used to empty the queue. While reading no more events can be added
func (e *ClientEventQueue) Read() []*ClientEvent {
	e.lock.Lock()
	defer e.lock.Unlock()

	events := make([]*ClientEvent, e.queue.Len())
	for i := 0; e.queue.Len() > 0; i++ {
		event := e.queue.Front()
		events[i] = e.queue.Remove(event).(*ClientEvent)
	}

	return events
}

// Push an event to the queue.
func (e *ClientEventQueue) Push(event *ClientEvent) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.queue.PushBack(event)
	
}

// NewClientEventQueue func
func NewClientEventQueue() *ClientEventQueue {
	return &ClientEventQueue{queue: list.New(), lock: &sync.Mutex{}}
}
