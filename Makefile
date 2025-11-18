.PHONY: build test clean lint run deps docker

# Build the application
build:
	go build -o bin/hulk ./cmd

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run the application
run:
	go run ./cmd/main.go

# Build Docker image
docker:
	docker build -t mcp-hulk:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 mcp-hulk:latest

# Generate code
generate:
	go generate ./...

# Format code
fmt:
	go fmt ./...

# Vet code for potential issues
vet:
	go vet ./...

# Run all checks (lint, vet, test)
check: lint vet test

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/hulk-linux ./cmd

# Install application locally
install:
	go install ./cmd

# Development setup
dev-setup: deps
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help target
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Run linter"
	@echo "  deps          - Install dependencies"
	@echo "  run           - Run the application"
	@echo "  docker        - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  generate      - Generate code"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code for potential issues"
	@echo "  check         - Run all checks (lint, vet, test)"
	@echo "  build-prod    - Build for production"
	@echo "  install       - Install application locally"
	@echo "  dev-setup     - Development setup"
	@echo "  help          - Show this help message"