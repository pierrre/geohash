package geohash

import (
	"testing"

	codefor_geohash "github.com/Codefor/geohash"
	tomi_hiltunen_geohash "github.com/TomiHiltunen/geohash-golang"
	broady_geohash "github.com/broady/gogeohash" //nolint: misspell
	fanixk_geohash "github.com/fanixk/geohash"
	mmcloughlin_geohash "github.com/mmcloughlin/geohash"
	the42_cartconvert_geohash "github.com/the42/cartconvert/cartconvert"
)

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Decode(testGeohash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNeighbors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetNeighbors(testGeohash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCodeforEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		codefor_geohash.Encode(testPoint.Lat, testPoint.Lon)
	}
}

func BenchmarkCodeforDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		codefor_geohash.Decode(testGeohash)
	}
}

func BenchmarkTomiHiltunenEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tomi_hiltunen_geohash.EncodeWithPrecision(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkTomiHiltunenDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tomi_hiltunen_geohash.Decode(testGeohash)
	}
}

func BenchmarkTomiHiltunenNeighbors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tomi_hiltunen_geohash.CalculateAllAdjacent(testGeohash)
	}
}

func BenchmarkBroadyEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		broady_geohash.Encode(testPoint.Lat, testPoint.Lon) //nolint: misspell
	}
}

func BenchmarkBroadyDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		broady_geohash.Decode(testGeohash) //nolint: misspell
	}
}

func BenchmarkFanixkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fanixk_geohash.PrecisionEncode(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkFanixkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fanixk_geohash.DecodeBoundingBox(testGeohash)
	}
}

func BenchmarkFanixkNeighbors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fanixk_geohash.Neighbors(testGeohash)
	}
}

func BenchmarkMmcloughlinEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mmcloughlin_geohash.EncodeWithPrecision(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkMmcloughlinDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mmcloughlin_geohash.BoundingBox(testGeohash)
	}
}

func BenchmarkMmcloughlinNeighbors(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
		the42_cartconvert_geohash.LatLongToGeoHashBits(pc, precision)
	}
}

func BenchmarkThe42CartconvertDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := the42_cartconvert_geohash.GeoHashToLatLong(testGeohash, the42_cartconvert_geohash.DefaultEllipsoid)
		if err != nil {
			b.Fatal(err)
		}
	}
}
