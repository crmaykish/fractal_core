package utils

func MapFloats(val, amin, amax, bmin, bmax float64) float64 {
	var R = (bmax - bmin) / (amax - amin)
	var result = (val-amin)*R + bmin

	return result
}

// func MapInts(val, amin, amax, bmin, bmax int) int {
// 	var R = (bmax - bmin) / (amax - amin)
// 	var result = (val-amin)*R + bmin

// 	return result
// }

// func MapComplex() complex128 {
// 	var a = MapFloats(float64(x), 0, float64(imageWidth), minX, maxX)
// 	var b = MapFloats(float64(y), 0, float64(imageHeight), minY, maxY)
// 	var p = complex(a, b)

// 	return p
// }
