package geometry

// Circle struct
type Circle struct {
	Point     // center
	r     int // radius
}

// NewCircle func
func NewCircle(x int, y int, r int) *Circle {
	return &Circle{Point{x, y}, r}
}

// GetRadius func
func (c *Circle) GetRadius() int {
	return c.r
}

/* ############################ HITBOXES ######################################*/

// HitsPoint func
func (c *Circle) HitsPoint(p *Point) bool {
	return distance2(c, p) < c.r*c.r
}

// HitsLine func
func (c *Circle) HitsLine(l *Line) bool {
	return l.HitsCircle(c)
}

// HitsRectangle func
func (c *Circle) HitsRectangle(r *Rectangle) bool {
	x := c.x
	y := c.y

	if c.x < r.x {
		x = r.x
	} else if c.x > r.x+r.w {
		x = r.x + r.w
	}

	if c.y < r.y {
		y = r.y
	} else if c.y > r.y+r.h {
		y = r.y + r.h
	}

	return distance2(c, &Point{x, y}) < c.r*c.r
}

// HitsCircle func
func (c *Circle) HitsCircle(other *Circle) bool {
	r := c.r + other.r
	return distance2(c, other) < r*r
}
