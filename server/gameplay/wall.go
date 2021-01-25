package gameplay

import "server/gameplay/geometry"

// Wall struct
type Wall struct {
	hitbox *geometry.Rectangle
}

// NewWall func
func NewWall(x int, y int, w int, h int) *Wall {
	return &Wall{hitbox: geometry.NewRectangle(x, y, w, h)}
}