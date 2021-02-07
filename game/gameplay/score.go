package gameplay

// Score struct
type Score struct {
	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Assists int `json:"assists"`
}

// NewScore func
func NewScore(kills int, deaths int, assists int) *Score {
	return &Score{Kills: kills, Deaths: deaths, Assists: assists}

}

// add a kill to the players score
func (score *Score) addKill() {
	score.Kills++
}

// add a death to the players score
func (score *Score) addDeath() {
	score.Deaths++
}

// add an assist to the player score
func (score *Score) addAssist() {
	score.Assists++
}
