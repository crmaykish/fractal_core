package mandelbrot

import (
	"math/cmplx"
	"sync"

	"github.com/crmaykish/fractals/utils"
)

const DefaultZoomLevel = 0.25
const DefaultMaxIterations = 1000

type Mandelbrot struct {
	ImageWidth             int
	ImageHeight            int
	center                 complex128
	zoomLevel              float64
	maxIterations          int
	buffer                 [][]uint32
	minX, minY, maxX, maxY float64
}

func Create(width, height int, center complex128) *Mandelbrot {
	// Create the main struct
	m := Mandelbrot{ImageWidth: width, ImageHeight: height, center: center}

	// Set up default configuration
	SetMaxIterations(&m, DefaultMaxIterations)
	SetZoom(&m, DefaultZoomLevel)

	// Create a buffer to store all pixels
	m.buffer = make([][]uint32, width)
	for i := 0; i < width; i++ {
		m.buffer[i] = make([]uint32, height)
	}

	return &m
}

func Generate(m *Mandelbrot) {
	var wg sync.WaitGroup

	for x := 0; x < m.ImageWidth; x++ {
		for y := 0; y < m.ImageHeight; y++ {
			// TODO make this mapping better with complex
			// Map this pixel to a complex number on the plane
			var a = utils.MapFloats(float64(x), 0, float64(m.ImageWidth), m.minX, m.maxX)
			var b = utils.MapFloats(float64(y), 0, float64(m.ImageHeight), m.minY, m.maxY)

			// p is a complex number of the form a+bi
			var p = complex(a, b)

			// Iterate this point
			wg.Add(1)
			go func(x, y int) {
				// The number of iterations this point endured is returned and stored in the blob array
				m.buffer[x][y] = iteratePoint(p, m.maxIterations)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()
}

func SetZoom(m *Mandelbrot, z float64) {
	m.zoomLevel = z

	offset := 1.0 / m.zoomLevel
	stretch := float64(m.ImageHeight) / float64(m.ImageWidth)

	// Set the range of the X axis
	m.minX = real(m.center) - offset
	m.maxX = real(m.center) + offset

	// Set the range of the Y access
	// Account for vertical stretch due to non-square image size
	m.minY = imag(m.center) - offset*stretch
	m.maxY = imag(m.center) + offset*stretch
}

func ScaleZoom(m *Mandelbrot, scale float64) {
	SetZoom(m, m.zoomLevel*scale)
}

func GetString() {

}

func GetBuffer(m *Mandelbrot) [][]uint32 {
	return m.buffer
}

func GetZoom(m *Mandelbrot) float64 {
	return m.zoomLevel
}

func GetMaxIterations(m *Mandelbrot) int {
	return m.maxIterations
}

func SetMaxIterations(m *Mandelbrot, i int) {
	m.maxIterations = i
}

// fc(z) = z^2 + c
func iteratePoint(val complex128, maxIterations int) uint32 {
	var curr complex128

	for i := 1; i <= maxIterations; i++ {
		curr = cmplx.Pow(curr, 2) + val

		if cmplx.Abs(curr) > 2.0 {
			return uint32(i)
		}
	}
	return 0
}
