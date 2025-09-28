.PHONY: build clean test install

# Binary name
BINARY_NAME=institutionalized

# Version - can be overridden by environment variable or command line
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0-dev")

# Build the binary
build:
	go build -ldflags="-X github.com/IanKnighton/institutionalized/cmd.Version=$(VERSION)" -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Install the binary to GOPATH/bin
install:
	go install -ldflags="-X github.com/IanKnighton/institutionalized/cmd.Version=$(VERSION)" .

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy modules
tidy:
	go mod tidy

# Run all checks
check: fmt tidy test

# Default target
all: clean build