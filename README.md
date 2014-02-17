# Geohash

A geohash library for Go (Golang)

## Features
- Encode latitude/longitude to geohash (optional precision)
- Decode geohash to latitude/longitude
- Round a geohash box to a single location

## Documentation
http://godoc.org/github.com/pierrre/geohash

## Benchmark and comparison
```
go test -v -bench=. -benchmem

BenchmarkEncode-8	 5000000	       352 ns/op	      32 B/op	       2 allocs/op
BenchmarkEncodeNotThreadSafe-8	10000000	       272 ns/op	      16 B/op	       1 allocs/op
BenchmarkDecode-8	 5000000	       324 ns/op	       0 B/op	       0 allocs/op

BenchmarkCodeforEncode-8	 5000000	       623 ns/op	      64 B/op	       6 allocs/op
BenchmarkCodeforDecode-8	 5000000	       330 ns/op	       0 B/op	       0 allocs/op

BenchmarkGnagelEncode-8	 5000000	       519 ns/op	     128 B/op	       2 allocs/op
BenchmarkGnagelDecode-8	 5000000	       563 ns/op	      32 B/op	       1 allocs/op

BenchmarkThe42CartconvertEncode-8	 1000000	      1555 ns/op	     320 B/op	      23 allocs/op
BenchmarkThe42CartconvertDecode-8	 2000000	       758 ns/op	      49 B/op	       2 allocs/op

BenchmarkBroadyEncode-8	 1000000	      1646 ns/op	     320 B/op	      23 allocs/op
BenchmarkBroadyDecode-8	 1000000	      1512 ns/op	     192 B/op	      12 allocs/op

CPU: Intel Core i7-3610QM
RAM: 8GB
```
