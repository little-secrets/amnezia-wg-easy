.PHONY: build build-wgpw test clean docker-build docker-run help

# Default target
.DEFAULT_GOAL := help

# Build main application
build:
	@echo "Building amnezia-wg-easy..."
	@go build -ldflags="-s -w" -o amnezia-wg-easy .

# Build wgpw utility
build-wgpw:
	@echo "Building wgpw..."
	@go build -ldflags="-s -w" -o wgpw ./cmd/wgpw

# Build all
build-all: build build-wgpw

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f amnezia-wg-easy wgpw
	@rm -f amnezia-wg-easy.exe wgpw.exe

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t amnezia-wg-easy:latest .

# Docker run (development)
docker-run:
	@echo "Starting container..."
	@docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Generate password hash
gen-password:
	@go run ./cmd/wgpw

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build main application"
	@echo "  build-wgpw     - Build wgpw utility"
	@echo "  build-all      - Build all binaries"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with docker-compose"
	@echo "  gen-password   - Generate bcrypt password hash"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  help           - Show this help"

