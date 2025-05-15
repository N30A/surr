BINARY_NAME = surr

.PHONY: build install run

build:
	go build -o bin/$(BINARY_NAME) .

install:
	go install .

run:
	go run .
