package utils

func MapIntToFloat(v, intMin, intMax int, floatMin, floatMax float64) float64 {
	intMinF := float64(intMin)
	var R = (floatMax - floatMin) / (float64(intMax) - intMinF)
	var result = (float64(v)-intMinF)*R + floatMin

	return result
}
