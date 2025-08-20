# Testing Strategy for Go Dating Web Application

This document outlines a comprehensive testing strategy for both the Go backend and React frontend components of the dating web application.

## Table of Contents
- [Backend Testing](#backend-testing)
  - [Unit Testing](#backend-unit-testing)
  - [Integration Testing](#backend-integration-testing)
  - [API Testing](#api-testing)
- [Frontend Testing](#frontend-testing)
  - [Component Testing](#component-testing)
  - [Integration Testing](#frontend-integration-testing)
  - [End-to-End Testing](#end-to-end-testing)
- [Special Testing Areas](#special-testing-areas)
  - [WebSocket Testing](#websocket-testing)
  - [SSE Testing](#sse-testing)
  - [AI Integration Testing](#ai-integration-testing)
  - [Database Testing](#database-testing)
- [Test Infrastructure](#test-infrastructure)
- [Implementation Plan](#implementation-plan)
- [Test Implementation Tracking](#test-implementation-tracking)

## Backend Testing

### Backend Unit Testing

#### HTTP Handlers

```go
// handlers_test.go
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockDB implements the database interface for testing
type MockDB struct {
    mock.Mock
}

func (m *MockDB) GetMatches(userID uint) ([]Match, error) {
    args := m.Called(userID)
    return args.Get(0).([]Match), args.Error(1)
}

func TestGetMatches(t *testing.T) {
    // Setup test environment
    gin.SetMode(gin.TestMode)
    
    // Create mock database
    mockDB := new(MockDB)
    
    // Setup test data
    testMatches := []Match{
        {ID: 1, Offered: 1, Accepted: 2, UnreadOffered: 3},
        {ID: 2, Offered: 2, Accepted: 1, UnreadAccepted: 5},
    }
    
    // Set expectations
    mockDB.On("GetMatches", uint(1)).Return(testMatches, nil)
    
    // Create test context
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    // Setup request
    req, _ := http.NewRequest("GET", "/matches/1", nil)
    c.Request = req
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Set("userID", uint(1)) // Simulate middleware setting user ID
    
    // Call handler with dependencies injected
    db = mockDB // Temporarily replace global db with mock
    GetMatchByPersonID(c)
    
    // Assertions
    assert.Equal(t, 200, w.Code)
    
    var response []Match
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, 2, len(response))
    assert.Equal(t, uint(1), response[0].ID)
}
```

#### WebSocket Implementation

```go
// websocket_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/stretchr/testify/assert"
)

func TestWebsocketHandler(t *testing.T) {
    // Setup test server
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/ws", WebsocketListener)
    
    server := httptest.NewServer(router)
    defer server.Close()
    
    // Convert http URL to ws URL
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=1&user_id=1"
    
    // Connect to WebSocket
    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    assert.NoError(t, err)
    defer ws.Close()
    
    // Test sending a message
    testMessage := ChatMessage{
        Message: "Hello, testing!",
        Who: 1,
        MatchID: 1,
    }
    
    err = ws.WriteJSON(testMessage)
    assert.NoError(t, err)
    
    // Test receiving a message
    var receivedMsg ChatMessage
    err = ws.ReadJSON(&receivedMsg)
    assert.NoError(t, err)
    assert.Equal(t, testMessage.Message, receivedMsg.Message)
}
```

#### Notification System

```go
// notifications_test.go
package main

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

func TestUpdateUnreadCounts(t *testing.T) {
    // Setup mock DB
    originalDB := db
    mockDB := setupMockDB()
    db = mockDB
    defer func() { db = originalDB }()
    
    // Create test match
    match := Match{
        ID: 1,
        Offered: 1,
        Accepted: 2,
        UnreadOffered: 0,
        UnreadAccepted: 0,
    }
    
    // Setup mock expectations
    mockDB.ExpectQuery("SELECT (.+) FROM `matches` WHERE").
        WillReturnRows(sqlmock.NewRows([]string{"id", "offered", "accepted", "unread_offered", "unread_accepted"}).
        AddRow(match.ID, match.Offered, match.Accepted, match.UnreadOffered, match.UnreadAccepted))
    
    mockDB.ExpectExec("UPDATE `matches` SET").WillReturnResult(sqlmock.NewResult(1, 1))
    
    // Call function with sender ID = 1 (the offering user)
    updateUnreadCounts(1, 1)
    
    // Verify expectations
    err := mockDB.ExpectationsWereMet()
    assert.NoError(t, err)
    
    // Repeat for other scenarios...
}
```

### Backend Integration Testing

```go
// integration_test.go
package main

import (
    "os"
    "testing"
    
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func TestMain(m *testing.M) {
    // Setup test database
    testDB, _ = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    
    // Migrate schemas
    testDB.AutoMigrate(&User{}, &Person{}, &Match{}, &ChatMessage{})
    
    // Setup test router
    gin.SetMode(gin.TestMode)
    testRouter = gin.New()
    
    // Setup routes with the test database
    db = testDB
    setupRoutes(testRouter)
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    // ...
    
    os.Exit(code)
}

// Test suites follow...
```

### API Testing

```go
// api_test.go
func TestAPI(t *testing.T) {
    // Setup test server
    ts := httptest.NewServer(testRouter)
    defer ts.Close()
    
    // Create test client
    client := &http.Client{}
    
    // Login to get JWT
    loginReq := LoginRequest{
        Email: "test@example.com",
        Password: "password",
    }
    loginJSON, _ := json.Marshal(loginReq)
    
    req, _ := http.NewRequest("POST", ts.URL+"/login", bytes.NewBuffer(loginJSON))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    var loginResp struct {
        Token string `json:"token"`
    }
    json.NewDecoder(resp.Body).Decode(&loginResp)
    resp.Body.Close()
    
    // Test authenticated endpoint
    req, _ = http.NewRequest("GET", ts.URL+"/matches", nil)
    req.Header.Set("Authorization", loginResp.Token)
    
    resp, err = client.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // More API tests...
}
```

## Frontend Testing

### Component Testing

#### Chat Component Testing

```jsx
// ChatModal.test.jsx
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ChatModal from './chatmodal';

// Mock dependencies
jest.mock('../api', () => ({
  fetchChatHistory: jest.fn().mockResolvedValue([
    { id: 1, message: 'Hello', who: 'Person 1', time: '2023-01-01T12:00:00Z' }
  ])
}));

describe('ChatModal', () => {
  const mockProps = {
    match: { ID: 1, Offered: 1, Accepted: 2, UnreadOffered: 3 },
    person: { id: 2, name: 'Test User', profile: { url: 'test.jpg' } },
    User: { id: 1 },
    jwt: 'test-token',
    clearChatNotification: jest.fn()
  };

  test('displays unread message badge', () => {
    render(<ChatModal {...mockProps} unreadmessages={3} />);
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  test('calls clearChatNotification when opened', () => {
    render(<ChatModal {...mockProps} unreadmessages={3} />);
    userEvent.click(screen.getByText(/Chat with Test User/i));
    expect(mockProps.clearChatNotification).toHaveBeenCalledWith(1);
  });

  test('renders message history', async () => {
    render(<ChatModal {...mockProps} />);
    userEvent.click(screen.getByText(/Chat with Test User/i));
    
    // Wait for message to appear
    const message = await screen.findByText('Hello');
    expect(message).toBeInTheDocument();
  });
});
```

#### Notification Service Testing

```jsx
// NotificationService.test.jsx
import React from 'react';
import { render, waitFor } from '@testing-library/react';
import NotificationService from './notificationservice';

// Mock EventSource
class MockEventSource {
  constructor(url) {
    this.url = url;
    this.onmessage = null;
    this.onerror = null;
    
    // Simulate connection
    setTimeout(() => {
      if (this.onopen) this.onopen();
    }, 0);
  }
  
  // Method to simulate receiving a message
  simulateMessage(data) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
  
  close() {}
}

// Replace browser's EventSource with mock
global.EventSource = MockEventSource;

describe('NotificationService', () => {
  test('registers SSE connection and handles messages', async () => {
    // Create mock handlers
    const mockNewMessage = jest.fn();
    const mockNewMatch = jest.fn();
    
    // Render component
    const { unmount } = render(
      <NotificationService 
        jwt="test-token" 
        user={{ id: 1 }} 
        onNewMessage={mockNewMessage} 
        onNewMatch={mockNewMatch} 
      />
    );
    
    // Access the EventSource instance
    const eventSource = global.EventSource.mock.instances[0];
    
    // Simulate receiving a message notification
    eventSource.simulateMessage({
      type: 'chat_message',
      match_id: 5,
      count: 3,
      timestamp: '2023-01-01T12:00:00Z'
    });
    
    // Check that handler was called with correct args
    await waitFor(() => {
      expect(mockNewMessage).toHaveBeenCalledWith(5, 3);
    });
    
    // Cleanup
    unmount();
  });
});
```

### Frontend Integration Testing

```jsx
// App.integration.test.jsx
import React from 'react';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { rest } from 'msw';
import { setupServer } from 'msw/node';
import App from './App';

// Setup mock server
const server = setupServer(
  // Mock login endpoint
  rest.post('/login', (req, res, ctx) => {
    return res(
      ctx.json({
        token: 'fake-jwt-token',
        person: { id: 1, name: 'Test User' }
      })
    );
  }),
  
  // Mock matches endpoint
  rest.get('/matches/:id', (req, res, ctx) => {
    return res(
      ctx.json([
        { ID: 1, Offered: 1, Accepted: 2, UnreadAccepted: 3 }
      ])
    );
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

test('login and notification flow', async () => {
  render(<App />);
  
  // Login
  userEvent.type(screen.getByLabelText(/email/i), 'test@example.com');
  userEvent.type(screen.getByLabelText(/password/i), 'password123');
  userEvent.click(screen.getByRole('button', { name: /login/i }));
  
  // Wait for login to complete
  await screen.findByText(/Chat Selector/i);
  
  // Check notification badge appears when matches load
  await screen.findByText('3');
  
  // Click on chat
  userEvent.click(screen.getByText(/Chat with/i));
  
  // Verify notification cleared
  expect(screen.queryByText('3')).not.toBeInTheDocument();
});
```

### End-to-End Testing

```js
// cypress/integration/chat_notification_spec.js
describe('Chat Notifications', () => {
  beforeEach(() => {
    // Login before each test
    cy.visit('/');
    cy.get('input[name="email"]').type('test@example.com');
    cy.get('input[name="password"]').type('password123');
    cy.get('button[type="submit"]').click();
    
    // Wait for dashboard to load
    cy.url().should('include', '/dashboard');
  });
  
  it('shows notification badge when new message arrives', () => {
    // Intercept API calls
    cy.intercept('GET', '/matches/*', { fixture: 'matches.json' });
    cy.intercept('GET', '/chat/*', { fixture: 'chat-history.json' });
    
    // Check no notification badge initially
    cy.get('.badge').should('not.exist');
    
    // Simulate incoming message via SSE
    cy.window().then(win => {
      const event = new MessageEvent('message', {
        data: JSON.stringify({
          type: 'chat_message',
          match_id: 1,
          count: 3,
          timestamp: new Date().toISOString()
        })
      });
      win.dispatchEvent(event);
    });
    
    // Verify notification badge appears
    cy.get('.badge').should('contain', '3');
    
    // Click on chat
    cy.contains('Chat with Test User').click();
    
    // Verify notification badge is cleared
    cy.get('.badge').should('not.exist');
  });
});
```

## Special Testing Areas

### WebSocket Testing

Testing WebSockets requires special handling to simulate the bidirectional communication:

```go
// Backend WebSocket Testing (Go)
func TestWebSocketBroadcast(t *testing.T) {
    // Create multiple test clients
    client1 := createTestWebSocketClient(t, 1, 1)
    client2 := createTestWebSocketClient(t, 1, 2)
    
    // Send message from client 1
    message := ChatMessage{MatchID: 1, Who: 1, Message: "Test broadcast"}
    client1.WriteJSON(message)
    
    // Verify client 2 receives it
    var received ChatMessage
    err := client2.ReadJSON(&received)
    assert.NoError(t, err)
    assert.Equal(t, message.Message, received.Message)
}
```

```jsx
// Frontend WebSocket Testing (React)
import WS from 'jest-websocket-mock';

test('WebSocket connection and messaging', async () => {
  // Create mock WebSocket server
  const server = new WS('ws://localhost:8080');
  
  // Render component that connects to WebSocket
  render(<ChatModal {...mockProps} />);
  
  // Wait for client to connect
  await server.connected;
  
  // Simulate receiving a message
  server.send(JSON.stringify({
    id: 123,
    message: "Hello from server",
    who: 2,
    time: new Date().toISOString()
  }));
  
  // Verify message appears in UI
  expect(await screen.findByText('Hello from server')).toBeInTheDocument();
  
  // Close the mock server
  server.close();
});
```

### SSE Testing

For Server-Sent Events, we need to mock the EventSource API:

```jsx
// React SSE Testing
import { render, act } from '@testing-library/react';

// Mock EventSource before tests
const mockEventSource = {
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  close: jest.fn()
};

// Store original EventSource
const OriginalEventSource = global.EventSource;

beforeEach(() => {
  // Replace global EventSource with mock constructor
  global.EventSource = jest.fn(() => mockEventSource);
});

afterEach(() => {
  // Restore original EventSource
  global.EventSource = OriginalEventSource;
});

test('SSE notification handling', async () => {
  // Render component
  render(<NotificationService jwt="token" user={{id: 1}} />);
  
  // Verify EventSource was constructed with correct URL
  expect(global.EventSource).toHaveBeenCalledWith(
    'http://localhost:8080/notifications/1?token=token'
  );
  
  // Get the message handler
  const messageHandler = mockEventSource.addEventListener.mock.calls.find(
    call => call[0] === 'message'
  )[1];
  
  // Simulate receiving an event
  act(() => {
    messageHandler({
      data: JSON.stringify({
        type: 'chat_message',
        match_id: 1,
        count: 3
      })
    });
  });
  
  // Assertions about component behavior...
});
```

### AI Integration Testing

AI integration testing requires special handling due to the non-deterministic nature of AI responses:

```go
// Mock AI responses for testing
func TestGenerateDateSuggestion(t *testing.T) {
    // Setup mock AI client
    originalClient := geminiClient
    mockClient := &MockAIClient{}
    geminiClient = mockClient
    defer func() { geminiClient = originalClient }()
    
    // Setup expected response
    mockClient.On("GenerateContent", mock.Anything, mock.Anything).Return(
        &genai.GenerateContentResponse{
            Candidates: []*genai.Candidate{
                {
                    Content: &genai.Content{
                        Parts: []genai.Part{
                            genai.Text("I suggest a coffee at Central Park at 3pm."),
                        },
                    },
                },
            },
        }, 
        nil,
    )
    
    // Call the function
    result, err := GenerateDateSuggestion(1, []ChatMessage{})
    
    // Verify results
    assert.NoError(t, err)
    assert.Contains(t, result.Message, "coffee at Central Park")
}
```

### Database Testing

For database testing, use in-memory databases or testcontainers:

```go
// GORM with SQLite for testing
func setupTestDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Migrate schemas
    err = db.AutoMigrate(&User{}, &Person{}, &Match{}, &ChatMessage{})
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

// MongoDB with mongodb-memory-server
func setupTestMongoDB(t *testing.T) {
    // Start MongoDB memory server
    mongod, err := mongodbMemoryServer.Start()
    require.NoError(t, err)
    
    // Connect to the in-memory database
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongod.URI()))
    require.NoError(t, err)
    
    // Store original client
    originalClient := mongoClient
    
    // Replace with test client
    mongoClient = client
    
    // Return cleanup function
    return func() {
        client.Disconnect(context.Background())
        mongod.Stop()
        mongoClient = originalClient
    }
}
```

## Test Infrastructure

### CI/CD Integration

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17
      
      - name: Install dependencies
        run: |
          cd webserverpoc
          go mod download
      
      - name: Run tests
        run: |
          cd webserverpoc
          go test -v ./... -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: actions/upload-artifact@v2
        with:
          name: backend-coverage
          path: webserverpoc/coverage.out
  
  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Node
        uses: actions/setup-node@v2
        with:
          node-version: '16'
      
      - name: Install dependencies
        run: |
          cd react-vite-poc
          npm ci
      
      - name: Run tests
        run: |
          cd react-vite-poc
          npm test -- --coverage
      
      - name: Upload coverage
        uses: actions/upload-artifact@v2
        with:
          name: frontend-coverage
          path: react-vite-poc/coverage
  
  e2e-tests:
    runs-on: ubuntu-latest
    needs: [backend-tests, frontend-tests]
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up environment
        run: docker-compose up -d
      
      - name: Run Cypress tests
        uses: cypress-io/github-action@v2
        with:
          working-directory: react-vite-poc
          start: npm run dev
          wait-on: 'http://localhost:5173'
      
      - name: Store test artifacts
        uses: actions/upload-artifact@v2
        with:
          name: cypress-videos
          path: react-vite-poc/cypress/videos
```

## Implementation Plan

### Phase 1: Core Backend Tests (Week 1-2)

1. Set up testing frameworks and mocks
   - Install Testify, gomock, sqlmock
   - Create mock implementations for database and external services

2. Implement unit tests for critical components:
   - Authentication middleware tests
   - Match handling logic tests
   - Notification system tests

3. Develop integration tests:
   - API endpoint tests with httptest
   - Database interaction tests with test database

### Phase 2: Frontend Component Tests (Week 3-4)

1. Set up testing framework:
   - Configure Jest and React Testing Library
   - Setup MSW for API mocking

2. Implement component tests:
   - Chat component tests
   - Notification component tests
   - Authentication flow tests

3. Develop integration tests:
   - User journey tests from login to chat
   - WebSocket/SSE integration tests

### Phase 3: End-to-End and Special Area Tests (Week 5-6)

1. Configure Cypress and end-to-end testing:
   - Setup test environment with Docker
   - Develop key user journey tests

2. Implement special area tests:
   - WebSocket tests
   - SSE notification tests
   - AI integration tests

3. Setup CI/CD integration:
   - GitHub Actions workflow
   - Test coverage reporting
   - Automated test runs on PRs

### Phase 4: Test Coverage Improvements (Week 7-8)

1. Identify coverage gaps:
   - Run coverage analysis tools
   - Address untested code paths

2. Performance and load testing:
   - Setup k6 or JMeter tests
   - Test concurrent WebSocket connections
   - Test notification delivery at scale

3. Test documentation:
   - Document testing approach
   - Create test execution guides
   - Add testing requirements to contribution guidelines

## Test Implementation Tracking

### Test Coverage Targets

| Component Area           | Line Coverage Target | Function Coverage Target | Critical Path Coverage |
|-------------------------|---------------------|--------------------------|------------------------|
| Authentication System    | 80%                 | 90%                      | 100%                   |
| Real-time Chat          | 75%                 | 85%                      | 100%                   |
| Notification System     | 80%                 | 90%                      | 100%                   |
| Match Handling          | 70%                 | 80%                      | 100%                   |
| Profile Management      | 60%                 | 70%                      | 90%                    |
| AI Integration          | 60%                 | 80%                      | 90%                    |
| Frontend Components     | 70%                 | 80%                      | 90%                    |
| Frontend Integration    | 60%                 | N/A                      | 80%                    |
| End-to-End Flows        | N/A                 | N/A                      | 70%                    |

### Test Infrastructure Setup

#### Backend Test Setup (Go)

```bash
# Install required testing packages
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/mock
go get -u github.com/DATA-DOG/go-sqlmock
go get -u github.com/golang/mock/mockgen
go get -u github.com/alicebob/miniredis/v2

# Generate mocks for interfaces
go generate ./...

# Create test configuration file
cat << EOF > config.test.json
{
  "mysql": {
    "host": "localhost",
    "port": 3306,
    "user": "test",
    "password": "test",
    "database": "gowebserver_test"
  },
  "mongodb": {
    "host": "localhost",
    "port": 27017,
    "user": "test",
    "password": "test",
    "database": "urmid_test"
  }
}
EOF

# Run tests with coverage
go test -v ./... -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

#### Frontend Test Setup (React)

```bash
# Install required testing packages
npm install --save-dev @testing-library/react @testing-library/user-event @testing-library/jest-dom
npm install --save-dev msw jest-websocket-mock
npm install --save-dev cypress @cypress/code-coverage

# Initialize Jest configuration
npx jest --init

# Setup MSW for API mocking
npx msw init public/ --save

# Setup Cypress
npx cypress open

# Create test configuration
cat << EOF > .env.test
REACT_APP_API_URL=http://localhost:8080
REACT_APP_MOCK_ENABLED=true
EOF

# Run tests with coverage
npm test -- --coverage
```

### Dependencies Tracking

#### Backend Testing Dependencies

| Dependency              | Purpose                         | Version    | Installation Method |
|-------------------------|----------------------------------|------------|--------------------|
| testify                 | Assertion library               | v1.8.0+    | go get             |
| go-sqlmock              | Mock SQL database               | v1.5.0+    | go get             |
| mockgen                 | Generate interface mocks        | v1.6.0+    | go get             |
| miniredis              | Mock Redis server               | v2.x       | go get             |
| mongo-driver mocks      | Mock MongoDB interactions       | built-in   | N/A                |
| httptest                | HTTP request/response testing   | built-in   | N/A                |
| gomock                  | Mocking framework               | v1.6.0+    | go get             |

#### Frontend Testing Dependencies

| Dependency              | Purpose                              | Version    | Installation Method   |
|-------------------------|---------------------------------------|------------|----------------------|
| @testing-library/react  | React component testing              | ^13.0.0    | npm install --save-dev |
| @testing-library/user-event | Simulating user interactions     | ^14.0.0    | npm install --save-dev |
| jest                    | Test runner                          | ^29.0.0    | npm install --save-dev |
| msw                     | API mocking                          | ^0.49.0    | npm install --save-dev |
| jest-websocket-mock     | WebSocket mocking                    | ^2.3.0     | npm install --save-dev |
| cypress                 | End-to-end testing                   | ^10.0.0    | npm install --save-dev |
| @cypress/code-coverage  | Code coverage for E2E tests          | ^3.10.0    | npm install --save-dev |

### Test Prioritization Matrix

The following matrix helps prioritize which tests to implement first based on feature complexity and business importance:

| Feature Area             | Complexity (1-5) | Importance (1-5) | Risk Score (C×I) | Priority |
|--------------------------|-----------------|------------------|-----------------|----------|
| Authentication           | 3               | 5                | 15               | 1        |
| Real-time Chat Core      | 5               | 5                | 25               | 1        |
| SSE Notifications        | 4               | 4                | 16               | 1        |
| Match Creation/Handling  | 3               | 4                | 12               | 2        |
| Profile Management       | 2               | 3                | 6                | 3        |
| AI Date Suggestions      | 4               | 2                | 8                | 3        |
| WebSocket Error Handling | 4               | 4                | 16               | 1        |
| Photo Management         | 3               | 3                | 9                | 2        |
| UI Components            | 2               | 2                | 4                | 4        |
| Advanced Search          | 3               | 2                | 6                | 3        |

Priority: 1 = Highest (Implement immediately), 4 = Lowest

### Mocking Strategy

#### Backend Mocking

1. **Database Mocking**
   - Use `sqlmock` for SQL database interactions
   - Define common mock setup functions in `testutils` package
   - Example implementation:

```go
// testutils/db_mock.go
package testutils

import (
    "database/sql"
    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, err
    }
    
    dialector := mysql.New(mysql.Config{
        Conn: sqlDB,
        SkipInitializeWithVersion: true,
    })
    
    db, err := gorm.Open(dialector, &gorm.Config{})
    return db, mock, err
}
```

2. **External Service Mocking**
   - Create interface wrappers for external dependencies
   - Generate mocks using `mockgen` or manual implementation
   - Inject mocks through constructors or global variables with restore functions

```go
// Define interface
type AIServiceClient interface {
    GenerateContent(context.Context, string) (string, error)
}

// Real implementation
type GeminiClient struct {
    client *genai.Client
}

func (g *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
    // Real implementation
}

// Mock implementation for tests
type MockAIClient struct {
    mock.Mock
}

func (m *MockAIClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
    args := m.Called(ctx, prompt)
    return args.String(0), args.Error(1)
}
```

3. **WebSocket Testing**
   - Create in-memory WebSocket connection pairs
   - Implement test helpers for WebSocket assertions

```go
// testutils/websocket_mock.go

func NewWebSocketPair() (*websocket.Conn, *websocket.Conn, error) {
    // Create in-memory connection pair for testing
    // ...
}
```

#### Frontend Mocking

1. **API Mocking with MSW**
   - Define handlers for all API endpoints in `src/mocks/handlers.js`
   - Setup mock service worker in `src/mocks/browser.js`
   - Use consistent response formats across mocks

```javascript
// src/mocks/handlers.js
import { rest } from 'msw'

export const handlers = [
  rest.get('/matches/:id', (req, res, ctx) => {
    const { id } = req.params
    return res(
      ctx.status(200),
      ctx.json([
        { ID: 1, Offered: id, Accepted: 2, UnreadOffered: 3 }
      ])
    )
  }),
  
  rest.post('/chat/markread/:id', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({ message: 'Messages marked as read' })
    )
  }),
  
  // Additional endpoints...
]
```

2. **WebSocket and SSE Mocking**
   - Create reusable mock implementations
   - Define standard test events

```javascript
// src/tests/mocks/eventSourceMock.js
export class EventSourceMock {
  constructor(url) {
    this.url = url;
    this.readyState = EventSource.CONNECTING;
    this.onmessage = null;
    this.onerror = null;
    this.onopen = null;
    
    // Simulate successful connection
    setTimeout(() => {
      if (this.onopen) {
        this.readyState = EventSource.OPEN;
        this.onopen({ type: 'open' });
      }
    }, 0);
  }
  
  // Helper to simulate incoming messages
  mockMessage(data) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
  
  close() {
    this.readyState = EventSource.CLOSED;
  }
}
```

3. **Component Mocking**
   - Create standardized mock props for common components
   - Use consistent mock data across tests

```javascript
// src/tests/mocks/componentProps.js
export const mockUserData = {
  id: 1,
  name: 'Test User',
  email: 'test@example.com',
  is_admin: false,
};

export const mockMatch = {
  ID: 1,
  Offered: 1,
  Accepted: 2,
  UnreadOffered: 3,
  OfferedProfile: {
    id: 1,
    name: 'Current User',
    profile: { url: 'user1.jpg' }
  },
  AcceptedProfile: {
    id: 2,
    name: 'Match User',
    profile: { url: 'user2.jpg' }
  }
};

// Common mock props for frequently tested components
export const mockChatModalProps = {
  match: mockMatch,
  person: mockMatch.AcceptedProfile,
  User: mockUserData,
  jwt: 'fake-token',
  clearChatNotification: jest.fn()
};
```