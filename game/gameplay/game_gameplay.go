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

	/* Projectile movement, collision, deletion */
	for team := range game.teams {
		for projectile := range team.projectiles {
			projectile.move()
			projectile.collisions(game)
			if projectile.done() {
				delete(team.projectiles, projectile)
			}
		}
	}

	deaths := 0
	for client, info := range game.clients {
		champion := info.champion
		if champion.health > 0 {
			continue
		}

		game.getClientScore(client).addDeath()
		champion.respawn(info.team.respawn)

		killer, assists := champion.death()
		game.getClientScore(killer).addKill()
		for _, assist := range assists {
			game.getClientScore(assist).addAssist()
		}
		deaths++

	}

	if deaths > 0 {
		game.global.Broadcast("update-teams", NewTeamsUpdate(game))
	}

	/* Team "ticks"
	- Calculate vision of projectiles, and other players
	*/
	for team := range game.teams {
		team.tick(game)
	}
}

func (game *Game) hasLineOfSight(line *geometry.Line) bool {
	for _, wall := range game.walls {
		if line.HitsRectangle(wall) {
			return false
		}
	}

	for _, bush := range game.bushes {
		sourceInBush := line.GetStart().HitsRectangle(bush)
		targetInBush := line.GetEnd().HitsRectangle(bush)

		if targetInBush && sourceInBush {
			// Both in the same bush then the bush has no effect
			continue
		}

		if targetInBush && !sourceInBush {
			// The target is in the bush but the source is not. There is no line of line of sight
			return false
		}

		if !targetInBush && !sourceInBush && line.HitsRectangle(bush) {
			// The bush is acting as a wall
			return false
		}
	}

	return true
}
