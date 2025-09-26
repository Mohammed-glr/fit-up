# Test Commands for Lornian Backend

.PHONY: test test-unit test-integration test-e2e test-auth test-user test-ai test-gateway
.PHONY: test-setup test-teardown test-db-setup test-db-teardown
.PHONY: test-coverage test-coverage-html test-lint test-fmt

# =============================================================================
# MAIN TEST COMMANDS
# =============================================================================

# Run all tests
test: test-unit test-integration

# Run only unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -race -timeout=30s ./tests/*/unit/... ./services/*/tests/unit/...

# Run only integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -race -timeout=60s ./tests/*/integration/... ./services/*/tests/integration/...

# Run end-to-end tests
test-e2e:
	@echo "Running e2e tests..."
	go test -v -race -timeout=120s ./tests/*/e2e/... ./services/*/tests/e2e/...

# =============================================================================
# SERVICE-SPECIFIC TESTS
# =============================================================================

# Test auth service only
test-auth:
	@echo "Running auth service tests..."
	go test -v -race ./services/auth-service/tests/...

# Test user service only  
test-user:
	@echo "Running user service tests..."
	go test -v -race ./services/user-service/tests/...

# Test AI service only
test-ai:
	@echo "Running AI service tests..."
	go test -v -race ./services/ai-service/tests/...

# Test API gateway only
test-gateway:
	@echo "Running API gateway tests..."
	go test -v -race ./services/api-gateway/tests/...

# =============================================================================
# TEST ENVIRONMENT SETUP
# =============================================================================

# Set up test environment (databases, etc.)
test-setup:
	@echo "Setting up test environment..."
	cp .env.test.example .env.test
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for services to be ready..."
	sleep 30
	@echo "Running database migrations..."
	@make test-db-setup

# Tear down test environment
test-teardown:
	@echo "Tearing down test environment..."
	docker-compose -f docker-compose.test.yml down -v
	rm -f .env.test

# Set up test databases
test-db-setup:
	@echo "Setting up test databases..."
	# Add database migration commands here when ready
	@echo "Test databases ready!"

# Clean up test databases
test-db-teardown:
	@echo "Cleaning up test databases..."
	# Add database cleanup commands here when ready

# =============================================================================
# TEST COVERAGE AND QUALITY
# =============================================================================

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./tests/... ./services/*/tests/...
	go tool cover -func=coverage.out

# Generate HTML coverage report
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
test-lint:
	@echo "Running linting..."
	golangci-lint run ./tests/... ./services/*/tests/...

# Format test code
test-fmt:
	@echo "Formatting test code..."
	go fmt ./tests/... ./services/*/tests/...

# =============================================================================
# TEST UTILITIES
# =============================================================================

# Run tests in watch mode (requires entr)
test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" -not -path "./vendor/*" | entr -r make test-unit

# Run specific test file
test-file:
	@test -n "$(FILE)" || (echo "Usage: make test-file FILE=path/to/test_file.go" && exit 1)
	go test -v -race $(FILE)

# Run specific test function
test-func:
	@test -n "$(FUNC)" || (echo "Usage: make test-func FUNC=TestFunctionName" && exit 1)
	go test -v -race -run $(FUNC) ./...

# =============================================================================
# BENCHMARKS
# =============================================================================

# Run benchmarks
test-bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./tests/... ./services/*/tests/...

# =============================================================================
# HELP
# =============================================================================

# Show available test commands
test-help:
	@echo "Available test commands:"
	@echo "  test              - Run all tests (unit + integration)"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  test-e2e          - Run end-to-end tests"
	@echo ""
	@echo "Service-specific tests:"
	@echo "  test-auth         - Test auth service only"
	@echo "  test-user         - Test user service only"
	@echo "  test-ai           - Test AI service only"
	@echo "  test-gateway      - Test API gateway only"
	@echo ""
	@echo "Environment management:"
	@echo "  test-setup        - Set up test environment"
	@echo "  test-teardown     - Tear down test environment"
	@echo "  test-db-setup     - Set up test databases"
	@echo "  test-db-teardown  - Clean up test databases"
	@echo ""
	@echo "Coverage and quality:"
	@echo "  test-coverage     - Run tests with coverage"
	@echo "  test-coverage-html- Generate HTML coverage report"
	@echo "  test-lint         - Run linting on test code"
	@echo "  test-fmt          - Format test code"
	@echo ""
	@echo "Utilities:"
	@echo "  test-watch        - Run tests in watch mode"
	@echo "  test-file FILE=   - Run specific test file"
	@echo "  test-func FUNC=   - Run specific test function"
	@echo "  test-bench        - Run benchmarks"
