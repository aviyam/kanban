BINARY_NAME=kanban-tui

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

tidy:
	go mod tidy

# Install dependencies and build
init: tidy build

.PHONY: build run test clean tidy init
