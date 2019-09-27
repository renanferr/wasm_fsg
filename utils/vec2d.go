package utils

// Vec2D represents a 2D Vector
type Vec2D struct {
	X float64
	Y float64
}

func NewVec2D(X float64, Y float64) *Vec2D {
	v := new(Vec2D)
	v.X = float64(X)
	v.Y = float64(Y)
	return v
}

func (v *Vec2D) SetX(X float64) {
	v.X = X
}

func (v *Vec2D) SetY(Y float64) {
	v.Y = Y
}
