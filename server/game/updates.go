package game

// Update struct
type Update struct {
	Champions []*ChampionJSON `json:"champions"`
}

// NewTickUpdate func
func (game *Game) NewTickUpdate() *Update {
	updates := []*ChampionJSON{}
	for team := range game.teams {
		for _, champ := range team.members {
			updates = append(updates, NewChampionJSON(champ, true))
		}
	}
	return &Update{Champions: updates}
}

// NewFullUpdate func
func (game *Game) NewFullUpdate() *Update {
	updates := []*ChampionJSON{}
	for team := range game.teams {
		for _, champ := range team.members {
			updates = append(updates, NewChampionJSON(champ, true))
		}
	}
	return &Update{Champions: updates}
}
