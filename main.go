// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/renanferr/wasm_fsg/structs"
)

var (
	mousePos   structs.Vec2D
	ctx        js.Value
	lineDistSq float64 = 100 * 100
)

func main() {

	doc := js.Global().Get("document")
	// width :=
	// height :=
	grid := structs.Grid{
		Gravity:         1,
		Size:            1,
		ShouldSpawnDots: false,
		Width:           doc.Get("body").Get("clientWidth").Float(),
		Height:          doc.Get("body").Get("clientHeight").Float(),
	}
	canvasEl := doc.Call("getElementById", "glCanvas")
	canvasEl.Set("width", grid.Width)
	canvasEl.Set("height", grid.Height)
	ctx = canvasEl.Call("getContext", "2d")
	grid.Ctx = ctx
	done := make(chan struct{}, 0)

	mouseDownEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		grid.SetShouldSpawnDots(true)
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
		grid.SetShouldSpawnDots(false)
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
		grid.Gravity = fval
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
		grid.Size = intVal
		return nil
	})
	defer sizeChangeEvt.Release()

	// Event handler for lines toggle
	lineChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		grid.Lines = evt.Get("target").Get("checked").Bool()
		return nil
	})
	defer lineChangeEvt.Release()

	// Event handler for dashed toggle
	dashedChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		grid.Dashed = evt.Get("target").Get("checked").Bool()
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
	grid.Lines = false
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
		if curBodyW != grid.Width || curBodyH != grid.Height {
			grid.Width = curBodyW
			grid.Height = curBodyH
			canvasEl.Set("width", grid.Width)
			canvasEl.Set("height", grid.Height)
		}

		grid.MousePos = mousePos

		grid.Update(tdiff / 1000)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done
}
