package game

import (
	"encoding/json"
	"fmt"
	"math"
	"server/game/geometry"
	"server/ws"

	"github.com/google/uuid"
)

// Champion struct
type Champion struct {
	id     uuid.UUID
	hitbox *geometry.Circle
	target *geometry.Point // moving towards this position
	health int
	speed  int
	stop   int
}

// NewChampion func
func NewChampion(id uuid.UUID) *Champion {
	return &Champion{
		id:     id,
		hitbox: geometry.NewCircle(50, 50, 50),
		health: 100,
		speed:  400, // units per second
	}
}

// SetMovementDirection func
func (champ *Champion) setMovementDirection(event *ws.ClientEvent) {
	movement := &ChampionMoveEvent{}
	if err := json.Unmarshal(event.GetData(), movement); err != nil {
		fmt.Println(err)
		return
	}

	champ.target = geometry.NewPoint(movement.X, movement.Y)
}

func direction(a int) int {
	if a == 0 {
		return 0
	} else if a < 0 {
		return -1
	} else {
		return 1
	}
}

// Move func
func (champ *Champion) move(game *Game) {
	if champ.target == nil {
		return
	}

	position := champ.hitbox.GetPosition()

	dx := champ.target.GetX() - position.GetX()
	dy := champ.target.GetY() - position.GetY()

	distance := math.Sqrt(float64(dx*dx + dy*dy))
	speed := float64(champ.speed) / float64(game.tps)

	if distance < speed {
		champ.hitbox.GetPosition().Move(champ.target.GetX(), champ.target.GetY())
		champ.target = nil
		return
	}

	speedI := int(math.RoundToEven(speed))                          // speed per tick
	speedX := int(math.RoundToEven(float64(dx) / distance * speed)) // speed per tick
	speedY := int(math.RoundToEven(float64(dy) / distance * speed)) // speed per tick

	champ.hitbox.GetPosition().Shift(speedX, speedY)
	dirX := direction(speedX)
	dirY := direction(speedY)

	// Collision with walls
	for _, wall := range game.walls {
		if !wall.hitbox.HitsCircle(champ.hitbox) {
			continue
		}

		// Reverse only horizontal translation
		champ.hitbox.GetPosition().Shift(-speedX, 0)
		if !wall.hitbox.HitsCircle(champ.hitbox) {
			champ.hitbox.GetPosition().Shift(0, -speedY)
			champ.hitbox.GetPosition().Shift(0, speedI*dirY)
			continue
		}

		// Reverse only vertical translation
		champ.hitbox.GetPosition().Shift(speedX, -speedY)
		if !wall.hitbox.HitsCircle(champ.hitbox) {
			champ.hitbox.GetPosition().Shift(-speedX, 0)
			champ.hitbox.GetPosition().Shift(speedI*dirX, 0)
			continue
		}

		champ.hitbox.GetPosition().Shift(-speedX, 0)
		champ.target = nil
	}
}

// ChampionJSON struct
type ChampionJSON struct {
	ID      uuid.UUID `json:"id"`
	Health  int       `json:"health"`
	X       int       `json:"x"`
	Y       int       `json:"y"`
}

// NewChampionJSON func
func NewChampionJSON(champ *Champion) *ChampionJSON {
	return &ChampionJSON{
		ID:      champ.id,
		Health:  champ.health,
		X:       champ.hitbox.GetPosition().GetX(),
		Y:       champ.hitbox.GetPosition().GetY(),
	}
}
