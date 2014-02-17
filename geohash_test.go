package geohash

import (
	codefor_geohash "github.com/Codefor/geohash"
	broady_geohash "github.com/broady/gogeohash"
	//gnagel_geohash "github.com/gnagel/go-geohash/ggeohash"
	the42_cartconvert_geohash "github.com/the42/cartconvert/cartconvert"
	"testing"
)

var (
	testGeohash   = "u09tvqxnnuph"
	testPoint     = Point{Lat: 48.86, Lon: 2.35}
	testPrecision = 12
)

func TestEncodeAuto(t *testing.T) {
	gh := EncodeAuto(testPoint.Lat, testPoint.Lon)
	if gh != testGeohash[:7] {
		t.Fatal("wrong geohash")
	}
}

func TestEncode(t *testing.T) {
	gh := Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	if gh != testGeohash {
		t.Fatal("wrong geohash")
	}
}

func TestEncodeNotThreadSafe(t *testing.T) {
	gh := EncodeNotThreadSafe(testPoint.Lat, testPoint.Lon, testPrecision)
	if gh != testGeohash {
		t.Fatal("wrong geohash")
	}
}

func TestDecode(t *testing.T) {
	box, err := Decode(testGeohash)
	if err != nil {
		t.Fatal(err)
	}
	if !(testPointIsInsideBox(testPoint, box)) {
		t.Fatal("point is outside")
	}
}

func TestDecodeInvalidCharacter(t *testing.T) {
	_, err := Decode("é")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestBoxCenter(t *testing.T) {
	box := Box{
		Lat: Range{
			Min: testPoint.Lat - 1,
			Max: testPoint.Lat + 1,
		},
		Lon: Range{
			Min: testPoint.Lon - 1,
			Max: testPoint.Lon + 1,
		},
	}
	if box.Center() != testPoint {
		t.Fatal("invalid center")
	}
}

func TestBoxRound(t *testing.T) {
	box, err := Decode(testGeohash)
	if err != nil {
		t.Fatal(err)
	}
	round := box.Round()
	if round != testPoint {
		t.Fatal("invalid round")
	}
	if round == box.Center() {
		t.Fatal("round is equal to center")
	}
}

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

func testPointIsInsideBox(p Point, b Box) bool {
	return testValueIsInsideRange(p.Lat, b.Lat) &&
		testValueIsInsideRange(p.Lon, b.Lon)
}

func testValueIsInsideRange(v float64, r Range) bool {
	return v >= r.Min && v <= r.Max
}
