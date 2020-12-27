package game

import (
	"encoding/json"
	"fmt"
	"server/game/geometry"
	"server/ws"

	"github.com/google/uuid"
)

// MoveEvent struct
type MoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Champion struct
type Champion struct {
	id            uuid.UUID
	hitbox        *geometry.Circle
	target        *geometry.Point
	speed         int
	health        int
	// speed
	// abilities
	// stunned, blinded, stealth...
}

// NewChampion func
func NewChampion(id uuid.UUID) *Champion {
	return &Champion{
		id:     id,
		hitbox: geometry.NewCircle(50, 50, 20),
		health: 100,
	}
}

// SetPosition func
func (champ *Champion) SetPosition(x int, y int) {
	champ.hitbox.GetPosition().Move(x, y)
}

// MoveEvent func
func (champ *Champion) MoveEvent(event *ws.ClientEvent) {
	movement := &MoveEvent{}
	if err := json.Unmarshal(event.GetData(), movement); err == nil {
		champ.SetPosition(movement.X, movement.Y)
	} else {
		fmt.Println(err)
	}
}

type ChampionJSON struct {
	ID      uuid.UUID `json:"id"`
	Health  int       `json:"health"`
	Visible bool      `json:"visible"`
	X       int       `json:"x"`
	Y       int       `json:"y"`
}

// NewChampionJSON func
func NewChampionJSON(champ *Champion, visible bool) *ChampionJSON {
	r := &ChampionJSON{
		ID:      champ.id,
		Health:  champ.health,
		Visible: visible,
		X:       champ.hitbox.GetPosition().GetX(),
		Y:       champ.hitbox.GetPosition().GetY(),
	}

	if !visible {
		r.X = 0
		r.Y = 0
		r.Health = 0
	}
	
	return r
}