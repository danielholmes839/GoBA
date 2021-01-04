package game

import (
	"server/game/geometry"
	"server/ws"
)

// Team struct
type Team struct {
	name    string
	color   string
	size    int
	events  *ws.Subscription
	members map[*ws.Client]*Champion
}

// NewTeam func
func NewTeam(name string, color string) *Team {
	return &Team{
		name: name,
		color: color,
		size: 0,
		events:  ws.NewSubscription("team-events"),
		members: make(map[*ws.Client]*Champion),
	}
}

// Tick func
func (team *Team) tick(game *Game) {
	visibleChampions := []*ChampionJSON{}
	for _, champ := range team.members {
		visibleChampions = append(visibleChampions, NewChampionJSON(champ))
	}

	for otherTeam := range game.teams { 						// for every other team
		if otherTeam == team {
			continue
		}

		for _, otherChampion := range otherTeam.members { // for every champion of the other teame
			p2 := otherChampion.hitbox.GetPosition()		// Line end point
			for _, teamChampion := range team.members {
				p1 := teamChampion.hitbox.GetPosition()
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					visibleChampions = append(visibleChampions, NewChampionJSON(otherChampion))
					break
				}
			}
		}
	}

	team.events.Broadcast("tick", NewTickUpdate(visibleChampions))
}

// AddPlayer to team
func (team *Team) addClient(client *ws.Client, champion *Champion) {
	client.Subscribe(team.events)
	team.members[client] = champion
	team.size++
}

// RemovePlayer from team
func (team *Team) removeClient(client *ws.Client) {
	client.Unsubscribe(team.events)
	delete(team.members, client)
	team.size--
}

// TeamJSON struct
type TeamJSON struct {
	Color string `json:"color"`
	Size  int    `json:"size"`
}
