package main

import (
	"fmt"
	"math/cmplx"
	"strconv"
	"sync"
	"time"

	"github.com/fogleman/gg"
)

var centerX, centerY = 0.252, 0.0

var offset = 4.0

var minX = centerX - offset
var maxX = centerX + offset

var minY = centerY - offset
var maxY = centerY + offset

const imageWidth = 1000
const imageHeight = 1000

var maxIterations = 60

var wg sync.WaitGroup

var blob [imageWidth][imageHeight]uint8

func main() {
	for i := 0; i < 100; i++ {
		render("pics/fractal" + strconv.Itoa(i) + ".png")

		// zoom in
		offset *= 0.90
		minX = centerX - offset
		maxX = centerX + offset
		minY = centerY - offset
		maxY = centerY + offset
	}

}

func render(filename string) {
	dc := gg.NewContext(imageWidth, imageHeight)

	fmt.Println("Rendering started...")
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
				blob[x][y] = uint8(iterate(p))
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()

	fmt.Print("Rendering finished in: ")
	fmt.Println(time.Since(start))

	fmt.Println("Saving image...")

	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			its := blob[x][y]
			if its == 0 {
				dc.SetRGB255(0, 0, 0)
			} else {
				dc.SetRGB255(int(255-its), int(255-its), int(255-its))
			}
			dc.SetPixel(x, imageHeight-y-1)
		}
	}

	dc.SavePNG(filename)

	fmt.Println("Image saved!")
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
