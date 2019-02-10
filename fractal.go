package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"sync"
	"time"

	"github.com/fogleman/gg"
)

// const imageWidth = 7680
// const imageHeight = 4320

const imageWidth = 6000
const imageHeight = 3000

const centerX, centerY = 0.05837764683046145, -0.6561039334139365
const framesToRender = 1

var zoomLevel = 0.25
var maxIterations = 6000

var minX float64
var maxX float64
var minY float64
var maxY float64

var buffer [imageWidth][imageHeight]uint8

var wg sync.WaitGroup

func main() {
	zoom(35000.0)

	totalTime := time.Now()

	for i := 0; i < framesToRender; i++ {
		filename := fmt.Sprintf("pics/fractal%03d.png", i)

		render(i+1, filename)
		zoom(1.1)
	}

	fmt.Print("Total render time: ")
	fmt.Println(time.Since(totalTime))

}

func zoom(scale float64) {
	zoomLevel *= scale

	maxIterations += 40

	offset := 1.0 / zoomLevel
	minX = centerX - offset
	maxX = centerX + offset

	// Scale Y to account for non-square images
	minY = centerY - ((offset) * float64(imageHeight) / imageWidth)
	maxY = centerY + ((offset) * float64(imageHeight) / imageWidth)
}

func render(frame int, filename string) {
	dc := gg.NewContext(imageWidth, imageHeight)

	fmt.Printf("Frame: %d/%d | Zoom: %f | Iterations: %d | Time: ", frame, framesToRender, zoomLevel, maxIterations)

	start := time.Now()

	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			// Map this pixel to a complex number on the plane
			var a = mapTo(float64(x), 0, float64(imageWidth), minX, maxX)
			var b = mapTo(float64(y), 0, float64(imageHeight), minY, maxY)

			// p is a complex number of the form a+bi
			var p = complex(a, b)

			// Iterate this point
			wg.Add(1)
			go func(x, y int) {
				// The number of iterations this point endured is returned and stored in the blob array
				buffer[x][y] = uint8(iterate(p))
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()

	// Save the buffer to an image
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			its := buffer[x][y]
			if its == 0 {
				dc.SetRGB255(0, 0, 0)
			} else {
				// Some hacky color smoothing
				j := math.Pow(float64(its), 1.2)
				v := int(j) & 0xFF

				r, g, b := pixelColor(int(v))

				dc.SetRGB255(int(r), int(g), int(b))
			}
			dc.SetPixel(x, y)
		}
	}

	dc.SavePNG(filename)

	fmt.Println(time.Since(start))
}

// fc(z) = z^2 + c
func iterate(val complex128) int {
	var curr complex128

	for i := 1; i <= maxIterations; i++ {
		curr = cmplx.Pow(curr, 2) + val

		if cmplx.Abs(curr) > 2.0 {
			return i
		}
	}
	return 0
}

func mapTo(val, amin, amax, bmin, bmax float64) float64 {
	var R = (bmax - bmin) / (amax - amin)
	var result = (val-amin)*R + bmin

	return result
}

func colorChannelValue(age int, start, end uint8) uint8 {
	var colorChannelValue uint8

	if start < end {
		// If channel increases with age
		if age >= (int(end) - int(start)) {
			colorChannelValue = end
		} else {
			colorChannelValue = start + uint8(age) // losing data by converting int to uint8, but the range checks should protect it
		}
	} else {
		// If channel decreases with age
		if age >= (int(start) - int(end)) {
			colorChannelValue = end
		} else {
			colorChannelValue = start - uint8(age)
		}
	}

	return colorChannelValue
}

func pixelColor(age int) (uint8, uint8, uint8) {
	var r1, g1, b1 uint8 = 0x00, 0x00, 0x50
	var r2, g2, b2 uint8 = 0xFF, 0xFF, 0xFF

	var red = colorChannelValue(age, r1, r2)
	var green = colorChannelValue(age, g1, g2)
	var blue = colorChannelValue(age, b1, b2)

	return red, green, blue
}
