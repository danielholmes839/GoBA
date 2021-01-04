package game

import "server/game/geometry"

// Bush struct
type Bush struct {
	hitbox *geometry.Rectangle
}

// NewBush func
func NewBush(x int, y int, w int, h int) *Bush {
	return &Bush{hitbox: geometry.NewRectangle(x, y, w, h)}
}

// BushJSON struct
type BushJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// NewBushJSON func
func NewBushJSON(bush *Bush) *BushJSON {
	return &BushJSON{
		X: bush.hitbox.GetPosition().GetX(),
		Y: bush.hitbox.GetPosition().GetY(),
		W: bush.hitbox.GetWidth(),
		H: bush.hitbox.GetHeight()}
}