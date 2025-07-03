.PHONY: test build lint fmt vet clean install-deps tidy

# Default target
all: test build

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Build the project
build:
	go build -v ./...

# Run linting
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	go clean ./...
	rm -f coverage.out coverage.html

# Install dependencies
install-deps:
	go mod download
	go mod verify

# Tidy dependencies
tidy:
	go mod tidy

# Generate coverage report
coverage: test
	go tool cover -html=coverage.out -o=coverage.html

# Run all checks
check: fmt vet lint test

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest 