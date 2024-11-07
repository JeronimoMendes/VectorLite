BINARY_NAME=vectorlite


all: build

build:
	go build -ldflags=-w -o bin/$(BINARY_NAME) ./cmd/vectorlite

test:
	go test ./... -v

run: build
	./bin/vectorlite

clean:
	rm -rf bin

.PHONY: build run clean
