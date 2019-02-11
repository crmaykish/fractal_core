package utils

func MapFloatToFloat(v, aMin, aMax, bMin, bMax float64) float64 {
	var R = (bMax - bMin) / (aMax - aMin)
	var result = (v-aMin)*R + bMin

	return result
}

func MapIntToFloat(v, intMin, intMax int, floatMin, floatMax float64) float64 {
	intMinF := float64(intMin)
	var R = (floatMax - floatMin) / (float64(intMax) - intMinF)
	var result = (float64(v)-intMinF)*R + floatMin

	return result
}
