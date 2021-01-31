package gameplay

import (
	"server/game/gameplay/geometry"
	"math"
)

// Projectile struct
type Projectile struct {
	hit             bool
	damage          int
	speedX          int
	speedY          int
	maxRangeSquared int
	team            *Team
	origin          *geometry.Point
	hitbox          *geometry.Circle
}

// NewProjectile function
func NewProjectile(origin *geometry.Point, target *geometry.Point, game *Game, team *Team) *Projectile {
	x, y := origin.GetX(), origin.GetY()

	dx := (target.GetX() - origin.GetX())
	dy := (target.GetY() - origin.GetY())

	speedPerSecond := projectileSpeed
	speedPerTick := float64(speedPerSecond) / float64(game.tps)
	distance := math.Sqrt(float64((dx * dx) + (dy * dy)))

	speedX := int(math.RoundToEven((float64(dx) / distance) * speedPerTick)) // speed per tick
	speedY := int(math.RoundToEven((float64(dy) / distance) * speedPerTick)) // speed per tick

	return &Projectile{
		speedX:          speedX,
		speedY:          speedY,
		damage:          projectileDamage,
		maxRangeSquared: projectileSpeed * projectileSpeed,
		origin:          geometry.NewPoint(x, y),
		hitbox:          geometry.NewCircle(x, y, projectileRadius),
		team:            team,
	}
}

func (projectile *Projectile) move() {
	projectile.hitbox.GetPosition().Shift(projectile.speedX, projectile.speedY)
}

func (projectile *Projectile) collisions(game *Game) {
	for client, info := range game.clients {
		champ := info.champion
		team := info.team

		if team == projectile.team {
			continue
		}

		if projectile.hitbox.HitsCircle(champ.hitbox) {
			projectile.hit = true
			champ.health -= projectile.damage

			if champ.health <= 0 {
				team := game.getClientTeam(client)
				champ.respawn(team.respawn)
			}
		}
	}
}

func (projectile *Projectile) done() bool {
	return projectile.hit || (projectile.hitbox.GetPosition().Distance2(projectile.origin) > projectile.maxRangeSquared)
}
