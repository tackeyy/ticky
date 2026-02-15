# Testing Guide

## Overview

ticky uses Go's standard `testing` package along with `net/http/httptest` for HTTP mocking. All code contributions must include comprehensive tests.

## Test Philosophy

- **Write tests first** when fixing bugs (TDD approach)
- **Test coverage** for all new features
- **No breaking changes** without tests proving backward compatibility
- **Fast execution** - unit tests should run in milliseconds

## Test Structure

### Directory Layout

```
internal/ticktick/
  priority_test.go   # Priority parser tests (37 subtests)
  date_test.go       # Date parser tests (31 subtests)
  token_test.go      # Token I/O and config tests (10 tests)
  client_test.go     # HTTP API client tests (16 tests)
  auth_test.go       # OAuth authentication tests (10 tests)
  export_test.go     # Test helpers (NewTestClient, SetBaseURL, SetTokenURL)
```

### Naming Conventions

- Test files: `*_test.go` (Go convention)
- Test functions: `TestFunctionName_Scenario` (e.g., `TestParsePriority_ValidInputs`)
- Subtests: Descriptive lowercase names (e.g., `"low keyword"`, `"empty string"`)

### Package Conventions

Tests that need access to unexported functions use the same package name:

```go
package ticktick  // in token_test.go — accesses unexported configDir()
```

Tests that only use exported APIs use the `_test` suffix:

```go
package ticktick_test  // in client_test.go — uses only exported Client methods
```

## Test Categories

### 1. Unit Tests

Test individual functions in isolation using table-driven tests.

**Example** (from `priority_test.go`):
```go
func TestParsePriority_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"none keyword", "none", PriorityNone},
		{"low keyword", "low", PriorityLow},
		{"medium keyword", "medium", PriorityMedium},
		{"high keyword", "high", PriorityHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if err != nil {
				t.Fatalf("ParsePriority(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParsePriority(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
```

### 2. HTTP Mock Tests (Integration)

Test API client methods using `httptest.NewServer` to mock the TickTick API.

**Example** (from `client_test.go`):
```go
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*ticktick.Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	restore := ticktick.SetBaseURL(server.URL)
	client := ticktick.NewTestClient(server.Client(), "test-token")
	return client, func() {
		restore()
		server.Close()
	}
}

func TestGetProjects_Success(t *testing.T) {
	// Arrange
	want := []ticktick.Project{
		{ID: "proj-1", Name: "Work"},
		{ID: "proj-2", Name: "Personal"},
	}
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	})
	defer cleanup()

	// Act
	got, err := client.GetProjects()

	// Assert
	if err != nil {
		t.Fatalf("GetProjects() returned unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("GetProjects() returned %d projects, want 2", len(got))
	}
}
```

### 3. Authentication Tests

Test OAuth token exchange, refresh, and state generation using `httptest.NewServer`.

**Example** (from `auth_test.go`):
```go
func TestExchangeToken_Success(t *testing.T) {
	// Arrange
	wantToken := ticktick.OAuthToken{
		AccessToken:  "access-token-123",
		TokenType:    "bearer",
		ExpiresIn:    3600,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Basic Auth credentials
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Basic Auth not set")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wantToken)
	}))
	defer server.Close()

	restore := ticktick.SetTokenURL(server.URL)
	defer restore()

	// Act
	got, err := ticktick.ExchangeToken("test-client-id", "test-client-secret", "test-code")

	// Assert
	if err != nil {
		t.Fatalf("ExchangeToken() returned unexpected error: %v", err)
	}
	if got.AccessToken != wantToken.AccessToken {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, wantToken.AccessToken)
	}
}
```

### 4. File I/O Tests

Test token storage and config file operations using isolated temp directories.

**Example** (from `token_test.go`):
```go
func TestSaveAndLoadToken_RoundTrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	// Arrange
	want := &OAuthToken{
		AccessToken:  "test-access-token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
	}

	// Act
	if err := SaveToken(want); err != nil {
		t.Fatalf("SaveToken() returned unexpected error: %v", err)
	}
	got, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() returned unexpected error: %v", err)
	}

	// Assert
	if got.AccessToken != want.AccessToken {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, want.AccessToken)
	}
}
```

## Writing Good Tests

### Follow the AAA Pattern

```go
func TestSomething(t *testing.T) {
	// Arrange: Set up test data
	input := "test"

	// Act: Execute the function
	result, err := MyFunction(input)

	// Assert: Verify the result
	if err != nil {
		t.Fatalf("MyFunction(%q) returned unexpected error: %v", input, err)
	}
	if result != "expected" {
		t.Errorf("MyFunction(%q) = %q, want %q", input, result, "expected")
	}
}
```

### Use Table-Driven Tests

Table-driven tests with `t.Run` subtests are the standard pattern in this project:

```go
func TestMyFunction(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid input", "hello", "HELLO"},
		{"empty input", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MyFunction(tt.input)
			if got != tt.want {
				t.Errorf("MyFunction(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
```

### Use Descriptive Test Names

Bad:
```go
func TestWorks(t *testing.T) {}
func Test1(t *testing.T) {}
```

Good:
```go
func TestParsePriority_ValidInputs(t *testing.T) {}
func TestExchangeToken_HTTPError(t *testing.T) {}
func TestLoadToken_FileNotFound(t *testing.T) {}
```

### Mock External Dependencies

- **HTTP APIs**: Use `net/http/httptest.NewServer` with custom handlers
- **Environment variables**: Use `t.Setenv()` for automatic cleanup
- **File system**: Use `t.TempDir()` for isolated temporary directories
- **Package-level variables**: Use `SetBaseURL()` / `SetTokenURL()` helpers that return restore functions

### Test Edge Cases

Always test:
- Valid input (happy path)
- Invalid input (error cases)
- Boundary values (empty strings, zero values)
- HTTP error responses (401, 500)
- Invalid JSON responses
- Missing files / missing environment variables

## Running Tests

### Basic Commands

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file (by package)
go test -v ./internal/ticktick/

# Run tests matching a pattern
go test -v -run TestParsePriority ./internal/ticktick/

# Run a specific subtest
go test -v -run TestParsePriority_ValidInputs/low_keyword ./internal/ticktick/
```

### Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# View per-function coverage
go tool cover -func=coverage.out
```

### Debugging Tests

```bash
# Run with verbose output
go test -v ./...

# Run single test with debug output
go test -v -run TestGetProjects_Success ./internal/ticktick/

# Use delve debugger
dlv test ./internal/ticktick/ -- -test.run TestGetProjects_Success
```

## Test Coverage Requirements

| Category | Requirement |
|----------|-------------|
| **New Features** | Tests required for all new code |
| **Bug Fixes** | Regression test required |
| **Refactoring** | Maintain existing coverage |
| **Overall Project** | Target: 80%+ (testable code) |

## Common Testing Patterns

### Testing HTTP API Methods

```go
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*ticktick.Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	restore := ticktick.SetBaseURL(server.URL)
	client := ticktick.NewTestClient(server.Client(), "test-token")
	return client, func() {
		restore()
		server.Close()
	}
}
```

### Testing Environment Variables

```go
func TestMyFunction(t *testing.T) {
	t.Setenv("TICKTICK_CLIENT_ID", "test-client-id")
	t.Setenv("TICKTICK_CLIENT_SECRET", "test-client-secret")
	// env vars are automatically restored after the test
}
```

### Testing File System Operations

```go
func TestTokenStorage(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	// All file operations now use an isolated temporary directory
	// Cleaned up automatically after the test
}
```

### Testing Error Responses

```go
func TestGetProjects_HTTPError(t *testing.T) {
	client, cleanup := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	})
	defer cleanup()

	_, err := client.GetProjects()
	if err == nil {
		t.Fatal("GetProjects() expected error for 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("error = %q, want to contain 'status 500'", err.Error())
	}
}
```

## Best Practices

### DO

- Write tests before or alongside code
- Use `t.Helper()` in test helper functions
- Use `t.TempDir()` for file system isolation
- Use `t.Setenv()` for environment variable isolation
- Test both success and failure paths
- Use meaningful test and subtest names
- Keep tests simple and focused
- Use `t.Fatalf` for setup failures, `t.Errorf` for assertion failures

### DON'T

- Skip writing tests ("I'll add them later")
- Test implementation details (test behavior, not internals)
- Write tests that depend on other tests
- Use real external APIs in tests
- Leave commented-out test code
- Write flaky tests (tests that sometimes fail)
- Ignore returned errors in test assertions

## Troubleshooting

### "Tests pass locally but fail in CI"

- Check for timezone differences (date tests use `time.Local`)
- Ensure all mock servers are properly closed with `defer cleanup()`
- Verify environment variables are isolated with `t.Setenv()`

### "Cannot access unexported function"

- If your test needs unexported functions, use the same package name (e.g., `package ticktick`)
- For tests of exported APIs only, use the `_test` suffix (e.g., `package ticktick_test`)
- Add test helpers to `export_test.go` if needed

### "Test isolation issues"

- Use `t.TempDir()` instead of hard-coded paths
- Use `t.Setenv()` instead of direct `os.Setenv()` (restores automatically)
- Use `SetBaseURL()` / `SetTokenURL()` with `defer restore()`

## Resources

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Go httptest Documentation](https://pkg.go.dev/net/http/httptest)
- [Table-Driven Tests in Go](https://go.dev/wiki/TableDrivenTests)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

## Questions?

If you have questions about testing:
1. Check existing test files for examples
2. Open an [issue](https://github.com/tackeyy/ticky/issues) with the `question` label
