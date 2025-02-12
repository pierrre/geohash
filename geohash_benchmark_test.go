package geohash

import (
	"testing"

	codefor_geohash "github.com/Codefor/geohash"
	tomi_hiltunen_geohash "github.com/TomiHiltunen/geohash-golang"
	broadygeohash "github.com/broady/gogeohash" //nolint:misspell
	fanixk_geohash "github.com/fanixk/geohash"
	mmcloughlin_geohash "github.com/mmcloughlin/geohash"
	"github.com/pierrre/assert"
	the42_cartconvert_geohash "github.com/the42/cartconvert/cartconvert"
)

func BenchmarkEncode(b *testing.B) {
	for b.Loop() {
		Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkDecode(b *testing.B) {
	for b.Loop() {
		_, err := Decode(testGeohash)
		if err != nil {
			assert.NoError(b, err)
		}
	}
}

func BenchmarkNeighbors(b *testing.B) {
	for b.Loop() {
		_, err := GetNeighbors(testGeohash)
		if err != nil {
			assert.NoError(b, err)
		}
	}
}

func BenchmarkCodeforEncode(b *testing.B) {
	for b.Loop() {
		codefor_geohash.Encode(testPoint.Lat, testPoint.Lon)
	}
}

func BenchmarkCodeforDecode(b *testing.B) {
	for b.Loop() {
		codefor_geohash.Decode(testGeohash)
	}
}

func BenchmarkTomiHiltunenEncode(b *testing.B) {
	for b.Loop() {
		tomi_hiltunen_geohash.EncodeWithPrecision(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkTomiHiltunenDecode(b *testing.B) {
	for b.Loop() {
		tomi_hiltunen_geohash.Decode(testGeohash)
	}
}

func BenchmarkTomiHiltunenNeighbors(b *testing.B) {
	for b.Loop() {
		tomi_hiltunen_geohash.CalculateAllAdjacent(testGeohash)
	}
}

func BenchmarkBroadyEncode(b *testing.B) {
	for b.Loop() {
		broadygeohash.Encode(testPoint.Lat, testPoint.Lon)
	}
}

func BenchmarkBroadyDecode(b *testing.B) {
	for b.Loop() {
		broadygeohash.Decode(testGeohash)
	}
}

func BenchmarkFanixkEncode(b *testing.B) {
	for b.Loop() {
		fanixk_geohash.PrecisionEncode(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkFanixkDecode(b *testing.B) {
	for b.Loop() {
		fanixk_geohash.DecodeBoundingBox(testGeohash)
	}
}

func BenchmarkFanixkNeighbors(b *testing.B) {
	for b.Loop() {
		fanixk_geohash.Neighbors(testGeohash)
	}
}

func BenchmarkMmcloughlinEncode(b *testing.B) {
	for b.Loop() {
		mmcloughlin_geohash.EncodeWithPrecision(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkMmcloughlinDecode(b *testing.B) {
	for b.Loop() {
		mmcloughlin_geohash.BoundingBox(testGeohash)
	}
}

func BenchmarkMmcloughlinNeighbors(b *testing.B) {
	for b.Loop() {
		fanixk_geohash.Neighbors(testGeohash)
	}
}

func BenchmarkThe42CartconvertEncode(b *testing.B) {
	pc := &the42_cartconvert_geohash.PolarCoord{
		Latitude:  testPoint.Lat,
		Longitude: testPoint.Lon,
		Height:    0,
		El:        the42_cartconvert_geohash.DefaultEllipsoid,
	}
	precision := byte(testPrecision)
	for b.Loop() {
		the42_cartconvert_geohash.LatLongToGeoHashBits(pc, precision)
	}
}

func BenchmarkThe42CartconvertDecode(b *testing.B) {
	for b.Loop() {
		_, err := the42_cartconvert_geohash.GeoHashToLatLong(testGeohash, the42_cartconvert_geohash.DefaultEllipsoid)
		if err != nil {
			assert.NoError(b, err)
		}
	}
}
