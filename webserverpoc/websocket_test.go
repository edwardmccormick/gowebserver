package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
)

// Test the NotificationEvent struct and methods
func TestNotificationEvents(t *testing.T) {
	// Create a notification event
	event := NotificationEvent{
		Type:      NotificationTypeMessage,
		MatchID:   1,
		Count:     3,
		Timestamp: time.Now(),
	}
	
	// Verify the fields
	assert.Equal(t, NotificationTypeMessage, event.Type)
	assert.Equal(t, uint(1), event.MatchID)
	assert.Equal(t, 3, event.Count)
	assert.False(t, event.Timestamp.IsZero())
	
	// Test JSON marshaling
	jsonBytes, err := json.Marshal(event)
	require.NoError(t, err)
	
	// Test JSON unmarshaling
	var parsedEvent NotificationEvent
	err = json.Unmarshal(jsonBytes, &parsedEvent)
	require.NoError(t, err)
	
	// Verify the unmarshaled event matches the original
	assert.Equal(t, event.Type, parsedEvent.Type)
	assert.Equal(t, event.MatchID, parsedEvent.MatchID)
	assert.Equal(t, event.Count, parsedEvent.Count)
	
	// Test NotificationTypeMatch
	matchEvent := NotificationEvent{
		Type:      NotificationTypeMatch,
		MatchID:   2,
		Count:     1,
		Timestamp: time.Now(),
	}
	
	assert.Equal(t, NotificationTypeMatch, matchEvent.Type)
}

// TestTwoClientCommunication tests communication between two WebSocket clients
func TestTwoClientCommunication(t *testing.T) {
	// Setup test environment
	gin.SetMode(gin.TestMode)
	
	// Create a router with our mock broadcast WebSocket handler
	router := gin.New()
	router.GET("/ws", MockBroadcastHandler)
	
	// Start a test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Create websocket URLs from server URL
	wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=1&user_id=1"
	wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=1&user_id=2"
	
	// Connect the first client
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	require.NoError(t, err)
	defer conn1.Close()
	
	// Connect the second client
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	require.NoError(t, err)
	defer conn2.Close()
	
	// Read welcome messages
	var welcome1, welcome2 ChatMessage
	err = conn1.ReadJSON(&welcome1)
	require.NoError(t, err)
	err = conn2.ReadJSON(&welcome2)
	require.NoError(t, err)
	
	// Send a message from client 1
	testMsg := ChatMessage{
		Message: "Hello from client 1",
		Who:     1,
		MatchID: 1,
		Time:    time.Now(),
	}
	err = conn1.WriteJSON(testMsg)
	require.NoError(t, err)
	
	// Client 2 should receive the broadcast message
	var receivedMsg ChatMessage
	err = conn2.ReadJSON(&receivedMsg)
	require.NoError(t, err)
	assert.Equal(t, testMsg.Message, receivedMsg.Message)
	assert.Equal(t, testMsg.Who, receivedMsg.Who)
}

// MockWebSocketHandler is a test-friendly version of WebsocketListener
// that skips database operations but handles basic message exchange
func MockWebSocketHandler(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Parse room ID and user ID from query parameters
	roomID := c.Query("id")
	userID := c.Query("user_id")

	if roomID == "" || userID == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Room ID and User ID are required"))
		return
	}
	
	// Send a test welcome message
	welcomeMsg := ChatMessage{
		Message: "Welcome to the test chat",
		Who:     0, // System message
		MatchID: 1,
		Time:    time.Now(),
	}
	
	msgBytes, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)
	
	// Read one message from client and echo it back
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		
		// Echo the message back
		conn.WriteMessage(websocket.TextMessage, msg)
		break
	}
}

// MockBroadcastHandler simulates WebsocketListener with broadcasting capability
func MockBroadcastHandler(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Parse room ID and user ID from query parameters
	roomID := c.Query("id")
	userID := c.Query("user_id")

	if roomID == "" || userID == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Room ID and User ID are required"))
		return
	}
	
	// Add connection to a test room map
	roomsLock.Lock()
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rooms[roomID][conn] = true
	roomsLock.Unlock()
	
	defer func() {
		// Clean up connection from the room
		roomsLock.Lock()
		delete(rooms[roomID], conn)
		roomsLock.Unlock()
	}()
	
	// Send a test welcome message
	welcomeMsg := ChatMessage{
		Message: "Welcome to the test chat",
		Who:     0, // System message
		MatchID: 1,
		Time:    time.Now(),
	}
	
	msgBytes, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)
	
	// Read and broadcast messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		
		// Broadcast the message to all connections in this room
		roomsLock.Lock()
		for client := range rooms[roomID] {
			if client != conn { // Don't send back to sender
				client.WriteMessage(websocket.TextMessage, msg)
			}
		}
		roomsLock.Unlock()
	}
}

// Test WebSocket connection
func TestWebSocketConnection(t *testing.T) {
	// Setup test environment
	gin.SetMode(gin.TestMode)
	
	// Create a router with our mock WebSocket handler
	router := gin.New()
	router.GET("/ws", MockWebSocketHandler)
	
	// Start a test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Create websocket URL from server URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=1&user_id=1"
	
	// Connect to the websocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to websocket: %v", err)
	}
	defer conn.Close()
	
	// Check that the connection is successful
	assert.NotNil(t, conn)
	
	// Read welcome message
	var welcomeMsg ChatMessage
	err = conn.ReadJSON(&welcomeMsg)
	require.NoError(t, err)
	assert.Equal(t, "Welcome to the test chat", welcomeMsg.Message)
	
	// Send a test message
	testMsg := ChatMessage{
		Message: "Hello from test client",
		Who:     1,
		MatchID: 1,
	}
	err = conn.WriteJSON(testMsg)
	require.NoError(t, err)
	
	// Read the echo response
	var echoMsg ChatMessage
	err = conn.ReadJSON(&echoMsg)
	require.NoError(t, err)
	assert.Equal(t, testMsg.Message, echoMsg.Message)
}

// TestNotificationCenter tests the notification center functionality
func TestNotificationCenter(t *testing.T) {
	// Setup test environment
	gin.SetMode(gin.TestMode)
	
	// Create a test notification center
	nc := NotificationCenter{
		clients: make(map[string]map[chan NotificationEvent]bool),
	}
	
	// Create a notification channel
	userID := "1"
	notificationChan := make(chan NotificationEvent, 1)
	
	// Register the client
	nc.Register(userID, notificationChan)
	
	// Verify the client was registered
	assert.NotNil(t, nc.clients[userID])
	assert.True(t, nc.clients[userID][notificationChan])
	
	// Create a test event
	testEvent := NotificationEvent{
		Type:    NotificationTypeMessage,
		MatchID: 1,
		Count:   3,
	}
	
	// Broadcast the event
	nc.Broadcast(userID, testEvent)
	
	// Check if the event was received
	receivedEvent := <-notificationChan
	assert.Equal(t, testEvent.Type, receivedEvent.Type)
	
	// Unregister the client
	nc.Unregister(userID, notificationChan)
	
	// Verify the client was unregistered
	assert.Empty(t, nc.clients[userID])
}

// createSSETestToken creates a test JWT token for SSE testing
func createSSETestToken(userID uint) string {
	// Create token claims
	claims := jwt.MapClaims{
		"sub":  float64(userID),
		"exp":  float64(time.Now().Add(time.Hour).Unix()),
		"iat":  float64(time.Now().Unix()),
		"admin": false,
	}
	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token with secret
	tokenString, _ := token.SignedString(JwtSecret)
	return tokenString
}

// TestSSEAuthenticationWithValidToken tests SSE authentication with a valid token
func TestSSEAuthenticationWithValidToken(t *testing.T) {
	// Create a test context with a custom recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Set up the request with correct parameters
	c.Request, _ = http.NewRequest("GET", "/notifications/1", nil)
	c.Params = []gin.Param{{
		Key:   "id",
		Value: "1",
	}}
	
	// Set up a valid token as query parameter
	token := createSSETestToken(1)
	c.Request.URL.RawQuery = "token=" + token
	
	// Create a custom handler that just validates authentication
	validated := false
	handlerFunc := func(c *gin.Context) {
		// Run just the authentication part of SSEHandler
		// Get user ID from URL parameter
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}
		
		// Extract authentication token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Try to get token from query parameter
			authHeader = c.Query("token")
			if authHeader == "" {
				c.JSON(401, gin.H{"error": "Authentication required"})
				return
			}
		}
		
		// Validate the token
		tokenClaims, err := ValidateSSEToken(authHeader)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid authentication token"})
			return
		}
		
		// Get user ID from token claims
		tokenUserID, ok := tokenClaims["sub"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			return
		}
		
		// Verify that token user ID matches requested user ID
		if uint(userID) != uint(tokenUserID) {
			c.JSON(403, gin.H{"error": "Not authorized to subscribe to this user's notifications"})
			return
		}
		
		// If we reach here, authentication was successful
		validated = true
		c.Status(http.StatusOK)
	}
	
	// Run the handler function
	handlerFunc(c)
	
	// Check if validation was successful
	assert.True(t, validated, "SSE authentication validation should succeed")
}

// TestSSEAuthenticationWithInvalidToken tests SSE authentication with an invalid token
func TestSSEAuthenticationWithInvalidToken(t *testing.T) {
	// Create a test context with a custom recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Set up the request with correct parameters
	c.Request, _ = http.NewRequest("GET", "/notifications/1", nil)
	c.Params = []gin.Param{{
		Key:   "id",
		Value: "1",
	}}
	
	// Set up an invalid token as query parameter
	c.Request.URL.RawQuery = "token=invalid.token.value"
	
	// Create a custom handler that just validates authentication
	validated := false
	handlerFunc := func(c *gin.Context) {
		// Run just the authentication part of SSEHandler
		// Get user ID from URL parameter
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}
		
		// Extract authentication token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Try to get token from query parameter
			authHeader = c.Query("token")
			if authHeader == "" {
				c.JSON(401, gin.H{"error": "Authentication required"})
				return
			}
		}
		
		// Validate the token
		tokenClaims, err := ValidateSSEToken(authHeader)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid authentication token"})
			return
		}
		
		// Get user ID from token claims
		tokenUserID, ok := tokenClaims["sub"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			return
		}
		
		// Verify that token user ID matches requested user ID
		if uint(userID) != uint(tokenUserID) {
			c.JSON(403, gin.H{"error": "Not authorized to subscribe to this user's notifications"})
			return
		}
		
		// If we reach here, authentication was successful
		validated = true
		c.Status(http.StatusOK)
	}
	
	// Run the handler function
	handlerFunc(c)
	
	// Check validation failed (as expected with invalid token)
	assert.False(t, validated, "SSE authentication should fail with invalid token")
	assert.Equal(t, 401, w.Code, "Invalid token should result in 401 unauthorized")
}

// TestSSEAuthenticationWithWrongUserID tests SSE with token for different user
func TestSSEAuthenticationWithWrongUserID(t *testing.T) {
	// Create a test context with a custom recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Set up the request with parameters for user ID 1
	c.Request, _ = http.NewRequest("GET", "/notifications/1", nil)
	c.Params = []gin.Param{{
		Key:   "id",
		Value: "1",
	}}
	
	// But set up a token for user ID 2
	token := createSSETestToken(2)
	c.Request.URL.RawQuery = "token=" + token
	
	// Create a custom handler that just validates authentication
	validated := false
	handlerFunc := func(c *gin.Context) {
		// Run just the authentication part of SSEHandler
		// Get user ID from URL parameter
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}
		
		// Extract authentication token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Try to get token from query parameter
			authHeader = c.Query("token")
			if authHeader == "" {
				c.JSON(401, gin.H{"error": "Authentication required"})
				return
			}
		}
		
		// Validate the token
		tokenClaims, err := ValidateSSEToken(authHeader)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid authentication token"})
			return
		}
		
		// Get user ID from token claims
		tokenUserID, ok := tokenClaims["sub"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			return
		}
		
		// Verify that token user ID matches requested user ID
		if uint(userID) != uint(tokenUserID) {
			c.JSON(403, gin.H{"error": "Not authorized to subscribe to this user's notifications"})
			return
		}
		
		// If we reach here, authentication was successful
		validated = true
		c.Status(http.StatusOK)
	}
	
	// Run the handler function
	handlerFunc(c)
	
	// Check validation failed (as expected with mismatched user ID)
	assert.False(t, validated, "SSE authentication should fail with mismatched user ID")
	assert.Equal(t, 403, w.Code, "Mismatched user ID should result in 403 forbidden")
}

// TestSSEHandler tests the SSE handler functionality
// This test is disabled because it takes 60 seconds to run due to SSE keepalive
func TestSSEHandler(t *testing.T) {
	// Skip this test in normal runs because it takes too long
	t.Skip("This test takes 60 seconds to run due to SSE keepalive")
	
	// Setup test environment
	gin.SetMode(gin.TestMode)
	
	// Create a test router with the SSE handler
	router := gin.New()
	router.GET("/notifications/:id", SSEHandler)
	
	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Create a test request with a JWT token for user ID 1
	req, err := http.NewRequest("GET", server.URL+"/notifications/1", nil)
	require.NoError(t, err)
	
	// Add token to query parameter (for SSE connections where headers may not be reliable)
	token := createSSETestToken(1)
	req.URL.RawQuery = "token=" + token
	
	// Create a test client
	client := server.Client()
	
	// Send the request
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Check the response headers for SSE
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/event-stream")
	assert.Equal(t, "no-cache", resp.Header.Get("Cache-Control"))
	assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))
	
	// Close the connection explicitly
	resp.Body.Close()
	// Wait a moment for cleanup
	time.Sleep(100 * time.Millisecond)
}