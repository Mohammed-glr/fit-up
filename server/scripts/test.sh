#!/bin/bash

# Test Runner Script for Lornian Backend
# This script provides a convenient way to run tests with proper setup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi
    
    print_success "All dependencies are available"
}

# Setup test environment
setup_test_env() {
    print_status "Setting up test environment..."
    
    # Copy test environment file if it doesn't exist
    if [ ! -f .env.test ]; then
        print_status "Creating .env.test from example..."
        cp .env.test.example .env.test
    fi
    
    # Start test services
    print_status "Starting test services with Docker Compose..."
    docker-compose -f docker-compose.test.yml up -d
    
    # Wait for services to be ready
    print_status "Waiting for services to be ready..."
    sleep 30
    
    # Check if services are running
    if ! docker-compose -f docker-compose.test.yml ps | grep -q "Up"; then
        print_error "Test services failed to start"
        exit 1
    fi
    
    print_success "Test environment is ready"
}

# Teardown test environment
teardown_test_env() {
    print_status "Tearing down test environment..."
    docker-compose -f docker-compose.test.yml down -v
    print_success "Test environment cleaned up"
}

# Run tests with coverage
run_tests_with_coverage() {
    local test_pattern="$1"
    local output_file="coverage.out"
    
    print_status "Running tests with coverage..."
    
    if [ -n "$test_pattern" ]; then
        go test -v -race -coverprofile="$output_file" "$test_pattern"
    else
        go test -v -race -coverprofile="$output_file" ./tests/... ./services/*/tests/...
    fi
    
    if [ -f "$output_file" ]; then
        print_status "Generating coverage report..."
        go tool cover -func="$output_file"
        
        # Generate HTML report
        go tool cover -html="$output_file" -o coverage.html
        print_success "Coverage report generated: coverage.html"
    fi
}

# Run specific test type
run_test_type() {
    local test_type="$1"
    local service="$2"
    
    case "$test_type" in
        "unit")
            if [ -n "$service" ]; then
                print_status "Running unit tests for $service service..."
                go test -v -race ./services/$service/tests/unit/...
            else
                print_status "Running all unit tests..."
                go test -v -race ./services/*/tests/unit/...
            fi
            ;;
        "integration")
            if [ -n "$service" ]; then
                print_status "Running integration tests for $service service..."
                go test -v -race ./services/$service/tests/integration/...
            else
                print_status "Running all integration tests..."
                go test -v -race ./services/*/tests/integration/...
            fi
            ;;
        "e2e")
            print_status "Running end-to-end tests..."
            go test -v -race ./tests/*/e2e/...
            ;;
        *)
            print_error "Unknown test type: $test_type"
            print_status "Available types: unit, integration, e2e"
            exit 1
            ;;
    esac
}

# Show usage information
show_usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  setup                    Set up test environment"
    echo "  teardown                 Tear down test environment"
    echo "  test [pattern]           Run tests (optionally with pattern)"
    echo "  test-type TYPE [SERVICE] Run specific test type (unit/integration/e2e)"
    echo "  coverage [pattern]       Run tests with coverage"
    echo "  lint                     Run linting"
    echo "  fmt                      Format code"
    echo "  clean                    Clean test artifacts"
    echo ""
    echo "Examples:"
    echo "  $0 setup                          # Set up test environment"
    echo "  $0 test                           # Run all tests"
    echo "  $0 test ./services/auth-service/  # Run auth service tests"
    echo "  $0 test-type unit auth-service    # Run auth service unit tests"
    echo "  $0 coverage                       # Run tests with coverage"
    echo "  $0 teardown                       # Clean up test environment"
}

# Main script logic
main() {
    case "${1:-}" in
        "setup")
            check_dependencies
            setup_test_env
            ;;
        "teardown")
            teardown_test_env
            ;;
        "test")
            check_dependencies
            if [ -n "${2:-}" ]; then
                print_status "Running tests with pattern: $2"
                go test -v -race "$2"
            else
                print_status "Running all tests..."
                go test -v -race ./tests/... ./services/*/tests/...
            fi
            ;;
        "test-type")
            check_dependencies
            if [ -z "${2:-}" ]; then
                print_error "Test type is required"
                show_usage
                exit 1
            fi
            run_test_type "$2" "${3:-}"
            ;;
        "coverage")
            check_dependencies
            run_tests_with_coverage "${2:-}"
            ;;
        "lint")
            print_status "Running linting..."
            if command -v golangci-lint &> /dev/null; then
                golangci-lint run ./tests/... ./services/*/tests/...
            else
                print_warning "golangci-lint not found, using go vet..."
                go vet ./tests/... ./services/*/tests/...
            fi
            ;;
        "fmt")
            print_status "Formatting code..."
            go fmt ./tests/... ./services/*/tests/...
            print_success "Code formatted"
            ;;
        "clean")
            print_status "Cleaning test artifacts..."
            rm -f coverage.out coverage.html
            docker system prune -f
            print_success "Test artifacts cleaned"
            ;;
        "help"|"-h"|"--help"|"")
            show_usage
            ;;
        *)
            print_error "Unknown command: $1"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
