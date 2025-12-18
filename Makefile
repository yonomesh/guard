.PHONY: build clean fmt test

build:
	- go build -o ./build/guard ./cmd/guard

clean:
	- rm -rf build/*

lint:
	- gofmt -w .
	- golangci-lint run

test:
	- ./build/guard --version
	- ./build/guard --help

