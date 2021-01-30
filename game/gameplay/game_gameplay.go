package gameplay

import (
	"fmt"
	"server/game/gameplay/geometry"
)

// Tick func
func (game *Game) tick() {
	// Empty the event queue
	for event := range game.events.Read() {
		champ := game.getClientChampion(event.Client)
		
		switch event.Name {
		case "move":
			champ.setMovementDirection(event)
		case "shoot":
			champ.shoot(event, game)
		case "dash":
			champ.dash()
			
		default:
			fmt.Printf("unknown game event name: '%s'\n", event.Name)
		}
	}

	for _, client := range game.clients {
		client.champion.move(game)
	}
	
	for team := range game.teams {
		for projectile := range team.projectiles {
			projectile.move()
			projectile.collisions(game)
			if projectile.done() {
				delete(team.projectiles, projectile)
			}	
		}
	}

	// Team tick
	for team := range game.teams {
		team.tick(game)
	}
}

func (game *Game) hasLineOfSight(line *geometry.Line) bool {
	for _, wall := range game.walls {
		if line.HitsRectangle(wall.hitbox) {
			return false
		}
	}

	for _, bush := range game.bushes {
		sourceInBush := line.GetStart().HitsRectangle(bush.hitbox)
		targetInBush := line.GetEnd().HitsRectangle(bush.hitbox)

		if targetInBush && sourceInBush {
			// Both in the same bush then the bush has no effect
			continue
		}

		if targetInBush && !sourceInBush {
			// The target is in the bush but the source is not. There is no line of line of sight
			return false
		}

		if !targetInBush && !sourceInBush && line.HitsRectangle(bush.hitbox) {
			// The bush is acting as a wall
			return false
		}
	}

	return true
}