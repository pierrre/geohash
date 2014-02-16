package geohash

import (
	"bytes"
	"fmt"
)

var (
	base32     = []byte("0123456789bcdefghjkmnpqrstuvwxyz")
	defaultBox = Box{Lat: Range{Min: -90, Max: 90}, Lon: Range{Min: -180, Max: 180}}
)

func Encode(lat, lon float64, precision int) string {
	ghb := make([]byte, precision)
	box := defaultBox
	even := true
	for i := 0; i < precision; i++ {
		ci := 0
		for j := 16; j != 0; j >>= 1 {
			var r *Range
			var u float64
			if even {
				r = &box.Lon
				u = lon
			} else {
				r = &box.Lat
				u = lat
			}
			if mid := r.Mid(); u > mid {
				ci += j
				r.Min = mid
			} else {
				r.Max = mid
			}
			even = !even
		}
		ghb[i] += base32[ci]
	}
	return string(ghb)
}

func Decode(gh string) (box Box, err error) {
	box = defaultBox
	even := true
	for i := 0; i < len(gh); i++ {
		ci := bytes.IndexByte(base32, gh[i])
		if ci == -1 {
			err = fmt.Errorf("invalid character at index %d", i)
			return
		}
		for j := 16; j != 0; j >>= 1 {
			var r *Range
			if even {
				r = &box.Lon
			} else {
				r = &box.Lat
			}
			if mid := r.Mid(); ci&j != 0 {
				r.Min = mid
			} else {
				r.Max = mid
			}
			even = !even
		}
	}
	return
}

type Box struct {
	Lat, Lon Range
}

func (b Box) Center() Point {
	return Point{Lat: b.Lat.Mid(), Lon: b.Lon.Mid()}
}

func (b Box) IsPointInside(p Point) bool {
	return b.Lat.IsInside(p.Lat) &&
		b.Lon.IsInside(p.Lon)
}

type Point struct {
	Lat, Lon float64
}

type Range struct {
	Min, Max float64
}

func (r Range) Mid() float64 {
	return (r.Min + r.Max) / 2
}

func (r Range) IsInside(v float64) bool {
	return v >= r.Min && v <= r.Max
}
