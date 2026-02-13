.PHONY: build test lint clean install plugin help

# Build variables
BINARY_NAME=linterlog
CMD_DIR=./cmd/linterlog
PLUGIN_DIR=./plugin
BUILD_DIR=./bin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the standalone linter binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

plugin: ## Build the golangci-lint plugin
	@echo "Building golangci-lint plugin..."
	@mkdir -p $(BUILD_DIR)
	go build -buildmode=plugin -o $(BUILD_DIR)/$(BINARY_NAME).so $(PLUGIN_DIR)
	@echo "Plugin built at $(BUILD_DIR)/$(BINARY_NAME).so"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

lint: ## Run linter on the project itself
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

install: build ## Install the linter to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(CMD_DIR)
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

example: build ## Run the linter on example code
	@echo "Running linter on testdata..."
	@$(BUILD_DIR)/$(BINARY_NAME) ./testdata/src/a/... || true
	@echo "Example complete (exit code ignored for demo purposes)"

.DEFAULT_GOAL := help
