package geohash

import (
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

func testPointIsInsideBox(p Point, b Box) bool {
	return testValueIsInsideRange(p.Lat, b.Lat) &&
		testValueIsInsideRange(p.Lon, b.Lon)
}

func testValueIsInsideRange(v float64, r Range) bool {
	return v >= r.Min && v <= r.Max
}
