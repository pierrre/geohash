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

func normalize(lat, lon float64) (float64, float64) {
	if lat > 90 || lat < -90 {
		lat = centerMod(lat, 360)
		invertLon := true
		if lat < -90 {
			lat = -180 - lat
		} else if lat > 90 {
			lat = 180 - lat
		} else {
			invertLon = false
		}
		if invertLon {
			if lon > 0 {
				lon -= 180
			} else {
				lon += 180
			}
		}
	}
	if lon > 180 || lon <= -180 {
		lon = centerMod(lon, 360)
	}
	return lat, lon
}

func centerMod(x, y float64) float64 {
	r := math.Mod(x, y)
	if r <= 0 {
		r += y
	}
	if r > y/2 {
		r -= y
	}
	return r
}
