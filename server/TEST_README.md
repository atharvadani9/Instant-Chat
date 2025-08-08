# Backend Unit Tests

This directory contains comprehensive unit tests for the instant-chat backend application.

## Test Coverage

The test suite covers the following components:

### 1. Crypto Package (`internal/crypto/encryption_test.go`)
- **Encryption/Decryption**: Tests round-trip encryption and decryption
- **Unique Results**: Verifies that encryption produces unique results for the same input
- **Invalid Data Handling**: Tests error handling for corrupted or invalid encrypted data
- **Edge Cases**: Tests empty strings and various content types
- **Benchmarks**: Performance tests for encryption and decryption operations

### 2. Utils Package (`internal/utils/utils_test.go`)
- **JSON Writing**: Tests HTTP JSON response writing with various data types
- **WebSocket Messaging**: Tests WebSocket message writing functionality
- **Error Handling**: Tests error scenarios for invalid data and closed connections
- **Benchmarks**: Performance tests for JSON operations

### 3. Store Package (`internal/store/`)
#### User Store (`user_store_test.go`)
- **Password Hashing**: Tests bcrypt password hashing with various password types
- **Password Verification**: Tests password checking with correct and incorrect passwords
- **Hash Consistency**: Verifies that multiple hashes of the same password are different but verify correctly
- **JSON Marshaling**: Tests that password hashes are excluded from JSON output
- **Edge Cases**: Tests empty passwords, unicode passwords, and long passwords

#### Message Store (`message_store_test.go`)
- **JSON Marshaling**: Tests message serialization and deserialization
- **Field Validation**: Tests message structure and field validation
- **Content Types**: Tests various message content types (unicode, special characters, etc.)
- **Security**: Verifies encrypted content is excluded from JSON output

### 4. API Package (`internal/api/user_handler_test.go`)
- **User Registration**: Tests successful registration and duplicate user handling
- **User Login**: Tests successful login and invalid credential handling
- **Input Validation**: Tests empty username/password validation
- **Error Handling**: Tests database errors and authentication failures
- **HTTP Responses**: Verifies correct HTTP status codes and response formats

### 5. Application Package (`internal/app/app_test.go`)
- **Health Check**: Tests the health check endpoint functionality
- **HTTP Methods**: Tests health check with various HTTP methods
- **Concurrent Access**: Tests concurrent requests to health check
- **Error Scenarios**: Tests behavior with nil logger
- **Response Consistency**: Verifies consistent response format

### 6. Routes Package (`routes/routes_test.go`)
- **Route Setup**: Tests router initialization and configuration
- **CORS Headers**: Tests Cross-Origin Resource Sharing configuration
- **Endpoint Existence**: Verifies all defined routes are accessible
- **Error Handling**: Tests 404 and method not allowed responses
- **WebSocket Routes**: Tests WebSocket endpoint availability
- **Concurrent Requests**: Tests router under concurrent load

## Running Tests

### Prerequisites
The tests require the `ENCRYPTION_KEY` environment variable to be set. This is automatically handled by the test scripts.

### Quick Start
```bash
# Windows PowerShell
.\test.ps1

# Or manually set environment variable and run tests
$env:ENCRYPTION_KEY="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
go test ./... -v
```

### Individual Package Testing
```bash
# Test specific packages
go test ./internal/crypto -v
go test ./internal/utils -v
go test ./internal/store -v
go test ./internal/api -v
go test ./internal/app -v
go test ./routes -v
```

### Coverage Reports
```bash
# Generate coverage report
go test ./... -cover

# Generate detailed coverage HTML report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Benchmarks
```bash
# Run performance benchmarks
go test ./... -bench=.
```

## Test Structure

### Mocking Strategy
- **User Store**: Uses testify/mock for database operations
- **Message Store**: Uses testify/mock for database operations
- **HTTP Testing**: Uses httptest for HTTP request/response testing
- **WebSocket Testing**: Uses gorilla/websocket test utilities

### Test Categories
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **HTTP Tests**: Test API endpoints and middleware
4. **Concurrent Tests**: Test thread safety and concurrent access
5. **Benchmark Tests**: Test performance characteristics

### Error Testing
- Invalid input validation
- Database error simulation
- Network error handling
- Authentication failures
- Encryption/decryption errors

## Dependencies

The test suite uses the following testing libraries:
- `github.com/stretchr/testify` - Assertions and mocking
- `net/http/httptest` - HTTP testing utilities
- `github.com/gorilla/websocket` - WebSocket testing
- Standard Go testing package

## Best Practices

1. **Test Isolation**: Each test is independent and can run in any order
2. **Mock Usage**: External dependencies are mocked to ensure unit test isolation
3. **Error Coverage**: Both success and failure scenarios are tested
4. **Edge Cases**: Boundary conditions and edge cases are thoroughly tested
5. **Performance**: Benchmarks ensure performance regressions are caught
6. **Concurrent Safety**: Tests verify thread-safe operations

## Continuous Integration

These tests are designed to run in CI/CD pipelines. The test script handles environment setup automatically, making it suitable for automated testing environments.
