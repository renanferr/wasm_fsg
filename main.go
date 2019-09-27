package main

import (
	"fmt"
	"strconv"
	"syscall/js"

	. "github.com/renanferr/wasm_fsg/scene"
	. "github.com/renanferr/wasm_fsg/utils"
)

var (
	mousePos   *Vec2D
	ctx        js.Value
	lineDistSq float64 = 100 * 100
)

func main() {

	doc := js.Global().Get("document")

	// scene := structs.Grid{
	// 	Gravity:         1,
	// 	Size:            1,
	// 	ShouldSpawnDots: false,
	// }
	canvasEl := doc.Call("getElementById", "glCanvas")
	canvasEl.Set("height", doc.Get("body").Get("clientHeight").Float())
	canvasEl.Set("width", doc.Get("body").Get("clientWidth").Float())
	scene := NewScene(
		canvasEl.Get("height").Float(),
		canvasEl.Get("width").Float(),
		canvasEl.Call("getContext", "2d"),
	)

	done := make(chan struct{}, 0)

	mouseDownEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		scene.SetShouldSpawnDots(true)
		return nil
	})
	defer mouseDownEvt.Release()

	mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		mousePos = NewVec2D(e.Get("clientX").Float(), e.Get("clientY").Float())
		return nil
	})
	defer mouseMoveEvt.Release()

	mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		scene.SetShouldSpawnDots(false)
		return nil
	})
	defer mouseUpEvt.Release()

	gravityInputEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		fval, err := strconv.ParseFloat(evt.Get("target").Get("value").String(), 64)
		if err != nil {
			println("invalid value", err)
			return nil
		}
		scene.SetGravity(fval)
		return nil
	})
	defer gravityInputEvt.Release()

	// sizeChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	evt := args[0]
	// 	intVal, err := strconv.Atoi(evt.Get("target").Get("value").String())
	// 	if err != nil {
	// 		println("invalid value", err)
	// 		return nil
	// 	}
	// 	scene.Size = intVal
	// 	return nil
	// })
	// defer sizeChangeEvt.Release()

	// Event handler for lines toggle
	// lineChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	evt := args[0]
	// 	scene.Lines = evt.Get("target").Get("checked").Bool()
	// 	return nil
	// })
	// defer lineChangeEvt.Release()

	// // Event handler for dashed toggle
	// dashedChangeEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	evt := args[0]
	// 	scene.Dashed = evt.Get("target").Get("checked").Bool()
	// 	return nil
	// })
	// defer dashedChangeEvt.Release()

	doc.Call("addEventListener", "mousedown", mouseDownEvt)
	doc.Call("addEventListener", "mouseup", mouseUpEvt)
	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
	// doc.Call("getElementById", "count").Call("addEventListener", "change", countChangeEvt)
	doc.Call("getElementById", "gravity").Call("addEventListener", "input", gravityInputEvt)
	// doc.Call("getElementById", "size").Call("addEventListener", "input", sizeChangeEvt)
	// doc.Call("getElementById", "dashed").Call("addEventListener", "change", dashedChangeEvt)
	// doc.Call("getElementById", "lines").Call("addEventListener", "change", lineChangeEvt)

	// scene.Lines = false
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
		if curBodyW != scene.Width || curBodyH != scene.Height {
			scene.Width = curBodyW
			scene.Height = curBodyH
			canvasEl.Set("width", scene.Width)
			canvasEl.Set("height", scene.Height)
		}

		scene.SetMousePos(mousePos)

		scene.Update(tdiff/1000, doc)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done
}
