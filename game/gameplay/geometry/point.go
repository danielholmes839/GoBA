package geometry

// Point struct
type Point struct {
	x int
	y int
}

// NewPoint func
func NewPoint(x int, y int) *Point {
	return &Point{x, y}
}

// GetX func
func (point *Point) GetX() int {
	return point.x
}

// GetY func
func (point *Point) GetY() int {
	return point.y
}

// Shift func
func (point *Point) Shift(dx int, dy int) {
	point.x += dx
	point.y += dy
}

// Move func
func (point *Point) Move(x int, y int) {
	point.x = x
	point.y = y
}

// Copy func
func (point *Point) Copy() *Point {
	return NewPoint(point.x, point.y)
}

// Distance2 (distance squared between two points)
func (point *Point) Distance2(point2 *Point) int {
	return distance2(point, point2)
}


/* ############################ HITBOXES ######################################*/

// HitsLine func
func (point *Point) HitsLine(l *Line) bool {
	return l.HitsPoint(point)
}

// HitsCircle func
func (point *Point) HitsCircle(c *Circle) bool {
	return c.HitsPoint(point)
}

// HitsRectangle func
func (point *Point) HitsRectangle(r *Rectangle) bool {
	return (
		point.x >= r.x &&
		point.x <= r.x + r.w &&
		point.y >= r.y &&
		point.y <= r.y + r.h)
}

