BINARY := ragcheck

.PHONY: build test example release-check snapshot

build:
	mkdir -p dist
	go build -o dist/$(BINARY) .

test:
	go test ./...

example:
	go run . score -qrels examples/qrels.json -run examples/run.json -k 3

release-check:
	goreleaser check

snapshot:
	goreleaser release --snapshot --clean

