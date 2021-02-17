package gameplay

import (
	"fmt"
	"server/game/gameplay/geometry"
)

// Tick func
func (game *Game) tick() {
	// Empty the event queue
	for _, event := range game.events.Read() {
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
	game.processProjectiles()

	/* Check deaths and update the scoreboard */
	if deaths := game.processDeaths(); deaths > 0 {
		game.global.Broadcast("update-teams", NewTeamsUpdate(game))
	}

	/* Team "ticks"
	- Calculate vision of projectiles, and other players
	*/
	for team := range game.teams {
		team.tick(game)
	}
}

func (game *Game) processProjectiles() {
	// Process projectile movement and  collisions
	for team := range game.teams {
		for projectile := range team.projectiles {
			projectile.move()
			projectile.collisions(game)
			if projectile.done() {
				delete(team.projectiles, projectile)
			}
		}
	}
}

func (game *Game) processDeaths() int {
	// Process any deaths
	deaths := 0
	for client, info := range game.clients {
		champion := info.champion
		if champion.health > 0 {
			continue
		}

		game.getClientScore(client).addDeath()
		champion.respawn(info.team.respawn)

		killer, assists := champion.death()

		// Add the kills
		if clientInfo := game.getClientInfo(killer); clientInfo != nil {
			clientInfo.score.addKill()
		}

		// Add assists
		for _, assist := range assists {
			if clientInfo := game.getClientInfo(assist); clientInfo != nil {
				clientInfo.score.addAssist()
			}
		}
		deaths++
	}

	return deaths
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
