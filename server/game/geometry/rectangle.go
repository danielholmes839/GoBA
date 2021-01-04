package geometry

// Rectangle struct
type Rectangle struct {
	p *Point
	w int
	h int
}

// NewRectangle func
func NewRectangle(x int, y int, w int, h int) *Rectangle {
	return &Rectangle{p: &Point{x: x, y: y}, w: w, h: h}
}

// GetPosition func
func (rect *Rectangle) GetPosition() *Point {
	return rect.p
}

// GetWidth func
func (rect *Rectangle) GetWidth() int {
	return rect.w
}

// GetHeight func
func (rect *Rectangle) GetHeight() int {
	return rect.h
}

/* ############################ HITBOXES ######################################*/

// HitsPoint func
func (rect *Rectangle) HitsPoint(p *Point) bool {
	return rect.HitsPoint(p)
}

// HitsLine func
func (rect *Rectangle) HitsLine(l *Line) bool {
	return l.HitsRectangle(rect)
}

// HitsRectangle func
func (rect *Rectangle) HitsRectangle(r *Rectangle) bool {
	r1 := rect
	r2 := r
	return (r1.p.x+r1.w >= r2.p.x &&
		r1.p.x <= r2.p.x+r2.w &&
		r1.p.y+r1.h >= r2.p.y &&
		r1.p.y <= r2.p.y+r2.h)
}

// HitsCircle func
func (rect *Rectangle) HitsCircle(c *Circle) bool {
	return c.HitsRectangle(rect)
}
