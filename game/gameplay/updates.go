package gameplay

import (
	"encoding/json"
	"server/ws"
)

// SetupUpdate struct - initial update sent to clients
type SetupUpdate struct {
	ID     string      `json:"id"`
	Walls  []*WallJSON `json:"walls"`
	Bushes []*BushJSON `json:"bushes"`
}

// NewSetupUpdate struct
func NewSetupUpdate(game *Game, client *ws.Client) []byte {
	walls := make([]*WallJSON, len(game.walls))
	for i, wall := range game.walls {
		walls[i] = NewWallJSON(wall)
	}

	bushes := make([]*BushJSON, len(game.bushes))
	for i, bush := range game.bushes {
		bushes[i] = NewBushJSON(bush)
	}

	data, _ := json.Marshal(&SetupUpdate{
		ID:    client.GetID().String(),
		Walls: walls,
		Bushes: bushes,
	})

	return data
}

// TeamsUpdate struct
type TeamsUpdate struct {
	Teams   map[string]*TeamJSON `json:"teams"`   // team-name: { color, size }
	Clients map[string]string    `json:"clients"` // id: team-name
}

/* 
NewTeamsUpdate func 
Sent when the teams change
*/
func NewTeamsUpdate(game *Game) []byte {
	r := &TeamsUpdate{
		Teams:   make(map[string]*TeamJSON),
		Clients: make(map[string]string),
	}

	for team := range game.teams {
		r.Teams[team.name] = &TeamJSON{Color: team.color, Size: team.size}

		for client := range team.players {
			r.Clients[client.GetID().String()] = team.name
		}
	}

	data, _ := json.Marshal(r)
	return data
}

// TickUpdate struct
type TickUpdate struct {
	Champions []*ChampionJSON `json:"champions"`
	Projectiles []*ProjectileJSON `json:"projectiles"`
}

// NewTickUpdate func
func NewTickUpdate(champions []*ChampionJSON, projectiles []*ProjectileJSON) []byte {
	r := &TickUpdate{
		Champions: champions,
		Projectiles: projectiles,
	}

	data, _ := json.Marshal(r)
	return data
}
