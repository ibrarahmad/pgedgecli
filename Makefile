# Project settings
APP_NAME := pgedgecli
BUILD_DIR := build
SRC_DIR := .
VERSION := 1.0.0

# Go settings
GO := go
GOFILES := $(shell find $(SRC_DIR) -name '*.go' -type f)

# Default target
.PHONY: all
all: build

# Build the CLI binary
.PHONY: build
build: $(BUILD_DIR)/$(APP_NAME)
$(BUILD_DIR)/$(APP_NAME): $(GOFILES)
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) -ldflags="-X 'main.Version=$(VERSION)'" $(SRC_DIR)/main.go
	@echo "Built $(APP_NAME) version $(VERSION) successfully!"

# Run the CLI
.PHONY: run
run: build
	./$(BUILD_DIR)/$(APP_NAME)

# Test the project
.PHONY: test
test:
	$(GO) test ./... -v

# Lint the project
.PHONY: lint
lint:
	@which golangci-lint >/dev/null || (echo "Please install golangci-lint first!"; exit 1)
	golangci-lint run

# Clean the build directory
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build directory."

# Display help
.PHONY: help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo "Usage:"
	@echo "  make build        Build the CLI binary"
	@echo "  make run          Run the CLI"
	@echo "  make test         Run tests"
	@echo "  make lint         Run linter"
	@echo "  make clean        Clean build artifacts"
	@echo "  make help         Show this help message"

