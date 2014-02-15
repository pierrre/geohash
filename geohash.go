package geohash

import (
	"fmt"
	"strings"
)

const (
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz"
)

func Encode(lat, lon float64, precision int) (string, error) {
	//TODO
	return "", nil
}

func Decode(gh string) (box Box, err error) {
	box = Box{Lat: Range{Min: -90, Max: 90}, Lon: Range{Min: -180, Max: 180}}

	even := true
	for i, c := range gh {
		v := strings.IndexRune(base32, c)
		if v == -1 {
			err = fmt.Errorf("invalid character '%c' at index %d", c, i)
			return
		}

		for j := 16; j != 0; j >>= 1 {
			var r *Range
			if even {
				r = &box.Lon
			} else {
				r = &box.Lat
			}

			var u *float64
			if v&j != 0 {
				u = &r.Min
			} else {
				u = &r.Max
			}
			*u = r.Mid()

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
