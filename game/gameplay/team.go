package gameplay

import (
	"server/game/gameplay/geometry"
	"server/ws"
)

// Team struct
type Team struct {
	name        string
	color       string
	size        int
	respawn     *geometry.Point
	events      *ws.Subscription
	projectiles map[*Projectile]struct{}
	players     map[*ws.Client]*Champion
}

// NewTeam func
func NewTeam(name string, color string, respawn *geometry.Point) *Team {
	return &Team{
		name:        name,
		color:       color,
		size:        0,
		respawn:     respawn,
		events:      ws.NewSubscription("team-events"),
		projectiles: make(map[*Projectile]struct{}),
		players:     make(map[*ws.Client]*Champion),
	}
}

func (team *Team) tick(game *Game) {
	// Tick
	visibleChampions := []*ChampionJSON{}
	visibleProjectiles := []*ProjectileJSON{}

	// Ally players are visible
	for client, champ := range team.players {
		visibleChampions = append(visibleChampions, NewChampionJSON(client, champ))
	}

	// Ally projectiles are visible
	for projectile := range team.projectiles {
		visibleProjectiles = append(visibleProjectiles, NewProjectileJSON(projectile))
	}

	for otherTeam := range game.teams {
		if otherTeam == team {
			continue
		}

		// Adding visible champions on other teams
		for client, otherChampion := range otherTeam.players {
			vision := false
			p2 := otherChampion.hitbox

			// Check players with vision
			for _, teamChampion := range team.players {
				p1 := teamChampion.hitbox
				// Line of sight champion to enemy
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					vision = true
					break
				}
			}

			if vision {
				visibleChampions = append(visibleChampions, NewChampionJSON(client, otherChampion))
				continue
			}

			for teamProjectile := range team.projectiles {
				p1 := teamProjectile.hitbox
				// Line of sight champion to enemy
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					vision = true
				}
			}

			if vision {
				visibleChampions = append(visibleChampions, NewChampionJSON(client, otherChampion))
				continue
			}
		}

		// Adding visible projectiles from other teams
		for otherProjectile := range otherTeam.projectiles {
			p2 := otherProjectile.hitbox
			for _, teamChampion := range team.players {
				p1 := teamChampion.hitbox

				// Line of sight champion to projectile
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					visibleProjectiles = append(visibleProjectiles, NewProjectileJSON(otherProjectile))
					break
				}
			}
		}
	}

	team.events.Broadcast("tick", NewTickUpdate(visibleChampions, visibleProjectiles))
}

func (team *Team) addClient(client *ws.Client, champion *Champion) {
	// Add a client to the game
	client.Subscribe(team.events)
	team.players[client] = champion
	team.size++
}

func (team *Team) removeClient(client *ws.Client) {
	// Remove a client to the game
	client.Unsubscribe(team.events)
	delete(team.players, client)
	team.size--
}
