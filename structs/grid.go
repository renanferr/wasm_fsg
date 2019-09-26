package structs

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)

// Grid manager
type Grid struct {
	Dots            []*Dot
	Dashed          bool
	Lines           bool
	Gravity         float64
	Size            int
	ShouldSpawnDots bool
	Width           float64
	Height          float64
	Ctx             js.Value
	MousePos        Vec2D
}

// Update updates the dot positions and draws
func (grid *Grid) Update(gridTime float64) {
	if grid.Dots == nil {
		return
	}
	grid.Ctx.Call("clearRect", 0, 0, grid.Width, grid.Height)
	grid.Ctx.Set("fillStyle", "rgb(0,0,0)")
	grid.Ctx.Call("fillRect", 0, 0, grid.Width, grid.Height)

	grid.SpawnDots(1, Vec2D{
		grid.MousePos.X,
		grid.MousePos.Y,
	})

	for _, dot := range grid.Dots {
		fmt.Println(len(grid.Dots))
		// if dot.pos.X < 0 {
		// 	dot.pos.X = 0
		// 	dot.dir.X = 0
		// }
		// if dot.pos.X > grid.Width {
		// 	dot.pos.X = grid.Width
		// 	dot.dir.X = 0
		// }

		// if dot.pos.Y < 0 {
		// 	dot.pos.Y = 0
		// 	dot.dir.Y = 0
		// }

		// if dot.pos.Y > grid.Height {
		// 	dot.pos.Y = grid.Height
		// 	dot.dir.Y = 0
		// }

		grid.Ctx.Call("beginPath")
		grid.Ctx.Set("fillStyle", fmt.Sprintf("#%06x", dot.color))
		grid.Ctx.Set("strokeStyle", fmt.Sprintf("#%06x", dot.color))

		grid.Ctx.Set("lineWigridh", grid.Size)
		grid.Ctx.Call("arc", dot.pos.X, dot.pos.Y, grid.Size, 0, 2*math.Pi)
		grid.Ctx.Call("fill")

		// if grid.lines {
		// 	for _, dot2 := range grid.Dots[i+1:] {
		// 		mx := dot2.pos.X - dot.pos.X
		// 		my := dot2.pos.Y - dot.pos.Y
		// 		d := mx*mx + my*my
		// 		if d < lineDistSq {
		// 			grid.Ctx.Set("globalAlpha", 1-d/lineDistSq)
		// 			grid.Ctx.Call("beginPath")
		// 			grid.Ctx.Call("moveTo", dot.pos.X, dot.pos.Y)
		// 			ctx.Call("lineTo", dot2.pos.X, dot2.pos.Y)
		// 			ctx.Call("stroke")
		// 		}
		// 	}
		// }

		dot.SetGravity(grid.Gravity)

		dot.pos.X += dot.dir.X * dot.speed * gridTime
		dot.pos.Y += dot.dir.Y * dot.speed * gridTime

		if dot.pos.Y >= grid.Width {
			grid.RemoveDot(dot)
		}
	}
}

// SetNDots reinitializes dots with n size
func (grid *Grid) SetNDots(n int) {
	grid.Dots = make([]*Dot, n)
	for i := 0; i < n; i++ {
		grid.Dots[i] = &Dot{
			pos: Vec2D{
				rand.Float64() * grid.Width,
				rand.Float64() * grid.Height,
			},
			dir: Vec2D{
				0,
				grid.Gravity,
			},
			color: uint32(rand.Intn(0xFFFFFF)),
			size:  10,
		}
	}
}

// GetNDots returns the current number of dots in dot thing
func (grid *Grid) GetNDots() int {
	return len(grid.Dots)
}

// SpawnDots spawns dots given n number of dots in mouse position
func (grid *Grid) SpawnDots(n int, pos Vec2D) {

	if grid.ShouldSpawnDots {
		for i := 0; i < n; i++ {
			grid.Dots = append(grid.Dots, &Dot{
				pos: pos,
				dir: Vec2D{
					0,
					grid.Gravity,
				},
				color: uint32(rand.Intn(0xFFFFFF)),
				size:  10,
				speed: 160,
			})
		}
	}
}

func (grid *Grid) SetShouldSpawnDots(v bool) {
	grid.ShouldSpawnDots = v
}

// RemoveDot deletes a given dot from the grid
func (grid *Grid) RemoveDot(dot *Dot) {
	i := 0

	for _, d := range grid.Dots {
		if d != dot {
			grid.Dots[i] = d
			i++
		}
	}
	grid.Dots = grid.Dots[:i]
}
