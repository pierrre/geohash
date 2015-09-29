package geohash

import (
	"reflect"
	"testing"
)

var (
	testGeohash   = "u09tvqxnnuph"
	testNeighbors = Neighbors{
		North: "u09tvqxnnupj",
		NorthEast: "u09tvqxnnupm",
		East: "u09tvqxnnupk",
		SouthEast: "u09tvqxnnup7",
		South: "u09tvqxnnup5",
		SouthWest: "u09tvqxnnung",
		West: "u09tvqxnnunu",
		NorthWest: "u09tvqxnnunv",
	}
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

func TestNeighborsInvalidCharacter(t *testing.T) {
	_, err := GetNeighbors("é")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestNeighbors(t *testing.T) {
	neighbors, err := GetNeighbors(testGeohash)

	if err != nil {
		t.Fatal("err from neighbors should not be nil")
	}

	if !reflect.DeepEqual(neighbors, testNeighbors) {
		t.Fatal("failed to return the correct neighbors")
	}
}
