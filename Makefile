build: build-cmd

build-cmd:
	go build -v -i -o build/geohash ./cmd/geohash

test:
	go test -v ./...

lint:
	go get -v -u github.com/alecthomas/gometalinter
	gometalinter --install --update --no-vendored-linters
	GOGC=800 gometalinter --enable-all -D dupl -D lll -D gas -D goconst -D gotype -D interfacer -D misspell -D safesql -D test -D testify -D vetshadow\
	 --tests --deadline=10m --concurrency=4 --enable-gc ./...

clean:
	rm -rf build

.PHONY: build build-cmd test lint clean
