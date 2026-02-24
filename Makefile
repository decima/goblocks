.PHONY: help build run test test-verbose test-coverage clean install dev fmt vet lint

# Variables
BINARY_NAME=goblocks
VERSION?=0.0.1
BUILD_DIR=./bin
GO=go
GOFLAGS=-v

# Default target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

build-mac: ## Build for macOS
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64"

build-all: build-linux build-windows build-mac ## Build for all platforms

##@ Development

dev: ## Run in development mode with hot reload
	@echo "Starting in development mode..."
	DEBUG=1 $(GO) run .

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

install: ## Install dependencies
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies installed"

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GO) fmt ./app/... ./libraries/... .
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./app/... ./libraries/... .
	@echo "Vet complete"

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

##@ Testing

test: ## Run tests
	@echo "Running tests..."
	$(GO) test ./app/...
	@echo "Tests complete"

test-verbose: ## Run tests in verbose mode
	@echo "Running tests (verbose)..."
	$(GO) test -v ./app/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GO) test -cover ./app/...
	@echo ""
	@echo "Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./app/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GO) test -race ./app/...

test-bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./app/...

##@ Docker

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "Docker image built: $(BINARY_NAME):$(VERSION)"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8000:8000 -v $(PWD)/data:/data $(BINARY_NAME):latest

##@ Cleanup

clean: ## Remove build artifacts and cache
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean -cache
	@echo "Clean complete"

clean-data: ## Remove data directory (WARNING: deletes all blocks)
	@echo "WARNING: This will delete all block data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		rm -rf ./data/*; \
		echo "Data directory cleaned"; \
	else \
		echo "Cancelled"; \
	fi

##@ Utilities

check: fmt vet test ## Run formatting, vetting, and tests

ci: install check ## Run CI pipeline (install, format, vet, test)

info: ## Display project information
	@echo "Project: $(BINARY_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Go version: $$($(GO) version)"
	@echo "Build dir: $(BUILD_DIR)"

deps-upgrade: ## Upgrade all dependencies
	@echo "Upgrading dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "Dependencies upgraded"

deps-graph: ## Show dependency graph (requires graphviz)
	@echo "Generating dependency graph..."
	@if command -v dot >/dev/null 2>&1; then \
		$(GO) mod graph | sed 's/@[^ ]*//g' | sort -u | \
		awk '{print "\"" $$1 "\" -> \"" $$2 "\""}' | \
		dot -Tpng -o deps.png; \
		echo "Dependency graph saved to deps.png"; \
	else \
		echo "graphviz not installed. Install it with:"; \
		echo "  sudo apt-get install graphviz  # Ubuntu/Debian"; \
		echo "  brew install graphviz          # macOS"; \
	fi
