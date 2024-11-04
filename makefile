BINARY_NAME=vectorlite


all: build

build:
	go build -ldflags=-w -o bin/$(BINARY_NAME) ./cmd/vectorlite

run: build
	./bin/vectorlite

clean:
	rm -rf bin

.PHONY: build run clean
