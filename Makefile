# Makefile for envm project

.PHONY: test test-verbose test-cover test-race clean build help test-windows test-windows-direct test-wsl test-unix-direct

# Default target
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-verbose      - Run tests with verbose output"
	@echo "  test-cover        - Run tests with coverage report"
	@echo "  test-race         - Run tests with race detector"
	@echo "  test-short        - Run tests excluding long-running tests"
	@echo "  test-unit         - Run only unit tests (exclude network tests)"
	@echo "  test-windows      - Run tests using PowerShell script (Windows)"
	@echo "  test-windows-direct - Run tests directly with Windows env vars"
	@echo "  test-wsl          - Run tests using bash script (WSL)"
	@echo "  test-unix-direct  - Run tests directly with Unix env vars"
	@echo "  clean             - Clean test artifacts"
	@echo "  build             - Build the project"
	@echo "  lint              - Run linters"

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	go test -race ./...

# Run short tests only (skip long-running tests)
test-short:
	go test -short ./...

# Run unit tests only (exclude network-dependent tests)
test-unit:
	go test -tags=unit ./...

# Test specific packages
test-cmd:
	go test -v ./cmd/...

test-config:
	go test -v ./internal/config/...

test-arch:
	go test -v ./internal/arch/...

test-commands:
	go test -v ./internal/commands/...

test-util:
	go test -v ./util/...

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html
	go clean -testcache

# Build the project
build:
	go build -o bin/envm .

# Run linters (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests in CI environment
ci-test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Check for test failures in parallel
test-parallel:
	go test -parallel 4 ./...

# Benchmark tests
benchmark:
	go test -bench=. ./...

# Windows specific test target
test-windows:
	@echo "Running tests with PowerShell script..."
	powershell.exe -ExecutionPolicy Bypass -File scripts/test.ps1

# Windows direct test (without script)
test-windows-direct:
	@echo "Running tests with proper Windows environment variables..."
	@powershell.exe -Command "$$env:ENVM_HOME='$$env:USERPROFILE\\.envm-test'; $$env:ENVM_GO_SYMLINK='$$env:USERPROFILE\\.envm-test\\go'; $$env:ENVM_JAVA_SYMLINK='$$env:USERPROFILE\\.envm-test\\java'; New-Item -ItemType Directory -Path '$$env:ENVM_HOME\\downloads\\go' -Force | Out-Null; New-Item -ItemType Directory -Path '$$env:ENVM_HOME\\downloads\\java' -Force | Out-Null; go test -v -coverprofile=coverage.out ./..."

# WSL specific test target  
test-wsl:
	@echo "Running tests with bash script..."
	bash scripts/test.sh

# Unix direct test (without script)
test-unix-direct:
	@echo "Running tests with proper Unix environment variables..."
	@ENVM_HOME="$$HOME/.envm-test" ENVM_GO_SYMLINK="$$HOME/.envm-test/go" ENVM_JAVA_SYMLINK="$$HOME/.envm-test/java" \
	bash -c 'mkdir -p "$$ENVM_HOME/downloads/go" "$$ENVM_HOME/downloads/java" && go test -v -coverprofile=coverage.out ./...'

# Test with different Go versions (requires multiple Go installations)
test-go-versions:
	@echo "Testing with different Go versions..."
	go1.19 version && go1.19 test ./... || echo "Go 1.19 not available"
	go1.20 version && go1.20 test ./... || echo "Go 1.20 not available"
	go1.21 version && go1.21 test ./... || echo "Go 1.21 not available"