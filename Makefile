# Makefile for Knetz

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet

# Binary name
BINARY_NAME = knetz
BINARY_PATH = bin/$(BINARY_NAME)

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Platforms for cross-compilation
PLATFORMS = darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build clean test coverage fmt vet lint install help

all: test build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/knetz

## run: Run the application
run: build
	./$(BINARY_PATH)

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf dist/

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

## tidy: Tidy go modules
tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy

## install: Install the binary to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	$(GOCMD) install $(LDFLAGS) ./cmd/knetz

## build-all: Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		output_name=dist/$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = "windows" ]; then output_name=$${output_name}.exe; fi; \
		echo "Building $$output_name..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) $(LDFLAGS) -o $$output_name ./cmd/knetz; \
	done

## release: Create release archives
release: build-all
	@echo "Creating release archives..."
	@cd dist && for binary in knetz-*; do \
		if [ -f "$$binary" ]; then \
			tar -czf "$$binary.tar.gz" "$$binary"; \
			echo "Created $$binary.tar.gz"; \
		fi \
	done

## help: Show this help message
help:
	@echo "Knetz - Cross-Cluster Dependency & Version Manager"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^##/ {desc = substr($$0, 4); getline; printf "  %-15s %s\n", $$2, desc}' $(MAKEFILE_LIST)

