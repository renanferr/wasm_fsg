package structs

// Dot represents a dot
type Dot struct {
	pos   Vec2D
	dir   Vec2D
	color uint32
	size  float64
	speed float64
}

// SetGravity sets dot gravity
func (d *Dot) SetGravity(v float64) {
	d.dir = Vec2D{
		X: 0,
		Y: v,
	}
}
