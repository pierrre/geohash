package geohash

import (
	"testing"
)

func TestDecode(t *testing.T) {
	box, err := Decode("ezs42")
	if err != nil {
		t.Fatal(err)
	}
	if !box.IsPointInside(Point{Lat: 42.6, Lon: -5.6}) {
		t.Fatal("decoding error: point is outside")
	}
}

func TestDecodeInvalidRune(t *testing.T) {
	_, err := Decode("é")
	if err == nil {
		t.Fatal("no error")
	}
}
