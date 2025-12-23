.PHONY: build install clean run test

# Binary name
BINARY_NAME=orbit

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Suppress directory messages
MAKEFLAGS += --no-print-directory

# Build directory
BUILD_DIR=bin

# Version info
VERSION=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

## build: Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Install the binary to /usr/local/bin
install: build
	@echo " Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

## run: Run the application
run:
	$(GORUN) .

## clean: Clean build files
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo " Clean complete!"

## test: Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

## deps: Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo " Dependencies downloaded!"

## lint: Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

## help: Show this help
help:
	@echo "🌍 Orbit - Project Manager"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

# Default target
.DEFAULT_GOAL := help
