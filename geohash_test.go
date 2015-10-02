package geohash

import (
	"reflect"
	"testing"
)

const (
	testGeohash   = "u09tvqxnnuph"
	testPrecision = 12
)

var (
	testPoint = Point{Lat: 48.86, Lon: 2.35}
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

func TestNeighbors(t *testing.T) {
	for _, tc := range []struct {
		gh       string
		expected Neighbors
	}{
		{
			gh: testGeohash,
			expected: Neighbors{
				North:     "u09tvqxnnupj",
				NorthEast: "u09tvqxnnupm",
				East:      "u09tvqxnnupk",
				SouthEast: "u09tvqxnnup7",
				South:     "u09tvqxnnup5",
				SouthWest: "u09tvqxnnung",
				West:      "u09tvqxnnunu",
				NorthWest: "u09tvqxnnunv",
			},
		},
		{
			gh: Encode(0, 0, 4), // s000
			expected: Neighbors{
				North:     "s001",
				NorthEast: "s003",
				East:      "s002",
				SouthEast: "kpbr",
				South:     "kpbp",
				SouthWest: "7zzz",
				West:      "ebpb",
				NorthWest: "ebpc",
			},
		},
		{
			gh: Encode(0, 180, 4), // xbpb
			expected: Neighbors{
				North:     "xbpc",
				NorthEast: "8001",
				East:      "8000",
				SouthEast: "2pbp",
				South:     "rzzz",
				SouthWest: "rzzx",
				West:      "xbp8",
				NorthWest: "xbp9",
			},
		},
		{
			gh: Encode(90, 0, 4), // upbp
			expected: Neighbors{
				North:     "bpbp",
				NorthEast: "bpbr",
				East:      "upbr",
				SouthEast: "upbq",
				South:     "upbn",
				SouthWest: "gzzy",
				West:      "gzzz",
				NorthWest: "zzzz",
			},
		},
	} {
		neighbors, err := GetNeighbors(tc.gh)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(neighbors, tc.expected) {
			t.Fatalf("unexpected result for %s: got %v, want %v", tc.gh, neighbors, tc.expected)
		}
	}
}

func TestNeighborsInvalidCharacter(t *testing.T) {
	_, err := GetNeighbors("é")
	if err == nil {
		t.Fatal("no error")
	}
}
