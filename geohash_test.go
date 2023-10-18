package geohash

import (
	"testing"

	"github.com/pierrre/assert"
)

const (
	testGeohash   = "u09tvqxnnuph"
	testPrecision = 12
)

var testPoint = Point{Lat: 48.86, Lon: 2.35}

func TestEncodeAuto(t *testing.T) {
	gh := EncodeAuto(testPoint.Lat, testPoint.Lon)
	assert.Equal(t, gh, testGeohash[:7])
}

func TestEncode(t *testing.T) {
	gh := Encode(testPoint.Lat, testPoint.Lon, testPrecision)
	assert.Equal(t, gh, testGeohash)
}

func TestEncodeMaxPrecision(t *testing.T) {
	gh := Encode(testPoint.Lat, testPoint.Lon, encodeMaxPrecision+1)
	assert.StringLen(t, gh, encodeMaxPrecision)
}

func TestDecode(t *testing.T) {
	box, err := Decode(testGeohash)
	assert.NoError(t, err)
	assert.True(t, testPointIsInsideBox(testPoint, box))
}

func TestDecodeInvalidCharacter(t *testing.T) {
	_, err := Decode("é")
	assert.Error(t, err)
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
	center := box.Center()
	assert.Equal(t, center, testPoint)
}

func TestBoxRound(t *testing.T) {
	box, err := Decode(testGeohash)
	assert.NoError(t, err)
	round := box.Round()
	assert.Equal(t, round, testPoint)
	assert.NotEqual(t, round, box.Center())
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
		t.Run(tc.gh, func(t *testing.T) {
			neighbors, err := GetNeighbors(tc.gh)
			assert.NoError(t, err)
			assert.Equal(t, neighbors, tc.expected)
		})
	}
}

func TestNeighborsInvalidCharacter(t *testing.T) {
	_, err := GetNeighbors("é")
	assert.Error(t, err)
}
