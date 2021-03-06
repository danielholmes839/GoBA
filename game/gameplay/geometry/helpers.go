package geometry

type IPoint interface {
	GetX() int
	GetY() int
}

func div(v1 int, v2 int) float64 {
	return float64(v1) / float64(v2)
}

func dot(line *Line, p3 IPoint) int {
	p1 := line.p1
	p2 := line.p2
	length := distance2(p1, p2)
	return ((p3.GetX()-p1.x)*(p3.GetX()-p2.x) + (p3.GetY()-p1.y)*(p3.GetY()-p2.y)) / length
}

// Distance2 (distance squared)
func distance2(p1 IPoint, p2 IPoint) int {
	dx := p1.GetX() - p2.GetX()
	dy := p1.GetY() - p2.GetY()
	return dx*dx + dy*dy
}
