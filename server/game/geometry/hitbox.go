package geometry

// Hitbox interface
type Hitbox interface {
	HitsPoint(p *Point) bool
	HitsLine(l *Line) bool
	HitsCircle(c *Circle) bool
	HitsRectangle(r *Rectangle) bool
}
