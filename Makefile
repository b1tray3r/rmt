.PHONY: all tidy build run clean docker docker-local docker-build docker-push dev run-docker docker-scan docker-inspect docker-clean version help

VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VCS_REF ?= $(shell git rev-parse HEAD)
DOCKER_IMAGE ?= aborgardt/rmt

all: tidy test build run clean

tidy:
	@echo "[tidy]"
	@go mod tidy

test:
	@echo "Running tests..."
	@go test ./... -cover
	@mkdir -p build
	@go test -coverprofile=build/coverage.out ./... > /dev/null
	@if [ -f build/coverage.out ]; then \
		go tool cover -func=build/coverage.out | grep total | awk '{print "Total coverage: " $$3}'; \
		rm -f build/coverage.out; \
	fi

build:
	@echo "[build]"
	@mkdir -p build
	@CGO_ENABLED=1 go build \
		-ldflags="-s -w -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitCommit=$(VCS_REF)" \
		-o ./build/rmt .

docker:
	@echo "[docker-local] Building local development image..."
	@echo "$(VERSION)-local" > VERSION
	@docker build \
		--build-arg VERSION="$(VERSION)-local" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg VCS_REF="$(VCS_REF)" \
		-t $(DOCKER_IMAGE):local \
		-t $(DOCKER_IMAGE):$(VERSION)-local \
		. --no-cache
	@rm -f VERSION

run:
	@echo "[run]"
	@./build/rmt

clean:
	@echo "[clean]"
	@rm -rf build
	@rm -f VERSION

# Development and testing targets
dev: tidy test build
	@echo "[dev] Development build complete"

run-docker:
	@echo "[run-docker] Running local Docker image..."
	@docker run --rm -it $(DOCKER_IMAGE):local

# Docker utility targets
docker-scan:
	@echo "[docker-scan] Scanning image for vulnerabilities..."
	@docker scout cves $(DOCKER_IMAGE):local || echo "Docker Scout not available, try: trivy image $(DOCKER_IMAGE):local"

docker-inspect:
	@echo "[docker-inspect] Inspecting image..."
	@docker inspect $(DOCKER_IMAGE):local | jq '.[0].Config.Labels'

docker-clean:
	@echo "[docker-clean] Cleaning up Docker images..."
	@docker images $(DOCKER_IMAGE) -q | xargs -r docker rmi -f

# Version and info targets
version:
	@echo "Version: $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "VCS Ref: $(VCS_REF)"
	@echo "Docker Image: $(DOCKER_IMAGE)"

help:
	@echo "Available targets:"
	@echo "  all           - Run tidy, test, build, run, clean"
	@echo "  dev           - Run tidy, test, build (development workflow)"
	@echo "  tidy          - Clean up go.mod"
	@echo "  test          - Run tests with coverage"
	@echo "  build         - Build binary with version info"
	@echo "  run           - Run the built binary"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-local  - Build local development Docker image"
	@echo "  docker-build  - Build multi-platform Docker image"
	@echo "  docker-push   - Build and push multi-platform Docker image"
	@echo "  run-docker    - Run local Docker image"
	@echo "  docker-scan   - Scan Docker image for vulnerabilities"
	@echo "  docker-inspect- Inspect Docker image labels"
	@echo "  docker-clean  - Remove Docker images"
	@echo ""
	@echo "Utility targets:"
	@echo "  version       - Show version information"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Variables (can be overridden):"
	@echo "  VERSION=$(VERSION)"
	@echo "  DOCKER_IMAGE=$(DOCKER_IMAGE)"