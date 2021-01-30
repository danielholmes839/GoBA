package gameplay

import (
	"server/ws"

	"github.com/google/uuid"
)

// TeamJSON struct
type TeamJSON struct {
	Color string `json:"color"`
	Size  int    `json:"size"`
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
		H: bush.hitbox.GetHeight(),
	}
}

// ChampionJSON struct
type ChampionJSON struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Health int       `json:"health"`
	Radius int       `json:"r"`
	X      int       `json:"x"`
	Y      int       `json:"y"`
}

// NewChampionJSON func
func NewChampionJSON(client *ws.Client, champ *Champion) *ChampionJSON {
	return &ChampionJSON{
		ID:     champ.id,
		Name:   client.GetName(),
		Health: champ.health,
		Radius: champ.hitbox.GetRadius(),
		X:      champ.hitbox.GetPosition().GetX(),
		Y:      champ.hitbox.GetPosition().GetY(),
	}
}

// ProjectileJSON struct
type ProjectileJSON struct {
	Team   string `json:"team"`
	Radius int    `json:"r"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// NewProjectileJSON func
func NewProjectileJSON(projectile *Projectile) *ProjectileJSON {
	return &ProjectileJSON{
		Team:   projectile.team.name,
		Radius: projectile.hitbox.GetRadius(),
		X:      projectile.hitbox.GetPosition().GetX(),
		Y:      projectile.hitbox.GetPosition().GetY(),
	}
}

// ChampionMoveEvent struct
type ChampionMoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ChampionShootEvent struct
type ChampionShootEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}
