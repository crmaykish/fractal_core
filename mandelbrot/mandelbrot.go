package mandelbrot

import (
	"math"
	"math/cmplx"
	"sync"

	"github.com/crmaykish/fractals/utils"
)

const mandelbrotEscapeRadius = 2.0

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
				m.buffer[x][y] = pointInSet(p, m.maxIterations)
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

// Check if the given complex number is in the Mandelbrot set
// If it is, return 0; if not, return the number of iterations
// it took to diverge outside of the escape radius
func pointInSet(val complex128, maxIterations int) uint32 {
	// Split the complex number into real and imaginary parts
	x := real(val)
	y := imag(val)

	// If the given point is in the main cardioid or the period 2 bulb,
	// it's definitely in the set. No need to iterate on it.
	// This is a huge optimization for points near the main cardioid
	if pointInCardioid(x, y) || pointInPeriod2Bulb(x, y) {
		return 0
	}

	// Keep track of the last two iterated points. If the current
	// point has already been seen, it cannot diverge and must be
	// in the set.
	// TODO: Look into generalizing this instead of just keeping
	// track of 2 points. See where the best tradeoff is
	last0 := complex(0, 0)
	last1 := complex(0, 0)

	var curr complex128

	// Iterate the given point through fc(z) = z^2 + c until it
	// diverges outside of the set or the max iteration has been reached
	for i := 1; i <= maxIterations; i++ {
		// Put the current point through the equation
		curr = cmplx.Pow(curr, 2) + val

		if curr == last0 || curr == last1 {
			// If we've seen this point before, it must be in the set
			return 0
		}

		if cmplx.Abs(curr) > mandelbrotEscapeRadius {
			// Point diverged, return the number of iterations it took
			return uint32(i)
		}

		// Update the last points before iterating again
		last1 = last0
		last0 = curr
	}

	// Point did not diverge, assume it's in the set
	return 0
}

func pointInCardioid(a, b float64) bool {
	p := math.Sqrt(math.Pow(a-(0.25), 2) + math.Pow(b, 2))
	comp := p - 2*math.Pow(p, 2) + (0.25)
	return a <= comp
}

func pointInPeriod2Bulb(a, b float64) bool {
	return math.Pow(a+1, 2)+math.Pow(b, 2) <= float64(1)/float64(16)
}
