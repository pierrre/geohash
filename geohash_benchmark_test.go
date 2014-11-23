package geohash

import (
	"testing"

	codefor_geohash "github.com/Codefor/geohash"
	tomi_hiltunen_geohash "github.com/TomiHiltunen/geohash-golang"
	broady_geohash "github.com/broady/gogeohash"
	the42_cartconvert_geohash "github.com/the42/cartconvert/cartconvert"
)

func BenchmarkEncode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Encode(testPoint.Lat, testPoint.Lon, testPrecision)
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Decode(testGeohash)
		}
	})
}

func BenchmarkCodeforEncode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			codefor_geohash.Encode(testPoint.Lat, testPoint.Lon)
		}
	})
}

func BenchmarkCodeforDecode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			codefor_geohash.Decode(testGeohash)
		}
	})
}

func BenchmarkTomiHiltunenEncode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tomi_hiltunen_geohash.Encode(testPoint.Lat, testPoint.Lon)
		}
	})
}

func BenchmarkTomiHiltunenDecode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tomi_hiltunen_geohash.Decode(testGeohash)
		}
	})
}

func BenchmarkThe42CartconvertEncode(b *testing.B) {
	pc := &the42_cartconvert_geohash.PolarCoord{
		Latitude:  testPoint.Lat,
		Longitude: testPoint.Lon,
		Height:    0,
		El:        the42_cartconvert_geohash.DefaultEllipsoid,
	}
	precision := byte(testPrecision)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			the42_cartconvert_geohash.LatLongToGeoHashBits(pc, precision)
		}
	})
}

func BenchmarkThe42CartconvertDecode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			the42_cartconvert_geohash.GeoHashToLatLong(testGeohash, the42_cartconvert_geohash.DefaultEllipsoid)
		}
	})
}

func BenchmarkBroadyEncode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			broady_geohash.Encode(testPoint.Lat, testPoint.Lon)
		}
	})
}

func BenchmarkBroadyDecode(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			broady_geohash.Decode(testGeohash)
		}
	})
}
