package geohash

import (
	"math"
)

func round(val float64) float64 {
	if val < 0 {
		return math.Ceil(val - 0.5)
	}
	return math.Floor(val + 0.5)
}

func roundDecimal(val float64, dec int) float64 {
	factor := math.Pow10(dec)
	return round(val*factor) / factor
}
