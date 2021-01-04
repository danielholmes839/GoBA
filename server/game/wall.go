package game

import "server/game/geometry"

// Wall struct
type Wall struct {
	hitbox *geometry.Rectangle
}

// NewWall func
func NewWall(x int, y int, w int, h int) *Wall {
	return &Wall{hitbox: geometry.NewRectangle(x, y, w, h)}
}

// WallJSON struct
type WallJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// NewWallJSON func
func NewWallJSON(wall *Wall) *WallJSON {
	return &WallJSON{
		X: wall.hitbox.GetPosition().GetX(),
		Y: wall.hitbox.GetPosition().GetY(),
		W: wall.hitbox.GetWidth(),
		H: wall.hitbox.GetHeight()}
}