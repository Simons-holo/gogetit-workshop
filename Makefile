.PHONY: build test lint fmt clean install help

# Binary name
BINARY_NAME=gogetit
MAIN_PATH=./cmd/gogetit

# Build directory
BUILD_DIR=dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run golangci-lint
	@echo "Running linter..."
	$(GOLINT) run ./...

lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linter with auto-fix..."
	$(GOLINT) run --fix ./...

fmt: ## Format code with gofmt
	@echo "Formatting code..."
	$(GOFMT) -s -w .

fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@! $(GOFMT) -s -d . | grep '^'

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	$(GOCLEAN)

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(MAIN_PATH)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

all: deps fmt lint test build ## Run all checks and build

ci: fmt-check lint test ## Run CI checks
