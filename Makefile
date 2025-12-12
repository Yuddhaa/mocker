# Default command
all: help

# Run the app in development mode with hot-reload
dev: # install-air .air.toml
	@echo "Starting server with hot-reload (Air)..."
	@air

# Build the Go binary locally
build:
	@echo "Building Go binary..."
	# @mkdir -p tmp
	@go build -o ./tmp/main ./...

run: build
	./tmp/main

# Install the 'air' hot-reload tool
install-air:
	@echo "Installing/Updating 'air'..."
	@go install github.com/air-verse/air@latest

build-macos-amd:
	@echo "building for macos amd"
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/mocker-macos-amd main.go 
	@echo "macos amd done\n"

build-macos:
	@echo "building for macos"
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/mocker-macos main.go 
	@echo "macos done\n"

build-linux:
	@echo "building for linux"
	@GOOS=linux GOARCH=amd64 go build -o ./bin/mocker-linux main.go 
	@echo "linux done\n"

build-windows:
	@echo "building for windows"
	@GOOS=windows GOARCH=amd64 go build -o ./bin/mocker-windows.exe main.go
	@echo "windows done\n"

build-all: build-macos build-macos-amd build-linux build-windows
# Show help menu
help:
	@echo "Available commands:"
	@echo "  make dev          - Start local dev server with hot-reload"
	@echo "  make build        - Build the Go binary locally"
	@echo "  make run          - Build and Run the Go binary locally"
	@echo "  make install-air  - Install the 'air' tool"

