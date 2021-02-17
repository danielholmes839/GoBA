package gameplay

import (
	"errors"
	"time"
)

// Cooldown struct
type Cooldown struct {
	start    time.Time
	duration time.Duration
}

// NewCooldown func
func NewCooldown(duration time.Duration) *Cooldown {
	return &Cooldown{
		start:    time.Time{},
		duration: duration,
	}
}

func (cd *Cooldown) use() error {
	if !cd.isReady() {
		return errors.New("Bruv we on cooldown")
	}

	cd.start = time.Now()
	return nil 
}

func (cd *Cooldown) isReady() bool {
	t := time.Time{}
	if cd.start == t {
		return true
	}

	return time.Now().Sub(cd.start) > cd.duration
}
