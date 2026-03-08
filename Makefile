.PHONY: build install test clean run fmt vet help

BINARY_NAME=envoke
VERSION?=0.1.0
BUILD_DIR=./bin
GOFLAGS=-ldflags "-X main.version=$(VERSION)"

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo " Building $(BINARY_NAME)..."
	@go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo " Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: ## Install the binary to $GOPATH/bin
	@echo " Installing $(BINARY_NAME)..."
	@go install $(GOFLAGS)
	@echo " Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

# test: ## Run tests
# 	@echo " Running tests..."
# 	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo " Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo " Coverage report: coverage.html"

run: build ## Build and run the binary
	@$(BUILD_DIR)/$(BINARY_NAME)

fmt: ## Format code
	@echo " Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo " Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo " Running golangci-lint..."
	@golangci-lint run

tidy: ## Run go mod tidy
	@echo " Tidying dependencies..."
	@go mod tidy

clean: ## Remove build artifacts
	@echo " Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo " Clean complete"

dev: ## Run in development mode (rebuilds on change)
	@echo " Starting development mode..."
	@echo "Press Ctrl+C to stop"
	@while true; do \
		$(MAKE) build && $(BUILD_DIR)/$(BINARY_NAME); \
		sleep 1; \
	done

release: test ## Build release binaries for multiple platforms
	@echo " Building release binaries..."
	@mkdir -p $(BUILD_DIR)/releases
	# Linux amd64
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)
	# Linux arm64
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)_$(VERSION)_linux_arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)
	# macOS amd64
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)
	# macOS arm64
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)_$(VERSION)_darwin_arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)
	# Windows amd64
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe
	cd $(BUILD_DIR) && zip -q releases/$(BINARY_NAME)_$(VERSION)_windows_amd64.zip $(BINARY_NAME).exe
	# Generate checksums
	cd $(BUILD_DIR)/releases && sha256sum * > checksums.txt
	@echo " Release artifacts:"
	@ls -lh $(BUILD_DIR)/releases

setup: ## Set up development environment
	@echo " Setting up development environment..."
	@go mod download
	@echo " Setup complete"

demo: build ## Run a quick demo
	@echo "🎬 Running demo..."
	@echo ""
	@echo "Initializing envoke..."
	@$(BUILD_DIR)/$(BINARY_NAME) init || true
	@echo ""
	@echo "Creating dev environment..."
	@$(BUILD_DIR)/$(BINARY_NAME) env create dev || true
	@echo ""
	@echo "Listing environments..."
	@$(BUILD_DIR)/$(BINARY_NAME) env list || true
	@echo ""
	@echo "✅ Demo complete. Try: make run"

.DEFAULT_GOAL := help
