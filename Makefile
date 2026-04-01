# Makefile for go-passgen
# Variables
APP_NAME ?= passgen
IMAGE_NAME ?= jurikolo/go-passgen
IMAGE_TAG ?= latest
PORT ?= 8080
CONTAINER_RUNTIME ?= podman

.PHONY: help build test docker-build docker-run docker-push clean

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) .

test: ## Run all tests
	@echo "Running tests..."
	go test ./...

docker-build: ## Build container image using $(CONTAINER_RUNTIME)
	@echo "Building container image $(IMAGE_NAME):$(IMAGE_TAG) using $(CONTAINER_RUNTIME)..."
	$(CONTAINER_RUNTIME) build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run: ## Run container locally
	@echo "Running container on port $(PORT)..."
	$(CONTAINER_RUNTIME) run --rm -p $(PORT):8080 --name $(APP_NAME) $(IMAGE_NAME):$(IMAGE_TAG)

docker-push: ## Push container image to registry
	@echo "Pushing $(IMAGE_NAME):$(IMAGE_TAG) to registry..."
	$(CONTAINER_RUNTIME) push $(IMAGE_NAME):$(IMAGE_TAG)

clean: ## Remove built binary and container images
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	@if $(CONTAINER_RUNTIME) images | grep -q "$(IMAGE_NAME)"; then \
		echo "Removing container image $(IMAGE_NAME):$(IMAGE_TAG)"; \
		$(CONTAINER_RUNTIME) rmi $(IMAGE_NAME):$(IMAGE_TAG) || true; \
	fi

# Convenience targets
all: test build docker-build ## Run tests, build binary, and build container image

local: build ## Build and run locally
	@echo "Running locally on port $(PORT)..."
	./$(APP_NAME)