package gameplay

import (
	"encoding/json"
	"math"
	"server/game/gameplay/geometry"
	"server/ws"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Champion struct
type Champion struct {
	id           uuid.UUID
	hitbox       *geometry.Circle
	target       *geometry.Point // moving towards this position
	maxHealth    int
	health       int
	stop         int
	speed        int
	movementLock *sync.Mutex

	shootCooldown *Cooldown
	dashCooldown  *Cooldown
}

// NewChampion func
func NewChampion(id uuid.UUID) *Champion {
	return &Champion{
		id:           id,
		hitbox:       geometry.NewCircle(championStartX, championStartY, championRadius),
		maxHealth:    championMaxHealth,
		health:       championMaxHealth,
		speed:        championSpeed, // units per second
		movementLock: &sync.Mutex{},

		shootCooldown: NewCooldown(shootCooldown),
		dashCooldown:  NewCooldown(dashCooldown),
	}
}

func (champ *Champion) respawn(point *geometry.Point) {
	x, y := point.GetX(), point.GetY()
	champ.health = champ.maxHealth
	champ.hitbox.GetPosition().Move(x, y)
}

func (champ *Champion) shoot(event *ws.ClientEvent, game *Game) {
	if !champ.shootCooldown.isReady() {
		return
	}

	go champ.shootCooldown.start()

	data := &ChampionShootEvent{}
	if err := json.Unmarshal(event.GetData(), data); err != nil {
		return
	}

	team := game.getClientTeam(event.Client)
	origin := champ.hitbox.GetPosition().Copy()
	target := geometry.NewPoint(data.X, data.Y)

	projectile := NewProjectile(origin, target, game, team)
	team.projectiles[projectile] = struct{}{}
}

func (champ *Champion) dash() {
	champ.movementLock.Lock()
	defer champ.movementLock.Unlock()

	if !champ.dashCooldown.isReady() {
		return
	}

	go champ.dashCooldown.start()
	champ.speed *= dashSpeedMultiplier

	go func() {
		time.Sleep(dashDuration)
		champ.movementLock.Lock()
		defer champ.movementLock.Unlock()
		champ.speed /= dashSpeedMultiplier
	}()
}

func (champ *Champion) move(game *Game) {
	champ.movementLock.Lock()
	defer champ.movementLock.Unlock()

	if champ.target == nil {
		return
	}

	position := champ.hitbox.GetPosition()

	dx := champ.target.GetX() - position.GetX()
	dy := champ.target.GetY() - position.GetY()

	if dx == 0 && dy == 0 {
		return
	}

	distance := math.Sqrt(float64(dx*dx + dy*dy))
	speed := float64(champ.speed) / float64(game.tps)
	speedX := int(math.RoundToEven(float64(dx) / distance * speed)) // speed per tick
	speedY := int(math.RoundToEven(float64(dy) / distance * speed)) // speed per tick

	dirX := direction(speedX)
	dirY := direction(speedY)

	// Collision with walls
	position.Shift(speedX, 0)
	for _, wall := range game.walls {
		if !wall.hitbox.HitsCircle(champ.hitbox) {
			continue
		}

		position.Shift(-speedX, 0)
		for i := 0; i < (speedX * dirX); i++ {
			position.Shift(dirX, 0)
			if wall.hitbox.HitsCircle(champ.hitbox) {
				position.Shift(-dirX, 0)
				break
			}
		}
		break
	}

	position.Shift(0, speedY)
	for _, wall := range game.walls {
		if !wall.hitbox.HitsCircle(champ.hitbox) {
			continue
		}

		position.Shift(0, -speedY)
		for i := 0; i < (speedY * dirY); i++ {
			position.Shift(0, dirY)
			if wall.hitbox.HitsCircle(champ.hitbox) {
				position.Shift(0, -dirY)
				break
			}
		}
		break
	}

	if distance < speed {
		champ.target = nil
	}
}

func (champ *Champion) setMovementDirection(event *ws.ClientEvent) {
	champ.movementLock.Lock()
	defer champ.movementLock.Unlock()

	movement := &ChampionMoveEvent{}
	if err := json.Unmarshal(event.GetData(), movement); err != nil {
		return
	}

	champ.target = geometry.NewPoint(movement.X, movement.Y)
}

func (champ *Champion) setMovementSpeed() {

}
