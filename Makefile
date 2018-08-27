build: build/geohash

build/geohash:
	go build -v -i -o build/geohash ./cmd/geohash

test:
	go test -v ./...

lint:
	go get -v -u github.com/alecthomas/gometalinter
	gometalinter --install
	GOGC=800 gometalinter --disable-all -E deadcode -E errcheck -E gocyclo -E gofmt -E goimports -E golint -E ineffassign -E megacheck -E misspell -E nakedret -E structcheck -E unconvert -E unparam -E varcheck -E vet\
 --tests --vendor --warn-unmatched-nolint --deadline=10m --concurrency=4 --enable-gc ./...

clean:
	rm -rf build

.PHONY: build test lint clean
