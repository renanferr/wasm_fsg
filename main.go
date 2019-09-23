// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"syscall/js"
)

var (
	width      float64
	height     float64
	mousePos   Vec2D
	ctx        js.Value
	lineDistSq float64 = 100 * 100
)

// Grid manager
type Grid struct {
	dots            []*Dot
	dashed          bool
	lines           bool
	gravity         float64
	size            int
	shouldSpawnDots bool
}

// Update updates the dot positions and draws
func (grid *Grid) Update(gridTime float64) {
	if grid.dots == nil {
		return
	}
	ctx.Call("clearRect", 0, 0, width, height)
	ctx.Set("fillStyle", "rgb(0,0,0)")
	ctx.Call("fillRect", 0, 0, width, height)

	grid.SpawnDots(1)

	for _, dot := range grid.dots {
		dir := Vec2D{}

		if dot.pos.X < 0 {
			dot.pos.X = 0
			dot.dir.X = 0
		}
		if dot.pos.X > width {
			dot.pos.X = width
			dot.dir.X = 0
		}

		if dot.pos.Y < 0 {
			dot.pos.Y = 0
			dot.dir.Y = 0
		}

		if dot.pos.Y > height {
			dot.pos.Y = height
			dot.dir.Y = 0
		}

		dir = dot.dir

		ctx.Call("beginPath")
		ctx.Set("fillStyle", fmt.Sprintf("#%06x", dot.color))
		ctx.Set("strokeStyle", fmt.Sprintf("#%06x", dot.color))

		ctx.Set("lineWigridh", grid.size)
		ctx.Call("arc", dot.pos.X, dot.pos.Y, grid.size, 0, 2*math.Pi)
		ctx.Call("fill")

		// if grid.lines {
		// 	for _, dot2 := range grid.dots[i+1:] {
		// 		mx := dot2.pos.X - dot.pos.X
		// 		my := dot2.pos.Y - dot.pos.Y
		// 		d := mx*mx + my*my
		// 		if d < lineDistSq {
		// 			ctx.Set("globalAlpha", 1-d/lineDistSq)
		// 			ctx.Call("beginPath")
		// 			ctx.Call("moveTo", dot.pos.X, dot.pos.Y)
		// 			ctx.Call("lineTo", dot2.pos.X, dot2.pos.Y)
		// 			ctx.Call("stroke")
		// 		}
		// 	}
		// }

		dot.pos.X += dir.X * dot.speed * gridTime
		dot.pos.Y += dir.Y * dot.speed * gridTime
	}
}

// SetNDots reinitializes dots with n size
func (grid *Grid) SetNDots(n int) {
	grid.dots = make([]*Dot, n)
	for i := 0; i < n; i++ {
		grid.dots[i] = &Dot{
			pos: Vec2D{
				rand.Float64() * width,
				rand.Float64() * height,
			},
			dir: Vec2D{
				rand.NormFloat64(),
				grid.gravity * -1,
			},
			color: uint32(rand.Intn(0xFFFFFF)),
			size:  10,
		}
	}
}

// GetNDots returns the current number of dots in dot thing
func (grid *Grid) GetNDots() int {
	return len(grid.dots)
}

// SpawnDots spawns dots given n number of dots in mouse position
func (grid *Grid) SpawnDots(n int) {

	if grid.shouldSpawnDots {
		for i := 0; i < n; i++ {
			grid.dots = append(grid.dots, &Dot{
				pos: Vec2D{
					mousePos.X,
					mousePos.Y,
				},
				dir: Vec2D{
					0,
					grid.gravity,
				},
				color: uint32(rand.Intn(0xFFFFFF)),
				size:  10,
				speed: 160,
			})
		}
	}
}

func (grid *Grid) setShouldSpawnDots(v bool) {
	grid.shouldSpawnDots = v
}

// Dot represents a dot
type Dot struct {
	pos   Vec2D
	dir   Vec2D
	color uint32
	size  float64
	speed float64
}

// Vec2D represents a 2D Vector
type Vec2D struct {
	X float64
	Y float64
}

func main() {

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "mycanvas")
	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	ctx = canvasEl.Call("getContext", "2d")
	done := make(chan struct{}, 0)
	grid := Grid{gravity: 1, size: 1, shouldSpawnDots: false}

	mouseDownEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// e := args[0]
		grid.setShouldSpawnDots(true)
		return nil
	})
	defer mouseDownEvt.Release()

	mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		mousePos.X = e.Get("clientX").Float()
		mousePos.Y = e.Get("clientY").Float()
		return nil
	})
	defer mouseMoveEvt.Release()

	mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		grid.setShouldSpawnDots(false)
		return nil
	})
	defer mouseUpEvt.Release()

	// Event handler for count range
	// countChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	evt := args[0]
	// 	intVal, err := strconv.Atoi(evt.Get("target").Get("value").String())
	// 	if err != nil {
	// 		println("Invalid value", err)
	// 		return nil
	// 	}
	// 	grid.SetNDots(intVal)
	// 	return nil
	// })
	// defer countChangeEvt.Release()

	// Event handler for speed range
	gravityInputEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		fval, err := strconv.ParseFloat(evt.Get("target").Get("value").String(), 64)
		if err != nil {
			println("invalid value", err)
			return nil
		}
		grid.gravity = fval
		return nil
	})
	defer gravityInputEvt.Release()

	// Event handler for size
	sizeChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		intVal, err := strconv.Atoi(evt.Get("target").Get("value").String())
		if err != nil {
			println("invalid value", err)
			return nil
		}
		grid.size = intVal
		return nil
	})
	defer sizeChangeEvt.Release()

	// Event handler for lines toggle
	lineChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		grid.lines = evt.Get("target").Get("checked").Bool()
		return nil
	})
	defer lineChangeEvt.Release()

	// Event handler for dashed toggle
	dashedChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		grid.dashed = evt.Get("target").Get("checked").Bool()
		return nil
	})
	defer dashedChangeEvt.Release()

	doc.Call("addEventListener", "mousedown", mouseDownEvt)
	doc.Call("addEventListener", "mouseup", mouseUpEvt)
	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
	// doc.Call("getElementById", "count").Call("addEventListener", "change", countChangeEvt)
	doc.Call("getElementById", "gravity").Call("addEventListener", "input", gravityInputEvt)
	// doc.Call("getElementById", "size").Call("addEventListener", "input", sizeChangeEvt)
	// doc.Call("getElementById", "dashed").Call("addEventListener", "change", dashedChangeEvt)
	// doc.Call("getElementById", "lines").Call("addEventListener", "change", lineChangeEvt)

	grid.SetNDots(0)
	grid.lines = false
	var renderFrame js.Func
	var tmark float64
	var markCount = 0
	var tdiffSum float64

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := args[0].Float()
		tdiff := now - tmark
		tdiffSum += now - tmark
		markCount++
		if markCount > 10 {
			doc.Call("getElementById", "fps").Set("innerHTML", fmt.Sprintf("FPS: %.01f", 1000/(tdiffSum/float64(markCount))))
			tdiffSum, markCount = 0, 0
		}
		tmark = now

		// Pull window size to handle resize
		curBodyW := doc.Get("body").Get("clientWidth").Float()
		curBodyH := doc.Get("body").Get("clientHeight").Float()
		if curBodyW != width || curBodyH != height {
			width, height = curBodyW, curBodyH
			canvasEl.Set("width", width)
			canvasEl.Set("height", height)
		}
		grid.Update(tdiff / 1000)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done

}
