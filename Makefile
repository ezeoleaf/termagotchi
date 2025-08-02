.PHONY: build run clean test help

# Build the application
build:
	go build -o bin/termagotchi cmd/termagotchi/main.go

# Run the application
run: build
	./bin/termagotchi

# Clean build artifacts
clean:
	rm -f bin/termagotchi

lint:
	golangci-lint run

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy

# Show help
help:
	@echo "Available commands:"
	@echo "  build  - Build the termagotchi application"
	@echo "  run    - Build and run the application"
	@echo "  clean  - Remove build artifacts"
	@echo "  test   - Run tests"
	@echo "  lint   - Run linter"
	@echo "  deps   - Install dependencies"
	@echo "  help   - Show this help message"

# Default target
all: build 