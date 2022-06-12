package main

import (
	"image/color"
	"syscall/js"

	"math/rand"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
)

var done chan struct{}

var grid []bool

var jsCanvas js.Value
var cvs *canvas.Canvas2d
var width int
var height int

var sizeX int = 50
var sizeY int = 50

func randBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func gridAt(i, j int) bool {
	if i < 0 || i >= sizeX || j < 0 || j >= sizeY {
		return false
	}
	return grid[i * sizeX + j]
}

func main() {
	cvs, _ = canvas.NewCanvas2d(false)

	height = 600
	width = 600

	grid = make([]bool, sizeX * sizeY)

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			grid[i * sizeX + j] = randBool()
		}
	}
	
	cvs.Start(30, Render)

	document := js.Global().Get("document")
	jsCanvas = document.Call("createElement", "canvas")

	// set the width and height
	jsCanvas.Set("width", width)
	jsCanvas.Set("height", height)

	// set id = "canvas"
	jsCanvas.Call("setAttribute", "id", "canvas")

	// apply classes
	jsCanvas.Get("classList").Call("add", "w-[80vw]", "h-[80vw]", "md:w-[80vh]", "md:h-[80vh]", "mx-0", "my-0")

	// append to body
	canvasContainer := document.Call("getElementById", "canvas-container")
	cvs.Set(jsCanvas, width, height)
	canvasContainer.Call("appendChild", jsCanvas)
	
	<-done
}

func Render(gc *draw2dimg.GraphicContext) bool {
	// gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Save()

	gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetLineWidth(1)

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			if gridAt(i, j) {
				gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
			} else {
				gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			}
			gc.BeginPath()
			draw2dkit.Rectangle(gc, float64(i)*float64(width)/float64(sizeX), float64(j)*float64(height)/float64(sizeY), float64(i+1)*float64(width)/float64(sizeX), float64(j+1)*float64(height)/float64(sizeY))
			gc.FillStroke()
			gc.Close()
		}
	}
	
	gc.Restore()

	nextGeneration := make([]bool, sizeX * sizeY)

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			alive := 0
			for x := -1; x <= 1; x++ {
				for y := -1; y <= 1; y++ {
					if x == 0 && y == 0 {
						continue
					}
					if gridAt(i+x, j+y) {
						alive++
					}
				}
			}
			/*
			https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life

			The rules of the game are:
			1. Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			2. Any live cell with two or three live neighbours lives on to the next generation.
			3. Any live cell with more than three live neighbours dies, as if by overpopulation.
			4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			*/
			if gridAt(i, j) {
				if alive < 2 || alive > 3 {
					nextGeneration[i * sizeX + j] = false
				} else {
					nextGeneration[i * sizeX + j] = true
				}
			} else {
				if alive == 3 {
					nextGeneration[i * sizeX + j] = true
				} else {
					nextGeneration[i * sizeX + j] = false
				}
			}
		}
	}

	grid = nextGeneration

	return true
}
