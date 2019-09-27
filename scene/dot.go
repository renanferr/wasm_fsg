package scene

import (
	"fmt"

	. "github.com/renanferr/wasm_fsg/utils"
)

// Dot represents a dot
type Dot struct {
	pos   *Vec2D
	dir   *Vec2D
	color uint32
	size  float64
	speed float64
}

func NewDot(pos *Vec2D, dir *Vec2D, color uint32) *Dot {
	d := new(Dot)
	d.pos = pos
	d.dir = dir
	d.color = color
	d.size = float64(1)
	d.speed = float64(10)
	return d
}

func (d *Dot) SetDirection(v *Vec2D) {
	d.dir = v
}

func (d *Dot) GetDirection() Vec2D {
	return *d.dir
}

func (d *Dot) SetPosition(X float64, Y float64) {
	d.pos.SetX(X)
	d.pos.SetY(Y)
}

func (d *Dot) SetPositionX(X float64) {
	d.pos.SetX(X)
}

func (d *Dot) SetPositionY(Y float64) {
	fmt.Println(Y)
	d.pos.SetY(Y)
}

func (d *Dot) GetPosition() Vec2D {
	return *d.pos
}

func (d *Dot) GetX() float64 {
	return d.pos.X
}

func (d *Dot) GetY() float64 {
	return d.pos.Y
}
