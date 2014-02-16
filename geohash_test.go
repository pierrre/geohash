package geohash

import (
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

func BenchmarkEncode0(b *testing.B) {
	benchmarkEncode(b, 0)
}

func BenchmarkEncode1(b *testing.B) {
	benchmarkEncode(b, 1)
}

func BenchmarkEncode2(b *testing.B) {
	benchmarkEncode(b, 2)
}

func BenchmarkEncode3(b *testing.B) {
	benchmarkEncode(b, 3)
}

func BenchmarkEncode4(b *testing.B) {
	benchmarkEncode(b, 4)
}

func BenchmarkEncode5(b *testing.B) {
	benchmarkEncode(b, 5)
}

func BenchmarkEncode6(b *testing.B) {
	benchmarkEncode(b, 6)
}

func BenchmarkEncode7(b *testing.B) {
	benchmarkEncode(b, 7)
}

func BenchmarkEncode8(b *testing.B) {
	benchmarkEncode(b, 8)
}

func BenchmarkEncode9(b *testing.B) {
	benchmarkEncode(b, 9)
}

func BenchmarkEncode10(b *testing.B) {
	benchmarkEncode(b, 10)
}

func BenchmarkEncode11(b *testing.B) {
	benchmarkEncode(b, 11)
}

func BenchmarkEncode12(b *testing.B) {
	benchmarkEncode(b, 12)
}

func benchmarkEncode(b *testing.B, precision int) {
	for i := 0; i < b.N; i++ {
		Encode(testPoint.Lat, testPoint.Lon, precision)
	}
}

func BenchmarkDecode0(b *testing.B) {
	benchmarDecode(b, 0)
}

func BenchmarkDecode1(b *testing.B) {
	benchmarDecode(b, 1)
}

func BenchmarkDecode2(b *testing.B) {
	benchmarDecode(b, 2)
}

func BenchmarkDecode3(b *testing.B) {
	benchmarDecode(b, 3)
}

func BenchmarkDecode4(b *testing.B) {
	benchmarDecode(b, 4)
}

func BenchmarkDecode5(b *testing.B) {
	benchmarDecode(b, 5)
}

func BenchmarkDecode6(b *testing.B) {
	benchmarDecode(b, 6)
}

func BenchmarkDecode7(b *testing.B) {
	benchmarDecode(b, 7)
}

func BenchmarkDecode8(b *testing.B) {
	benchmarDecode(b, 8)
}

func BenchmarkDecode9(b *testing.B) {
	benchmarDecode(b, 9)
}

func BenchmarkDecode10(b *testing.B) {
	benchmarDecode(b, 10)
}

func BenchmarkDecode11(b *testing.B) {
	benchmarDecode(b, 11)
}

func BenchmarkDecode12(b *testing.B) {
	benchmarDecode(b, 12)
}

func benchmarDecode(b *testing.B, precision int) {
	gh := testGeohash[:precision]
	for i := 0; i < b.N; i++ {
		Decode(gh)
	}
}
