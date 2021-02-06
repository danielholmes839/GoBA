package gameplay

import (
	"sync"
	"time"
)

// Cooldown struct
type Cooldown struct {
	ready    bool
	duration time.Duration
	lock     *sync.Mutex
}

// NewCooldown func
func NewCooldown(duration time.Duration) *Cooldown {
	return &Cooldown{
		ready:    true,
		duration: duration,
		lock:     &sync.Mutex{},
	}
}

func (cd *Cooldown) start() {
	cd.setReady(false)
	time.Sleep(cd.duration)
	cd.setReady(true)
}

func (cd *Cooldown) setReady(value bool) {
	cd.lock.Lock()
	defer cd.lock.Unlock()
	cd.ready = value
}

func (cd *Cooldown) isReady() bool {
	cd.lock.Lock()
	defer cd.lock.Unlock()
	return cd.ready
}
