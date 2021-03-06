package geometry

// Rectangle struct
type Rectangle struct {
	Point
	w int
	h int
}

// NewRectangle func
func NewRectangle(x, y, w, h int) *Rectangle {
	return &Rectangle{Point{x, y}, w, h}
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
	return (r1.x+r1.w >= r2.x &&
		r1.x <= r2.x+r2.w &&
		r1.y+r1.h >= r2.y &&
		r1.y <= r2.y+r2.h)
}

// HitsCircle func
func (rect *Rectangle) HitsCircle(c *Circle) bool {
	return c.HitsRectangle(rect)
}
