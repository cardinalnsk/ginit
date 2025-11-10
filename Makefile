.PHONY: build clean test run install uninstall help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ginit
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
all: build

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/ginit

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf bin/

# Run tests
test:
	$(GOTEST) -v ./...

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/ginit
	./$(BINARY_NAME)

# Install the application
install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/

# Uninstall the application
uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Cross-platform builds
build-all: build-linux build-windows build-darwin

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/ginit
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/ginit

# Build for Windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/ginit
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-windows-arm64.exe ./cmd/ginit

# Build for macOS
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/ginit
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/ginit

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  run           - Build and run the application"
	@echo "  install       - Install to /usr/local/bin"
	@echo "  uninstall     - Remove from /usr/local/bin"
	@echo "  build-all     - Build for all platforms (Linux, Windows, macOS)"
	@echo "  build-linux   - Build for Linux (amd64, arm64)"
	@echo "  build-windows - Build for Windows (amd64, arm64)"
	@echo "  build-darwin  - Build for macOS (amd64, arm64)"
	@echo "  help          - Show this help message"