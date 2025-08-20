package testutils

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// WebSocketTestClient represents a WebSocket test client for simulating WebSocket connections
type WebSocketTestClient struct {
	*websocket.Conn
	T *testing.T
}

// MockEventSource mocks the browser's EventSource for Server-Sent Events testing
type MockEventSource struct {
	Events      []map[string]interface{}
	OnOpen      func()
	OnMessage   func(data map[string]interface{})
	OnError     func(err error)
	IsConnected bool
}

// NewMockEventSource creates a new mock EventSource
func NewMockEventSource() *MockEventSource {
	return &MockEventSource{
		Events:      []map[string]interface{}{},
		IsConnected: true,
	}
}

// SendEvent simulates receiving an SSE event
func (m *MockEventSource) SendEvent(data map[string]interface{}) {
	m.Events = append(m.Events, data)
	if m.OnMessage != nil {
		m.OnMessage(data)
	}
}

// Close simulates closing the EventSource connection
func (m *MockEventSource) Close() {
	m.IsConnected = false
}

// Open simulates opening the EventSource connection
func (m *MockEventSource) Open() {
	m.IsConnected = true
	if m.OnOpen != nil {
		m.OnOpen()
	}
}

// CreateWebSocketTestPair creates a pair of connected WebSocket clients for testing
func CreateWebSocketTestPair(t *testing.T, handler func(c *gin.Context), matchID, userID uint) (*WebSocketTestClient, *WebSocketTestClient) {
	// Setup test server
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ws", handler)

	server := httptest.NewServer(router)
	defer server.Close()

	// Create WebSocket URL for client 1 (user1)
	wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=" + 
		strconv.FormatUint(uint64(matchID), 10) + "&user_id=" + strconv.FormatUint(uint64(userID), 10)

	// Create WebSocket URL for client 2 (user2)
	wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?id=" + 
		strconv.FormatUint(uint64(matchID), 10) + "&user_id=" + strconv.FormatUint(uint64(userID+1), 10)

	// Connect both clients
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	require.NoError(t, err)

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	require.NoError(t, err)

	client1 := &WebSocketTestClient{Conn: conn1, T: t}
	client2 := &WebSocketTestClient{Conn: conn2, T: t}

	return client1, client2
}

// SendJSON sends a JSON message through the WebSocket connection
func (c *WebSocketTestClient) SendJSON(v interface{}) {
	err := c.Conn.WriteJSON(v)
	require.NoError(c.T, err)
}

// ReceiveJSON receives a JSON message from the WebSocket connection
func (c *WebSocketTestClient) ReceiveJSON(v interface{}) {
	err := c.Conn.ReadJSON(v)
	require.NoError(c.T, err)
}

// Close closes the WebSocket connection
func (c *WebSocketTestClient) Close() {
	c.Conn.Close()
}

// SSEResponseRecorder is a custom ResponseRecorder for SSE testing
type SSEResponseRecorder struct {
	*httptest.ResponseRecorder
	closeNotify chan bool
}

// NewSSEResponseRecorder creates a custom response recorder for SSE testing
func NewSSEResponseRecorder() *SSEResponseRecorder {
	return &SSEResponseRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closeNotify:      make(chan bool, 1),
	}
}

// CloseNotify implements http.CloseNotifier interface
func (r *SSEResponseRecorder) CloseNotify() <-chan bool {
	return r.closeNotify
}

// Close simulates closing the connection
func (r *SSEResponseRecorder) Close() {
	r.closeNotify <- true
}

// CreateTestContext creates a Gin context with SSE support for testing
func CreateTestContext(t *testing.T) (*gin.Context, *SSEResponseRecorder) {
	recorder := NewSSEResponseRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	
	return c, recorder
}