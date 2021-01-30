package gameplay

import "server/game/gameplay/geometry"

// Bush struct
type Bush struct {
	hitbox *geometry.Rectangle
}

// NewBush func
func NewBush(x int, y int, w int, h int) *Bush {
	return &Bush{hitbox: geometry.NewRectangle(x, y, w, h)}
}
