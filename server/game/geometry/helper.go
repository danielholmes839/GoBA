package geometry

func division(v1 int, v2 int) float64 {
	return float64(v1) / float64(v2)
}

func dot(line *Line, p3 *Point) int {
	p1 := line.p1
	p2 := line.p2
	length := distance2(p1, p2)
	return ((p3.x-p1.x)*(p3.x-p2.x) + (p3.y-p1.y)*(p3.y-p2.y)) / length
}

// Distance2 (distance squared)
func distance2(p1 *Point, p2 *Point) int {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return dx*dx + dy*dy
}
