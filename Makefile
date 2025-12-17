.PHONY: build clean fmt test

build:
	- go build -o ./build/guard ./cmd/guard

clean:
	- rm -rf build/*

fmt:
	gofmt -w .

test:
	- ./build/guard --version
	- ./build/guard --help

