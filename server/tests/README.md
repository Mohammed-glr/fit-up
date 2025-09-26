# Testing Structure

This directory contains all tests for the Lornian Backend services.

## Directory Structure

```
tests/
├── shared/                     # Shared test utilities and mocks
│   ├── testutils/             # Common test utilities
│   └── mocks/                 # Shared mock implementations
├── auth-service/              # Authentication service tests (legacy)
├── user-service/              # User service tests (legacy)
├── ai-service/                # AI service tests (legacy)
├── api-gateway/               # API Gateway tests (legacy)
├── Makefile                   # Test automation commands
└── README.md                  # This file

services/
├── auth-service/tests/        # Auth service tests (recommended)
│   ├── unit/                  # Unit tests
│   ├── integration/           # Integration tests
│   ├── mocks/                 # Service-specific mocks
│   └── fixtures/              # Test data and fixtures
├── user-service/tests/        # User service tests (when created)
├── ai-service/tests/          # AI service tests (when created)
└── api-gateway/tests/         # API Gateway tests (when created)
```

## Test Organization

### Recommended Approach
Tests are now organized within each service directory under a `tests/` subdirectory. This approach:
- Keeps tests close to the code they test
- Allows access to internal packages
- Makes service-specific testing easier
- Follows Go best practices

### Legacy Structure
The top-level `tests/` directory is maintained for:
- Shared test utilities
- Cross-service integration tests
- E2E tests that span multiple services

## Test Types

### Unit Tests
- Test individual functions and methods in isolation
- Use mocks for external dependencies
- Fast execution, no external dependencies
- Located in `services/*/tests/unit/`

### Integration Tests
- Test interactions between components
- May use real databases (test instances)
- Test service integrations
- Located in `services/*/tests/integration/`

### End-to-End Tests
- Test complete user workflows
- Test through HTTP APIs
- Use real or containerized dependencies
- Located in `tests/*/e2e/` for cross-service tests

## Running Tests

### Quick Start
```bash
# Set up test environment
make test-setup

# Run all tests
make test

# Run service-specific tests
make test-auth
make test-user

# Tear down test environment
make test-teardown
```

### Detailed Commands

#### All Tests
```bash
go test ./tests/... ./services/*/tests/...
```

#### Service-Specific Tests
```bash
# Auth service tests
go test ./services/auth-service/tests/...

# User service tests (when available)
go test ./services/user-service/tests/...
```

#### Test Type Specific
```bash
# Unit tests only
go test ./services/auth-service/tests/unit/...

# Integration tests only
go test ./services/auth-service/tests/integration/...
```

#### With Coverage
```bash
# Generate coverage report
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

## Test Environment

### Docker Test Environment
Use the provided docker-compose file for test dependencies:

```bash
# Start test services
docker-compose -f docker-compose.test.yml up -d

# Stop test services
docker-compose -f docker-compose.test.yml down -v
```

### Environment Variables
Copy `.env.test.example` to `.env.test` and adjust values:

```bash
cp .env.test.example .env.test
```

Key test environment variables:
- `DATABASE_URL`: Test PostgreSQL connection
- `MONGODB_URI`: Test MongoDB connection (for user service)
- `JWT_SECRET`: Test JWT secret
- `SMTP_HOST`: Test email service (MailHog)

## Test Conventions

1. **File Naming**: Test files should end with `_test.go`
2. **Function Naming**: Test functions should start with `Test`
3. **Package Naming**: Test packages should match the package being tested or use `_test` suffix
4. **Mocks**: Place service-specific mocks in `services/*/tests/mocks/`
5. **Fixtures**: Place test data in `services/*/tests/fixtures/`
6. **Setup/Teardown**: Use `TestMain` for setup and teardown when needed

## Available Make Commands

Run `make test-help` to see all available test commands:

- `test` - Run all tests
- `test-unit` - Run unit tests only
- `test-integration` - Run integration tests only
- `test-e2e` - Run end-to-end tests
- `test-auth` - Test auth service only
- `test-setup` - Set up test environment
- `test-teardown` - Tear down test environment
- `test-coverage` - Run tests with coverage
- `test-lint` - Run linting on tests
- And many more...

## Mock Guidelines

### Service-Specific Mocks
- Located in `services/*/tests/mocks/`
- Implement the actual interfaces used by the service
- Support error injection for testing error scenarios
- Include helper methods for easy test setup

### Shared Mocks
- Located in `tests/shared/mocks/`
- For interfaces used across multiple services
- Database mocks, HTTP client mocks, etc.

## Test Data and Fixtures

### Fixtures
- Located in `services/*/tests/fixtures/`
- Provide factory functions for creating test data
- Should return valid, realistic test objects
- Support variations (admin users, invalid data, etc.)

### Example Usage
```go
// Create test user
user := fixtures.TestUser()

// Create admin user
admin := fixtures.TestAdminUser()

// Create test request
loginReq := fixtures.TestLoginRequest()
```

## CI/CD Integration

The test structure is designed to work with CI/CD pipelines:

1. **Test Environment Setup**: Automated via `make test-setup`
2. **Parallel Execution**: Tests can run in parallel per service
3. **Coverage Reports**: Generated automatically
4. **Test Results**: Standard Go test output for CI integration

## Future Enhancements

- [ ] Add benchmark tests
- [ ] Implement property-based testing
- [ ] Add mutation testing
- [ ] Performance regression tests
- [ ] Visual regression tests for UI components
- [ ] Load testing for API endpoints
