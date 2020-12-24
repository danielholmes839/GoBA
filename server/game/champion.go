package game

// Point struct

// Champion struct
type Champion struct {
	position *Point
	// speed
	// abilities
	// stunned, blinded, stealth...
}

// SetPosition func
func (champ *Champion) SetPosition(newPosition *Point) {
	champ.position = newPosition
}


// NewChampion func
func NewChampion() *Champion {
	return &Champion{
		position: &Point{0, 0},
	}
}
