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

func (score *Score) addKill() {
	score.Kills++
}

func (score *Score) addDeath() {
	score.Deaths++
}

func (score *Score) addAssist() {
	score.Assists++
}
