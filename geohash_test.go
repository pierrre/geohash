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

func TestHeight(t *testing.T) {
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

	height := box.Height()

	if height != 2 {
		t.Fatal("height is not calculated correctly")
	}
}

func TestWidth(t *testing.T) {
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

	width := box.Width()

	if width != 2 {
		t.Fatal("width is not calculated correctly")
	}
}

func TestNeighborsInvalidCharacter(t *testing.T) {
	_, err := Neighbors("é")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestNeighbors(t *testing.T) {
	neighbors, err := Neighbors(testGeohash)

	if err != nil {
		t.Fatal("err from neighbors should not be nil")
	}

	if len(neighbors) != 8 {
		t.Fatal("return the wrong number of neighbors")
	}

	if neighbors[0] != "u09tvqxnnupj" {
		t.Fatal("failed to return the correct north neighbor:", neighbors[1])
	}

	if neighbors[1] != "u09tvqxnnupm" {
		t.Fatal("failed to return the correct northeast neighbor:", neighbors[1])
	}

	if neighbors[2] != "u09tvqxnnupk" {
		t.Fatal("failed to return the correct east neighbor:", neighbors[2])
	}

	if neighbors[3] != "u09tvqxnnup7" {
		t.Fatal("failed to return the correct southeast neighbor:", neighbors[3])
	}

	if neighbors[4] != "u09tvqxnnup5" {
		t.Fatal("failed to return the correct south neighbor:", neighbors[4])
	}

	if neighbors[5] != "u09tvqxnnung" {
		t.Fatal("failed to return the correct southwest neighbor:", neighbors[5])
	}

	if neighbors[6] != "u09tvqxnnunu" {
		t.Fatal("failed to return the correct west neighbor:", neighbors[6])
	}

	if neighbors[7] != "u09tvqxnnunv" {
		t.Fatal("failed to return the correct northwest neighbor:", neighbors[7])
	}
}
