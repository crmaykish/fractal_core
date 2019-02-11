package main

import (
	"fmt"
	"time"

	mb "github.com/crmaykish/fractals/mandelbrot"
	"github.com/fogleman/gg"
)

const imageWidth = 800
const imageHeight = 600

// const center = complex(0, 0)
// const center = complex(0.25, 0)
const center = complex(0.05837764683046145, -0.6561039334139365)

var iterations = 20000
var startingZoom = 2000000.0
var zoomScale = 1.2
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

		renderImage(mb.GetBuffer(m), i+1, filename, m)

		fmt.Printf("Fr %d / %d | Z: %f | It: %d | Render time: ", i+1, framesToRender, mb.GetZoom(m), mb.GetMaxIterations(m))
		fmt.Println(time.Since(frameTimeStart))

		mb.ScaleZoom(m, zoomScale)

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

func renderImage(buffer [][]uint32, frame int, filename string, m *mb.Mandelbrot) {
	dc := gg.NewContext(imageWidth, imageHeight)

	// histogram setup crap - move this
	histogram := mb.GetHistogram(m)
	maxIterations := mb.GetMaxIterations(m)

	var histTotal uint32
	for i := 0; i < maxIterations; i++ {
		histTotal += histogram[i]
	}

	// Save the buffer to an image
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			iterations := buffer[x][y]
			if iterations == uint32(maxIterations) {
				dc.SetRGB255(0, 0, 0)
			} else {
				r, g, b := getColor(iterations, histogram, histTotal)

				dc.SetRGB255(r, g, b)
			}
			dc.SetPixel(x, y)
		}
	}

	dc.SavePNG(filename)
}

func getColor(v uint32, histogram []uint32, histTotal uint32) (int, int, int) {

	palette := make([]uint32, 16)

	// Seems pretty good
	// Need to smooth the transitions between colors
	// rotate the indexes of the array as we descend or the inside and
	// outside is always the same color

	palette[0] = 0x001F3F
	palette[1] = 0x0074D9
	palette[2] = 0x7FDBFF
	palette[3] = 0x39CCCC
	palette[4] = 0x3D9970
	palette[5] = 0x2ECC40
	palette[6] = 0x01FF70
	palette[7] = 0xFFDC00
	palette[8] = 0xFF851B
	palette[9] = 0xFF4136
	palette[10] = 0x85144b
	palette[11] = 0xF012BE
	palette[12] = 0xB10DC9

	var hue float64

	for i := 0; i <= int(v); i++ {
		hue += float64(histogram[i]) / float64(histTotal)
	}

	index := int(12 * hue)

	selected := uint32(palette[index])

	// Split selected into R, G, B channels
	r := (selected & 0xFF0000) >> 16
	g := (selected & 0xFF00) >> 8
	b := (selected & 0xFF)

	return int(r), int(g), int(b)
}
