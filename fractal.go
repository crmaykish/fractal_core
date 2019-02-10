package main

import (
	"fmt"
	"time"

	mb "github.com/crmaykish/fractals/mandelbrot"
	"github.com/fogleman/gg"
)

const imageWidth = 1000
const imageHeight = 1000

// const center = complex(0, 0)

const center = complex(0.05837764683046145, -0.6561039334139365)

var iterations = 200000
var startingZoom = 0.3
var framesToRender = 1

func main() {
	renderTimeStart := time.Now()

	m := mb.Create(imageWidth, imageHeight, center)
	mb.SetZoom(m, startingZoom)
	mb.SetMaxIterations(m, iterations)

	fmt.Printf("Starting render of %d frame(s)...\n", framesToRender)

	for i := 0; i < framesToRender; i++ {
		frameTimeStart := time.Now()
		filename := fmt.Sprintf("assets/fractal%03d.png", i+101)

		mb.Generate(m)

		renderImage(mb.GetBuffer(m), i+1, filename)

		fmt.Printf("Fr %d / %d | Z: %f | It: %d | Render time: ", i+1, framesToRender, mb.GetZoom(m), mb.GetMaxIterations(m))
		fmt.Println(time.Since(frameTimeStart))

		mb.ScaleZoom(m, 1.2)

		// Hacky solution to increasing iterations as we zoom
		z := mb.GetZoom(m)
		if z < 5000 {
			mb.SetMaxIterations(m, mb.DefaultMaxIterations)
		} else if z < 100000 {
			mb.SetMaxIterations(m, 2500)
		} else if z < 1000000 {
			mb.SetMaxIterations(m, 5000)
		} else if z < 20000000 {
			mb.SetMaxIterations(m, 10000)
		} else if z < 400000000 {
			mb.SetMaxIterations(m, 20000)
		}
	}

	fmt.Print("Total render time: ")
	fmt.Println(time.Since(renderTimeStart))
}

func renderImage(buffer [][]uint32, frame int, filename string) {
	dc := gg.NewContext(imageWidth, imageHeight)

	// Save the buffer to an image
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			iterations := buffer[x][y]
			if iterations == 0 {
				dc.SetRGB255(0, 0, 0)
			} else {
				v := int(iterations)

				r := getColor(v)
				g := getColor(v)
				b := getColor(v)

				dc.SetRGB255(r, g, b)
			}
			dc.SetPixel(x, y)
		}
	}

	dc.SavePNG(filename)
}

func getColor(v int) int {

	var color int

	var temp = v % 511

	if temp <= 255 {
		color = temp
	} else {
		color = temp - 255
	}

	return color
}
