BINARY_NAME=bin/kanban-tui

build:
	mkdir -p bin
	go build -o $(BINARY_NAME) ./src

run: build
	./$(BINARY_NAME)

test:
	go test -v ./src/...

clean:
	go clean
	rm -f $(BINARY_NAME)

tidy:
	go mod tidy

# Install dependencies and build
init: tidy build

.PHONY: build run test clean tidy init
