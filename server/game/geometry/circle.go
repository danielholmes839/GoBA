package geometry

// Circle struct
type Circle struct {
	p  *Point // center
	r  int    // radius
}

// NewCircle func
func NewCircle(x int, y int, r int) *Circle {
	return &Circle{p: &Point{x, y}, r: r}
}

// GetPosition func
func (circle *Circle) GetPosition() *Point {
	return circle.p
}

// GetRadius func
func (circle *Circle) GetRadius() int {
	return circle.r
}

/* ############################ HITBOXES ######################################*/

// OverlapsPoint func
func (circle *Circle) HitsPoint(p *Point) bool {
	return distance2(circle.p, p) < circle.r*circle.r
}

// OverlapsLine func
func (circle *Circle) OverlapsLine(l *Line) bool {
	return l.HitsCircle(circle)
}

// HitsRectangle func
func (circle *Circle) HitsRectangle(r *Rectangle) bool {
	c := circle.p
	x := circle.p.x
	y := circle.p.y
	
	if (c.x < r.p.x) {
		x = r.p.x
	} else if (c.x > r.p.x + r.w) {
		x = r.p.x + r.w
	}

	if (c.y < r.p.y) {
		y = r.p.y
	} else if (c.y > r.p.y + r.h) {
		y = r.p.y + r.h
	}

	return distance2(c, &Point{x, y}) < circle.r*circle.r
}

// HitsCircle func
func (circle *Circle) HitsCircle(c *Circle) bool {
	r := circle.r + c.r
	return distance2(circle.p, c.p) < r*r
}