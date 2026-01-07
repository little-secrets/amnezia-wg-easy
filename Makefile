.PHONY: build build-wgpw test clean docker-build docker-run docker-login docker-registry-build docker-registry-build-multiplatform docker-push docker-release docker-release-multiplatform docker-tag-latest help

# Default target
.DEFAULT_GOAL := help

# Docker configuration
DOCKER_REGISTRY ?= docker.highres.team
DOCKER_IMAGE_NAME ?= amnezia/go-amnezia-wg-easy
DOCKER_TAG ?= latest
LOCAL_IMAGE ?= go-amnezia-wg-easy

# Full image names
LOCAL_IMAGE_FULL = $(LOCAL_IMAGE):$(DOCKER_TAG)
REGISTRY_IMAGE_FULL = $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Build main application
build:
	@echo "Building amnezia-wg-easy..."
	@go build -buildvcs=false -ldflags="-s -w" -o amnezia-wg-easy .

# Build wgpw utility
build-wgpw:
	@echo "Building wgpw..."
	@go build -buildvcs=false -ldflags="-s -w" -o wgpw ./cmd/wgpw

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

# Docker build (local)
docker-build:
	@echo "Building local Docker image: $(LOCAL_IMAGE_FULL)"
	@docker build --no-cache -t $(LOCAL_IMAGE_FULL) .

# Docker login to registry
docker-login:
	@echo "Logging in to $(DOCKER_REGISTRY)..."
	@docker login $(DOCKER_REGISTRY)

# Docker build for registry
docker-registry-build:
	@echo "Building Docker image for registry: $(REGISTRY_IMAGE_FULL)"
	@docker build --no-cache -t $(REGISTRY_IMAGE_FULL) .

# Docker build multi-platform for registry
docker-registry-build-multiplatform:
	@echo "Building multi-platform image for registry: $(REGISTRY_IMAGE_FULL)"
	@docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(REGISTRY_IMAGE_FULL) \
		--push .

# Docker push to registry
docker-push:
	@echo "Pushing to registry: $(REGISTRY_IMAGE_FULL)"
	@docker push $(REGISTRY_IMAGE_FULL)

# Docker build and push (release)
docker-release: docker-registry-build docker-push
	@echo "Successfully released $(REGISTRY_IMAGE_FULL)"

# Docker multi-platform release
docker-release-multiplatform:
	@echo "Building and pushing multi-platform image..."
	@$(MAKE) docker-registry-build-multiplatform
	@echo "Successfully released multi-platform $(REGISTRY_IMAGE_FULL)"

# Docker tag and push latest + version
docker-tag-latest:
	@if [ "$(DOCKER_TAG)" != "latest" ]; then \
		echo "Tagging $(REGISTRY_IMAGE_FULL) as latest..."; \
		docker tag $(REGISTRY_IMAGE_FULL) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):latest; \
		docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):latest; \
		echo "Successfully tagged and pushed latest"; \
	else \
		echo "Skipping: Already using 'latest' tag"; \
	fi

# Docker run (development)
docker-run:
	@echo "Starting container..."
	@docker compose -f docker-compose.yml up --build

# Generate password hash
gen-password:
	@go run -buildvcs=false ./cmd/wgpw

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
	@echo ""
	@echo "Build:"
	@echo "  build                 - Build main application"
	@echo "  build-wgpw            - Build wgpw utility"
	@echo "  build-all             - Build all binaries"
	@echo ""
	@echo "Docker (Local):"
	@echo "  docker-build          - Build local Docker image ($(LOCAL_IMAGE_FULL))"
	@echo "  docker-run            - Run with docker-compose"
	@echo ""
	@echo "Docker (Registry):"
	@echo "  docker-login                      - Login to Docker registry"
	@echo "  docker-registry-build             - Build image for registry"
	@echo "  docker-registry-build-multiplatform - Build multi-platform image (amd64+arm64)"
	@echo "  docker-push                       - Push image to registry"
	@echo "  docker-release                    - Build and push to registry"
	@echo "  docker-release-multiplatform      - Build and push multi-platform image"
	@echo "  docker-tag-latest                 - Tag current image as latest and push"
	@echo ""
	@echo "Development:"
	@echo "  test                  - Run tests"
	@echo "  clean                 - Remove build artifacts"
	@echo "  gen-password          - Generate bcrypt password hash"
	@echo "  deps                  - Download and tidy dependencies"
	@echo "  fmt                   - Format code"
	@echo "  lint                  - Lint code"
	@echo ""
	@echo "Docker Registry Configuration:"
	@echo "  DOCKER_REGISTRY       - Registry URL (default: $(DOCKER_REGISTRY))"
	@echo "  DOCKER_IMAGE_NAME     - Image name (default: $(DOCKER_IMAGE_NAME))"
	@echo "  DOCKER_TAG            - Image tag (default: $(DOCKER_TAG))"
	@echo ""
	@echo "Examples:"
	@echo "  # Build and push to default registry"
	@echo "  make docker-release"
	@echo ""
	@echo "  # Build and push with custom tag"
	@echo "  make docker-release DOCKER_TAG=v1.2.0"
	@echo ""
	@echo "  # Build and push to custom registry"
	@echo "  make docker-release DOCKER_REGISTRY=myregistry.com DOCKER_TAG=v1.2.0"
	@echo ""
	@echo "  # Build and push multi-platform"
	@echo "  make docker-release-multiplatform DOCKER_TAG=v1.2.0"
	@echo ""
	@echo "  # Login to registry first"
	@echo "  make docker-login && make docker-release"

