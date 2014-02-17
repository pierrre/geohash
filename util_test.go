package geohash

import (
	"testing"
)

func TestRound(t *testing.T) {
	for _, test := range []struct {
		value    float64
		expected float64
	}{
		{-1.7, -2},
		{-1.5, -2},
		{-1.3, -1},
		{-1, -1},
		{0, 0},
		{1, 1},
		{1.3, 1},
		{1.5, 2},
		{1.7, 2},
	} {
		if result := round(test.value); result != test.expected {
			t.Fatalf("round %f, got %f instead of %f", test.value, result, test.expected)
		}
	}
}

func TestRoundDecimal(t *testing.T) {
	value := 12.345678
	for _, test := range []struct {
		dec      int
		expected float64
	}{
		{0, 12},
		{1, 12.3},
		{2, 12.35},
		{3, 12.346},
		{4, 12.3457},
		{5, 12.34568},
		{6, 12.345678},
	} {
		if result := roundDecimal(value, test.dec); result != test.expected {
			t.Fatalf("round %f with %d decimal, got %f instead of %f", value, test.dec, result, test.expected)
		}
	}
}
