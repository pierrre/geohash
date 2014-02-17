package geohash

import (
	codefor_geohash "github.com/Codefor/geohash"
	broady_geohash "github.com/broady/gogeohash"
	//gnagel_geohash "github.com/gnagel/go-geohash/ggeohash"
	the42_cartconvert_geohash "github.com/the42/cartconvert/cartconvert"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkEncodeNotThreadSafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeNotThreadSafe(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode(testGeohash)
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

//This test is disabled because the project contains a directory (example-scripts) with a wrong import (./go_geohash)
/*
func BenchmarkGnagelEncode(b *testing.B) {
	precision := uint8(testPrecision)
	for i := 0; i < b.N; i++ {
		gnagel_geohash.Encode(testPoint.Lat, testPoint.Lon, precision)
	}
}

func BenchmarkGnagelDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gnagel_geohash.DecodeBoundBox(testGeohash)
	}
}
*/

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
		the42_cartconvert_geohash.GeoHashToLatLong(testGeohash, the42_cartconvert_geohash.DefaultEllipsoid)
	}
}

func BenchmarkBroadyEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		broady_geohash.Encode(testPoint.Lat, testPoint.Lon)
	}
}

func BenchmarkBroadyDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		broady_geohash.Decode(testGeohash)
	}
}
