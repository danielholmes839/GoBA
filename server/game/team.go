package game

import (
	"encoding/json"
	"server/game/geometry"
	"server/ws"
)

// Team struct
type Team struct {
	number  int
	events  *ws.Subscription
	members map[*ws.Client]*Champion
	size    int
	color   string
}

// NewTeam func
func NewTeam() *Team {
	return &Team{
		events:  ws.NewSubscription("team"),
		members: make(map[*ws.Client]*Champion),
		size:    0,
		color
	}
}

// Update func
func (team *Team) Update(game *Game) {
	visibleChampions := []*ChampionJSON{}
	for _, champ := range team.members {
		visibleChampions = append(visibleChampions, NewChampionJSON(champ, true))
	}

	for otherTeam := range game.teams { // for every other team
		if otherTeam != team {
			for _, otherChampion := range otherTeam.members { // for every champion of the other team
				visible := false
				p1 := otherChampion.hitbox.GetPosition()
				for _, teamChampion := range team.members {
					p2 := teamChampion.hitbox.GetPosition()
					line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
					if game.LineOfSight(line) {
						visible = true
						break
					}
				}

				visibleChampions = append(visibleChampions, NewChampionJSON(otherChampion, visible))
			}
		}
	}

	data, _ := json.Marshal(&Update{Champions: visibleChampions})
	team.events.Broadcast("update", data)
}

// AddClient to team
func (team *Team) AddClient(client *ws.Client, champ *Champion) {
	client.Subscribe(team.events)
	team.members[client] = champ
	team.size++
}

// RemoveClient client from team
func (team *Team) RemoveClient(client *ws.Client) *Champion {
	temp := team.members[client]
	client.Unsubscribe(team.events)
	delete(team.members, client)
	team.size--
	return temp
}

// GetChampion func
func (team *Team) GetChampion(client *ws.Client) *Champion {
	return team.members[client]
}
