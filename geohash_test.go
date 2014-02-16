package geohash

import (
	codefor_geohash "github.com/Codefor/geohash"
	broady_geohash "github.com/broady/gogeohash"
	gnagel_geohash "github.com/gnagel/go-geohash/ggeohash"
	the42_cartconvert_geohash "github.com/the42/cartconvert/cartconvert"
	"testing"
)

var (
	testGeohash   = "u09tvqxnnuph"
	testPoint     = Point{Lat: 48.86, Lon: 2.35}
	testPrecision = 12
)

func TestEncode(t *testing.T) {
	gh := Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	if gh != testGeohash {
		t.Fatal("wrong geohash")
	}
}

func TestDecode(t *testing.T) {
	box, err := Decode(testGeohash)
	if err != nil {
		t.Fatal(err)
	}
	if !box.IsPointInside(testPoint) {
		t.Fatal("point is outside")
	}
}

func TestDecodeInvalidCharacter(t *testing.T) {
	_, err := Decode("é")
	if err == nil {
		t.Fatal("no error")
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode(testGeohash)
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
