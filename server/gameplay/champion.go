package gameplay

import (
	"encoding/json"
	"fmt"
	"math"
	"server/gameplay/geometry"
	"server/ws"
	"time"

	"github.com/google/uuid"
)

// Champion struct
type Champion struct {
	id        uuid.UUID
	hitbox    *geometry.Circle
	target    *geometry.Point // moving towards this position
	maxHealth int
	health    int
	speed     int
	stop      int

	shootCooldown *Cooldown
	dashCooldown  *Cooldown
}

// NewChampion func
func NewChampion(id uuid.UUID) *Champion {
	return &Champion{
		id:        id,
		hitbox:    geometry.NewCircle(1500, 1500, 75),
		maxHealth: 100,
		health:    100,
		speed:     700, // units per second

		shootCooldown: NewCooldown(time.Second / 10),
		dashCooldown:  NewCooldown(time.Second * 2),
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

func (champ *Champion) shoot(event *ws.ClientEvent, game *Game) {
	if !champ.shootCooldown.isReady() {
		return
	}

	go champ.shootCooldown.start()

	data := &ChampionShootEvent{}
	if err := json.Unmarshal(event.GetData(), data); err != nil {
		fmt.Println(err)
		return
	}

	team := game.getClientTeam(event.Client)
	origin := champ.hitbox.GetPosition().Copy()
	target := geometry.NewPoint(data.X, data.Y)

	projectile := NewProjectile(origin, target, game, team)
	team.projectiles[projectile] = struct{}{}
}

func (champ *Champion) dash() {
	if !champ.dashCooldown.isReady() {
		return
	}

	go champ.dashCooldown.start()
	champ.speed *= 3

	go func() {
		time.Sleep(time.Second / 5)
		champ.speed /= 3
	}()
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
