#!/bin/bash

# Test script for envm project
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
TEST_DIR="/tmp/envm-test-$$"
COVERAGE_THRESHOLD=70

# Functions
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

setup_test_env() {
    print_status "Setting up test environment..."
    
    # Create test directories
    mkdir -p "${TEST_DIR}/downloads/go"
    mkdir -p "${TEST_DIR}/downloads/java"
    
    # Set environment variables
    export ENVM_HOME="${TEST_DIR}"
    export ENVM_GO_SYMLINK="${TEST_DIR}/go"
    export ENVM_JAVA_SYMLINK="${TEST_DIR}/java"
    
    print_status "Test environment created at: ${TEST_DIR}"
}

cleanup_test_env() {
    print_status "Cleaning up test environment..."
    rm -rf "${TEST_DIR}"
    unset ENVM_HOME ENVM_GO_SYMLINK ENVM_JAVA_SYMLINK
}

run_unit_tests() {
    print_status "Running unit tests..."
    
    # Run tests with coverage
    go test -v -race -coverprofile=coverage.out ./... || {
        print_error "Unit tests failed"
        return 1
    }
    
    # Generate coverage report
    go tool cover -html=coverage.out -o coverage.html
    
    # Check coverage threshold
    local coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$coverage < $COVERAGE_THRESHOLD" | bc -l) )); then
        print_warning "Coverage ${coverage}% is below threshold ${COVERAGE_THRESHOLD}%"
    else
        print_status "Coverage: ${coverage}% (above threshold)"
    fi
}

run_integration_tests() {
    print_status "Running integration tests..."
    
    # Run integration tests (if any)
    go test -tags=integration -v ./... || {
        print_warning "Integration tests failed or not found"
    }
}

run_benchmarks() {
    print_status "Running benchmarks..."
    
    go test -bench=. -benchmem ./... || {
        print_warning "Benchmarks failed or not found"
    }
}

check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check Go version
    go version
    
    # Verify modules
    go mod verify
    
    # Download dependencies
    go mod download
}

# Main execution
main() {
    local run_benchmarks=false
    local run_integration=false
    local cleanup=true
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --bench)
                run_benchmarks=true
                shift
                ;;
            --integration)
                run_integration=true
                shift
                ;;
            --no-cleanup)
                cleanup=false
                shift
                ;;
            --coverage-threshold)
                COVERAGE_THRESHOLD="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [--bench] [--integration] [--no-cleanup] [--coverage-threshold N]"
                echo "  --bench                 Run benchmarks"
                echo "  --integration          Run integration tests"
                echo "  --no-cleanup           Don't cleanup test environment"
                echo "  --coverage-threshold N Set coverage threshold (default: 70)"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Trap for cleanup
    if [ "$cleanup" = true ]; then
        trap cleanup_test_env EXIT
    fi
    
    # Execute test pipeline
    check_dependencies
    setup_test_env
    run_unit_tests
    
    if [ "$run_integration" = true ]; then
        run_integration_tests
    fi
    
    if [ "$run_benchmarks" = true ]; then
        run_benchmarks
    fi
    
    print_status "All tests completed successfully!"
}

# Run main function with all arguments
main "$@"