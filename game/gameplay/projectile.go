package gameplay

import (
	"math"
	"server/game/gameplay/geometry"
	"server/ws"
)

// Projectile struct
type Projectile struct {
	hit    bool
	speedX int
	speedY int
	origin *geometry.Point
	hitbox *geometry.Circle
	team   *Team
	client *ws.Client
}

// NewProjectile function
func NewProjectile(origin *geometry.Point, target *geometry.Point, game *Game, client *ws.Client) *Projectile {
	x, y := origin.GetX(), origin.GetY()

	dx := (target.GetX() - origin.GetX())
	dy := (target.GetY() - origin.GetY())

	speedPerSecond := projectileSpeed
	speedPerTick := float64(speedPerSecond) / float64(game.tps)
	distance := math.Sqrt(float64((dx * dx) + (dy * dy)))

	speedX := int(math.RoundToEven((float64(dx) / distance) * speedPerTick)) // speed per tick
	speedY := int(math.RoundToEven((float64(dy) / distance) * speedPerTick)) // speed per tick
	
	return &Projectile{
		speedX: speedX,
		speedY: speedY,
		origin: geometry.NewPoint(x, y),
		hitbox: geometry.NewCircle(x, y, projectileRadius),
		team:   game.getClientTeam(client),
		client: client,
	}
}

func (projectile *Projectile) move() {
	projectile.hitbox.GetPosition().Shift(projectile.speedX, projectile.speedY)
}

// Check for collisions with other players
func (projectile *Projectile) collisions(game *Game) {
	for _, info := range game.clients {
		champ := info.champion
		team := info.team

		if team == projectile.team {
			continue
		}

		// The projectiles hit a champion
		if projectile.hitbox.HitsCircle(champ.hitbox) {
			projectile.hit = true
			champ.damage(projectileDamage, projectile.client)
		}
	}
}

// The projectile should be deleted
func (projectile *Projectile) done() bool {
	return projectile.hit || (projectile.hitbox.GetPosition().Distance2(projectile.origin) > (projectileRange * projectileRange))
}
