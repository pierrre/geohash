package geohash

import (
	"testing"
)

func TestRound(t *testing.T) {
	for _, tc := range []struct {
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
		if result := round(tc.value); result != tc.expected {
			t.Fatalf("round %f, got %f instead of %f", tc.value, result, tc.expected)
		}
	}
}

func TestRoundDecimal(t *testing.T) {
	value := 12.345678
	for _, tc := range []struct {
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
		if result := roundDecimal(value, tc.dec); result != tc.expected {
			t.Fatalf("round %f with %d decimal, got %f instead of %f", value, tc.dec, result, tc.expected)
		}
	}
}

func TestNormalize(t *testing.T) {
	for _, tc := range []struct {
		lat, lon                 float64
		expectedLat, expectedLon float64
	}{
		{testPoint.Lat, testPoint.Lon, testPoint.Lat, testPoint.Lon},
		{0, 0, 0, 0},
		{45, 90, 45, 90},
		{-45, -90, -45, -90},
		{90, 0, 90, 0},
		{-90, 0, -90, 0},
		{0, 180, 0, 180},
		{1, -180, 1, 180},
		{91, 0, 89, 180},
		{91, 1, 89, -179},
		{-91, -1, -89, 179},
		{0, 181, 0, -179},
		{0, -181, 0, 179},
		{270, 1, -90, 1},
	} {
		if resultLat, resultLon := normalize(tc.lat, tc.lon); resultLat != tc.expectedLat || resultLon != tc.expectedLon {
			t.Fatalf("unexpected result for %f,%f: got %f,%f, want %f,%f", tc.lat, tc.lon, resultLat, resultLon, tc.expectedLat, tc.expectedLon)
		}
	}
}

func TestCenterMod(t *testing.T) {
	for _, tc := range []struct {
		value    float64
		expected float64
	}{
		{0, 0},
		{45, 45},
		{180, 180},
		{181, -179},
		{-181, 179},
		{-180, 180},
		{-45, -45},
	} {
		if result := centerMod(tc.value, 360); result != tc.expected {
			t.Fatalf("unexpected result for %f: got %f, want %f", tc.value, result, tc.expected)
		}
	}
}
