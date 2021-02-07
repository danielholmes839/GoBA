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
	id        uuid.UUID
	maxHealth int
	health    int
	stop      int
	speed     int

	hitbox *geometry.Circle
	target *geometry.Point // moving towards this position

	movementLock  *sync.Mutex
	shootCooldown *Cooldown
	dashCooldown  *Cooldown

	lastHit  *ws.Client               // The last client to hit this champion
	lastHits map[*ws.Client]time.Time // The times hit by any other clients
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

		lastHit:  nil,
		lastHits: make(map[*ws.Client]time.Time),
	}
}

func (champ *Champion) damage(damage int, enemy *ws.Client) {
	champ.health -= damage
	champ.lastHit = enemy
	champ.lastHits[enemy] = time.Now()
}

func (champ *Champion) respawn(point *geometry.Point) {
	x, y := point.GetX(), point.GetY()
	champ.health = champ.maxHealth
	champ.hitbox.GetPosition().Move(x, y)
}

func (champ *Champion) death() (*ws.Client, []*ws.Client) {
	now := time.Now()
	assists := make([]*ws.Client, 0)

	for client, timeHit := range champ.lastHits {
		if client != champ.lastHit && now.Before(timeHit.Add(time.Second*10)) {
			assists = append(assists, client)
		}
	}

	return champ.lastHit, assists
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

	projectile := NewProjectile(origin, target, game, event.Client)
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

	// The champion isn't moving
	if champ.target == nil {
		return
	}

	// Calculate the difference between current and target position
	position := champ.hitbox.GetPosition()
	dx := champ.target.GetX() - position.GetX()
	dy := champ.target.GetY() - position.GetY()

	// The target is the current position
	if dx == 0 && dy == 0 {
		return
	}

	// Calculate the speed per tick in each direction
	distance := math.Sqrt(float64(dx*dx + dy*dy))
	speed := float64(champ.speed) / float64(game.tps)               // speed per tick
	speedX := int(math.RoundToEven(float64(dx) / distance * speed)) // speed per tick X axis
	speedY := int(math.RoundToEven(float64(dy) / distance * speed)) // speed per tick Y axis

	// Move the champion
	champ.moveX(game, speedX)
	champ.moveY(game, speedY)

	if distance < speed {
		champ.target = nil
	}
}

func (champ *Champion) moveX(game *Game, speedX int) {
	dirX := direction(speedX)
	position := champ.hitbox.GetPosition()

	position.Shift(speedX, 0)
	for _, wall := range game.walls {
		if !wall.HitsCircle(champ.hitbox) {
			continue
		}

		position.Shift(-speedX, 0)
		for i := 0; i < (speedX * dirX); i++ {
			position.Shift(dirX, 0)
			if wall.HitsCircle(champ.hitbox) {
				position.Shift(-dirX, 0)
				break
			}
		}
		break
	}
}

func (champ *Champion) moveY(game *Game, speedY int) {
	dirY := direction(speedY)
	position := champ.hitbox.GetPosition()

	position.Shift(0, speedY)
	for _, wall := range game.walls {
		if !wall.HitsCircle(champ.hitbox) {
			continue
		}

		position.Shift(0, -speedY)
		for i := 0; i < (speedY * dirY); i++ {
			position.Shift(0, dirY)
			if wall.HitsCircle(champ.hitbox) {
				position.Shift(0, -dirY)
				break
			}
		}
		break
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
