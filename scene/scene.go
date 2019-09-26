package scene

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"

	. "github.com/renanferr/wasm_fsg/utils"
)

// Scene manager
type Scene struct {
	dots            []*Dot
	gravity         *Vec2D
	ShouldSpawnDots bool
	Width           float64
	Height          float64
	ctx             js.Value
	mousePos        *Vec2D
}

// NewScene creates a new Scene
func NewScene(h float64, w float64, ctx js.Value) *Scene {
	s := new(Scene)
	s.gravity = NewVec2D(0, 1)
	s.dots = []*Dot{}
	s.Width = w
	s.Height = h
	s.ctx = ctx
	s.ShouldSpawnDots = false
	return s
}

func (s *Scene) GetDots() []*Dot {
	return s.dots
}

func (s *Scene) GetGravity() *Vec2D {
	return s.gravity
}

func (s *Scene) SetGravity(g float64) {
	s.gravity = NewVec2D(0, g)
}

func (s *Scene) GetContext() js.Value {
	return s.ctx
}

func (s *Scene) SetContext(ctx js.Value) {
	s.ctx = ctx
}

func (s *Scene) SetMousePos(pos *Vec2D) {
	s.mousePos = pos
}

// Update updates the dot positions and draws
func (s *Scene) Update(gridTime float64) {
	if s.dots == nil {
		return
	}
	s.ctx.Call("clearRect", 0, 0, s.Width, s.Height)
	s.ctx.Set("fillStyle", "rgb(0,0,0)")
	s.ctx.Call("fillRect", 0, 0, s.Width, s.Height)

	s.SpawnDots(1, s.mousePos)

	for _, dot := range s.dots {
		s.ctx.Call("beginPath")
		s.ctx.Set("fillStyle", fmt.Sprintf("#%06x", dot.color))
		s.ctx.Set("strokeStyle", fmt.Sprintf("#%06x", dot.color))

		// s.ctx.Set("lineWigridh", s.size)
		s.ctx.Call("arc", dot.GetX(), dot.GetY(), 1, 0, 2*math.Pi)
		s.ctx.Call("fill")

		// if s.lines {
		// 	for _, dot2 := range s.dots[i+1:] {
		// 		mx := dot2.pos.X - dot.pos.X
		// 		my := dot2.pos.Y - dot.pos.Y
		// 		d := mx*mx + my*my
		// 		if d < lineDistSq {
		// 			s.ctx.Set("globalAlpha", 1-d/lineDistSq)
		// 			s.ctx.Call("beginPath")
		// 			s.ctx.Call("moveTo", dot.pos.X, dot.pos.Y)
		// 			ctx.Call("lineTo", dot2.pos.X, dot2.pos.Y)
		// 			ctx.Call("stroke")
		// 		}
		// 	}
		// }

		dot.SetDirection(s.gravity)
		fmt.Println(dot.GetDirection().X, dot.GetDirection().Y)

		dot.SetPosition(dot.GetDirection().X*dot.speed*gridTime, dot.GetDirection().Y*dot.speed*gridTime)

		// dot.pos.X +=
		// dot.pos.Y +=

		if dot.GetY() >= s.Width {
			s.RemoveDot(dot)
		}
	}
}

// SetNDots reinitializes dots with n size
func (s *Scene) SetNDots(n int) {
	s.dots = make([]*Dot, n)
	for i := 0; i < n; i++ {
		s.dots[i] = NewDot(
			NewVec2D(rand.Float64()*s.Width, rand.Float64()*s.Height),
			s.gravity,
			uint32(rand.Intn(0xFFFFFF)),
		)
	}
}

// GetNDots returns the current number of dots in dot thing
func (s *Scene) GetNDots() int {
	return len(s.dots)
}

// SpawnDots spawns dots given n number of dots in mouse position
func (s *Scene) SpawnDots(n int, pos *Vec2D) {

	if s.ShouldSpawnDots {
		for i := 0; i < n; i++ {
			s.dots = append(s.dots, NewDot(pos, s.GetGravity(), uint32(rand.Intn(0xFFFFFF))))
		}
	}
}

func (s *Scene) SetShouldSpawnDots(v bool) {
	s.ShouldSpawnDots = v
}

// RemoveDot deletes a given dot from the s
func (s *Scene) RemoveDot(dot *Dot) {
	i := 0

	for _, d := range s.dots {
		if d != dot {
			s.dots[i] = d
			i++
		}
	}
	s.dots = s.dots[:i]
}
