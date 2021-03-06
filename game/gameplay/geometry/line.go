package geometry

import (
	"math"
)

// Line struct
type Line struct {
	p1 *Point
	p2 *Point
}

// NewLine func
func NewLine(x1 int, y1 int, x2 int, y2 int) *Line {
	return &Line{
		p1: &Point{x1, y1},
		p2: &Point{x2, y2},
	}
}

// Shift func
func (line *Line) Shift(dx int, dy int) {
	line.p1.Shift(dx, dy)
	line.p2.Shift(dx, dy)
}

// Length2 func - Length2 func
func (line *Line) Length2() int {
	return distance2(line.p1, line.p2)
}

// GetStart func
func (line *Line) GetStart() *Point {
	return line.p1
}

// GetEnd func
func (line *Line) GetEnd() *Point {
	return line.p2
}

/* ############################ HITBOXES ######################################*/

// HitsPoint func
func (line *Line) HitsPoint(p *Point) bool {
	buffer := 10
	length := distance2(line.p1, line.p2)                 // Length of the line squared
	d1d2 := distance2(line.p1, p) + distance2(line.p2, p) // Distance squared of
	return (d1d2 >= length-buffer) && (d1d2 <= length+buffer)
}

// HitsLine func
func (l *Line) HitsLine(other *Line) bool {
	p1 := l.p1
	p2 := l.p2
	p3 := other.p1
	p4 := other.p2

	uA := div(((p4.x-p3.x)*(p1.y-p3.y))-((p4.y-p3.y)*(p1.x-p3.x)), ((p4.y-p3.y)*(p2.x-p1.x))-((p4.x-p3.x)*(p2.y-p1.y)))
	uB := div(((p4.x-p3.x)*(p1.y-p3.y))-((p4.y-p3.y)*(p1.x-p3.x)), ((p4.y-p3.y)*(p2.x-p1.x))-((p4.x-p3.x)*(p2.y-p1.y)))

	if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
		x := int(math.Round(float64(p1.x) + (uA * float64(p2.x-p1.x))))
		y := int(math.Round(float64(p1.y) + (uA * float64(p2.y-p1.y))))

		if (x > p1.x && x > p2.x) || (x < p1.x && x < p2.x) {
			return false
		}

		if (y > p1.y && y > p2.y) || (y < p1.y && y < p2.y) {
			return false
		}

		if (x > p3.x && x > p4.x) || (x < p3.x && x < p4.x) {
			return false
		}

		if (y > p3.y && y > p4.y) || (y < p3.y && y < p4.y) {
			return false
		}

		return true
	}
	return false
}

// HitsCircle func
func (l *Line) HitsCircle(c *Circle) bool {
	p1 := l.p1
	p2 := l.p2
	r2 := c.r * c.r

	// Check if the end points of the line hit the circle first
	if distance2(c, p1) < r2 || distance2(c, p2) < r2 {
		return true
	}

	// Closest point
	d := dot(l, c)
	point := &Point{x: p1.x + (d * (p2.x - p1.x)), y: p1.y + (d * (p2.y - p1.y))}
	return distance2(c, point) < r2
}

// HitsRectangle func
func (l *Line) HitsRectangle(r *Rectangle) bool {
	return (l.HitsLine(NewLine(r.x, r.y, r.x+r.w, r.y)) ||
		l.HitsLine(NewLine(r.x, r.y, r.x, r.y+r.h)) ||
		l.HitsLine(NewLine(r.x, r.y+r.h, r.x+r.w, r.y+r.h)) ||
		l.HitsLine(NewLine(r.x+r.w, r.y+r.h, r.x+r.w, r.y)))
}
