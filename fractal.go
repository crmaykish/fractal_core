package main

import (
	"fmt"
	"time"

	mb "github.com/crmaykish/fractals/mandelbrot"
	"github.com/crmaykish/fractals/utils"
	"github.com/fogleman/gg"
)

const imageWidth = 4000
const imageHeight = 3200
const center = complex(0, 0)

// const center = complex(0.05837764683046145, -0.6561039334139365)

var iterations = 1000
var startingZoom = 0.4
var framesToRender = 1

func main() {
	renderTimeStart := time.Now()

	m := mb.Create(imageWidth, imageHeight, center)
	mb.SetZoom(m, startingZoom)
	mb.SetMaxIterations(m, iterations)

	fmt.Printf("Starting render of %d frame(s)...\n", framesToRender)

	for i := 0; i < framesToRender; i++ {
		frameTimeStart := time.Now()
		filename := fmt.Sprintf("assets/fractal%03d.png", i+1)

		mb.Generate(m)

		renderImage(mb.GetBuffer(m), i+1, filename)

		fmt.Printf("Fr %d / %d | Z: %f | It: %d | Render time: ", i+1, framesToRender, mb.GetZoom(m), mb.GetMaxIterations(m))
		fmt.Println(time.Since(frameTimeStart))

		mb.ScaleZoom(m, 2.0)

		// Hacky solution to increasing iterations as we zoom
		z := mb.GetZoom(m)
		if z < 5000 {
			mb.SetMaxIterations(m, mb.DefaultMaxIterations)
		} else if z < 100000 {
			mb.SetMaxIterations(m, 2000)
		} else if z < 1000000 {
			mb.SetMaxIterations(m, 5000)
		} else if z < 20000000 {
			mb.SetMaxIterations(m, 20000)
		} else if z < 400000000 {
			mb.SetMaxIterations(m, 50000)
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
				v := int(iterations % 255)

				r := v
				g := v
				b := int(utils.MapFloats(float64(v), 0, 255, 0x50, 0xFF))

				dc.SetRGB255(r, g, b)
			}
			dc.SetPixel(x, y)
		}
	}

	dc.SavePNG(filename)
}
