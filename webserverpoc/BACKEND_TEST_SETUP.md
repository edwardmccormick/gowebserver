# Backend Testing Setup Guide

This guide provides detailed instructions for setting up and running tests for the Go backend of our dating web application.

## Table of Contents

1. [Installation](#installation)
2. [Running Tests](#running-tests)
3. [Test Coverage](#test-coverage)
4. [Testing Components](#testing-components)
5. [Mocking Strategies](#mocking-strategies)
6. [Common Issues](#common-issues)

## Installation

First, make sure all required testing dependencies are installed:

```bash
# Navigate to the backend directory
cd webserverpoc

# Install testing dependencies
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/mock
go get -u github.com/DATA-DOG/go-sqlmock
```

## Running Tests

To run all tests:

```bash
cd webserverpoc
go test -v ./...
```

To run a specific test file:

```bash
go test -v ./path/to/file_test.go
```

To run a specific test function:

```bash
go test -v -run TestFunctionName
```

## Test Coverage

To run tests with coverage:

```bash
go test -v ./... -coverprofile=coverage.out
```

To view the coverage report:

```bash
go tool cover -html=coverage.out
```

The coverage report will open in your default browser, showing which lines are covered by tests.

## Testing Components

### Testing HTTP Handlers

For HTTP handlers, use the `httptest` package and `gin.CreateTestContext`:

```go
func TestYourHandler(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    // Create mock DB
    mockDB := new(MockDBForHandlers)
    originalDB := db
    db = mockDB
    defer func() { db = originalDB }()
    
    // Set up mock expectations
    mockDB.On("Some method", mock.Anything).Return(someResult)
    
    // Set up handler with middleware if needed
    router.GET("/your-endpoint", func(c *gin.Context) {
        c.Set("userID", uint(1)) // Mock JWT middleware if needed
        YourHandler(c)
    })
    
    // Create request
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/your-endpoint", nil)
    
    // Perform the request
    router.ServeHTTP(w, req)
    
    // Check response
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response YourResponseType
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, expectedValue, response.SomeField)
}
```

### Testing WebSockets

For WebSocket testing, use the `testutils/websocket_mock.go` utilities:

```go
func TestWebSocketHandler(t *testing.T) {
    // Setup test environment
    gin.SetMode(gin.TestMode)
    
    // Create a router with the WebSocket handler
    router := gin.New()
    router.GET("/ws", WebsocketHandler)
    
    // Start a test server
    server := httptest.NewServer(router)
    defer server.Close()
    
    // Create WebSocket URL
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=1&user_id=1"
    
    // Connect to the WebSocket
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("Failed to connect to WebSocket: %v", err)
    }
    defer conn.Close()
    
    // Send a test message
    testMessage := ChatMessage{
        Message: "Test message",
        Who: 1,
        MatchID: 1,
    }
    err = conn.WriteJSON(testMessage)
    assert.NoError(t, err)
    
    // For more complex tests, use the WebSocketTestClient from testutils
}
```

### Testing Database Operations

Use the `MockDB` in `testutils/mockdb.go`:

```go
func TestDatabaseOperation(t *testing.T) {
    // Setup mock DB
    mockDB, mockResult := testutils.CreateMockDB()
    originalDB := db
    db = mockDB
    defer func() { db = originalDB }()
    
    // Set up mock expectations
    mockDB.On("First", mock.AnythingOfType("*main.User"), uint(1)).Return(mockResult).Run(func(args mock.Arguments) {
        // Populate the destination object
        arg := args.Get(0).(*User)
        *arg = User{
            ID:    1,
            Email: "test@example.com",
        }
    })
    
    // Call the function that uses the database
    user, err := GetUserByID(1)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    
    // Verify all expectations were met
    mockDB.AssertExpectations(t)
}
```

### Testing Middleware

```go
func TestJwtMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    // Create a router with a protected endpoint
    router := gin.New()
    router.GET("/protected", JwtMiddleware, func(c *gin.Context) {
        c.String(http.StatusOK, "protected")
    })
    
    // Create a valid token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": float64(1),
        "exp": time.Now().Add(time.Hour).Unix(),
    })
    tokenString, _ := token.SignedString(JwtSecret)
    
    // Test with valid token in header
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", tokenString)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Test with no token
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/protected", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```

## Mocking Strategies

### Database Mocking

Use the `MockDB` from `testutils/mockdb.go`:

```go
mockDB := new(testutils.MockDB)
originalDB := db
db = mockDB
defer func() { db = originalDB }()

// Set expectations
mockDB.On("First", mock.AnythingOfType("*main.User"), mock.Anything).Return(&gorm.DB{Error: nil})
```

### External API Mocking

For external APIs, create interface wrappers and mock implementations:

```go
// Define interface
type AIClient interface {
    GenerateContent(prompt string) (string, error)
}

// Mock implementation
type MockAIClient struct {
    mock.Mock
}

func (m *MockAIClient) GenerateContent(prompt string) (string, error) {
    args := m.Called(prompt)
    return args.String(0), args.Error(1)
}

// In tests
mockClient := new(MockAIClient)
originalClient := aiClient
aiClient = mockClient
defer func() { aiClient = originalClient }()

mockClient.On("GenerateContent", "test prompt").Return("mock response", nil)
```

### HTTP Client Mocking

Use `httptest.Server` for mocking HTTP responses:

```go
// Create a test server that returns a specific response
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"result":"success"}`))
}))
defer server.Close()

// Use server.URL as the base URL for your HTTP client
```

## Common Issues

### Race Conditions in Tests

When testing concurrent code (like WebSockets), use the `-race` flag to detect race conditions:

```bash
go test -race ./...
```

### Dependency Injection

For easier testing, consider refactoring code to use dependency injection:

```go
// Before
func Handler(c *gin.Context) {
    // Uses global db variable
}

// After
func Handler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Uses injected db
    }
}
```

### Mock Verification Failures

If you're seeing errors like "mock: Unexpected Method Call", ensure that:

1. All expected methods are defined before calling the function under test
2. The method signatures match exactly
3. The number and types of arguments match
4. You're checking the correct object

Example fix:

```go
// Incorrect
mockDB.On("First", mock.Anything, 1).Return(mockResult)

// Correct
mockDB.On("First", mock.AnythingOfType("*main.User"), uint(1)).Return(mockResult)
```