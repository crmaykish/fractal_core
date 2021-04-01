package fractal_core

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

func MapIntToInt(v, aMin, aMax, bMin, bMax int) int {
	var R = (bMax - bMin) / (aMax - aMin)
	var result = (v-aMin)*R + bMin

	return result
}

func InterpColors(colorA, colorB uint32, hue float64) (uint8, uint8, uint8) {
	var ra = uint8(colorA & (0xFF << 16) >> 16)
	var ga = uint8(colorA & (0xFF << 8) >> 8)
	var ba = uint8(colorA & 0xFF)
	var rb = uint8(colorB & (0xFF << 16) >> 16)
	var gb = uint8(colorB & (0xFF << 8) >> 8)
	var bb = uint8(colorB & 0xFF)

	return uint8(float64(rb-ra)*hue) + ra, uint8(float64(gb-ga)*hue) + ga, uint8(float64(bb-ba)*hue) + ba
}
