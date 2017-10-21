build: build/geohash

build/geohash:
	go build -v -i -o build/geohash ./cmd/geohash

test:
	go test -v ./...

lint:
	go get -v -u github.com/alecthomas/gometalinter
	gometalinter --install
	GOGC=800 gometalinter --enable-all -D dupl -D lll -D gas -D goconst -D gotype -D interfacer -D safesql -D test -D testify -D vetshadow\
	 --tests --warn-unmatched-nolint --deadline=10m --concurrency=4 --enable-gc ./...

clean:
	rm -rf build

.PHONY: build test lint clean
