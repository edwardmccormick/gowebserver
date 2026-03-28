package gowebserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/realtime"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func signJWTForTests(t *testing.T, userID uint) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(signingKey())
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return tokenString
}

func TestHashPassword(t *testing.T) {
	hash := HashPassword("password123")
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "password123" {
		t.Fatal("expected password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("password123")); err != nil {
		t.Fatalf("expected hash to validate: %v", err)
	}
}

func TestDeepCopyPhotoArray(t *testing.T) {
	original := []ProfilePhoto{
		{S3Key: "1/image1", Caption: "first"},
		{S3Key: "1/image2", Caption: "second"},
	}

	cloned := deepCopyPhotoArray(original)
	cloned[0].Caption = "changed"

	if original[0].Caption != "first" {
		t.Fatalf("expected original slice to remain unchanged, got %q", original[0].Caption)
	}
}

func TestCalculateHaversineDistance(t *testing.T) {
	if got := CalculateHaversineDistance(29.52, -98.49, 29.52, -98.49); got != 0 {
		t.Fatalf("expected zero distance, got %f", got)
	}

	got := CalculateHaversineDistance(29.52016959410149, -98.49401109752402, 29.453596593823395, -98.47166788793534)
	if got <= 0 {
		t.Fatalf("expected positive distance, got %f", got)
	}
	if got < 4 || got > 6 {
		t.Fatalf("expected distance in a reasonable range, got %f", got)
	}
}

func TestCalculateBoundingBox(t *testing.T) {
	minLat, maxLat, minLong, maxLong := CalculateBoundingBox(29.42, -98.49, 10)
	if minLat >= maxLat {
		t.Fatalf("expected latitude bounds to expand, got %f >= %f", minLat, maxLat)
	}
	if minLong >= maxLong {
		t.Fatalf("expected longitude bounds to expand, got %f >= %f", minLong, maxLong)
	}
	if !(29.42 > minLat && 29.42 < maxLat) {
		t.Fatalf("expected original latitude to remain inside bounds: %f..%f", minLat, maxLat)
	}
	if !(-98.49 > minLong && -98.49 < maxLong) {
		t.Fatalf("expected original longitude to remain inside bounds: %f..%f", minLong, maxLong)
	}
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	configJSON := `{
		"mysql": {"host": "localhost", "port": 3306, "user": "urmid", "password": "pw", "database": "urmid"},
		"mongodb": {"host": "localhost", "port": 27017, "user": "root", "password": "pw", "database": "urmid"}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.MySQL.Host != "localhost" || cfg.Mongo.Port != 27017 {
		t.Fatalf("unexpected config contents: %+v", cfg)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestDetailToText(t *testing.T) {
	text := DetailToText("dogs", 10)
	if !strings.Contains(text, "dogs") && !strings.Contains(text, "golden retriever") {
		t.Fatalf("expected dog description, got %q", text)
	}

	if got := DetailToText("cats", 7); got != "7" {
		t.Fatalf("expected numeric fallback, got %q", got)
	}
}

func TestFormatDetailsForPrompt(t *testing.T) {
	person := Person{
		Details: Details{
			Dogs:        8,
			Cats:        4,
			Drinking:    6,
			EnergyLevel: 7,
		},
	}

	formatted := formatDetailsForPrompt(person)
	for _, expected := range []string{"Details:", "Dogs:", "Cats:", "Drinking:", "Energy Level:"} {
		if !strings.Contains(formatted, expected) {
			t.Fatalf("expected %q in formatted details: %q", expected, formatted)
		}
	}
}

func TestSortMessagesByTime(t *testing.T) {
	later := time.Now()
	earlier := later.Add(-time.Minute)
	messages := []ChatMessage{
		{ID: 2, Time: later},
		{ID: 1, Time: earlier},
	}

	sortMessagesByTime(messages)

	if messages[0].ID != 1 || messages[1].ID != 2 {
		t.Fatalf("expected messages sorted by time, got %+v", messages)
	}
}

func TestNotificationCenterBroadcast(t *testing.T) {
	nc := realtime.NewNotificationCenter()
	ch := make(chan realtime.NotificationEvent, 1)

	nc.Register("5", ch)
	nc.Broadcast("5", realtime.NotificationEvent{Type: realtime.NotificationTypeMessage, MatchID: 9, Count: 2})

	select {
	case event := <-ch:
		if event.MatchID != 9 || event.Count != 2 {
			t.Fatalf("unexpected event payload: %+v", event)
		}
	default:
		t.Fatal("expected notification to be delivered")
	}

	nc.Unregister("5", ch)
}

func TestJwtMiddlewareRejectsMissingToken(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request = req

	JwtMiddleware(ctx)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
	if !ctx.IsAborted() {
		t.Fatal("expected context to be aborted")
	}
}

func TestJwtMiddlewareAcceptsAuthorizationHeader(t *testing.T) {
	token := signJWTForTests(t, 42)
	calledNext := false
	router := gin.New()
	router.Use(JwtMiddleware)
	router.GET("/", func(c *gin.Context) {
		calledNext = true
		userID, exists := c.Get("userID")
		if !exists || userID.(uint) != 42 {
			t.Fatalf("expected userID 42 in context, got %+v (exists=%v)", userID, exists)
		}
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !calledNext {
		t.Fatal("expected middleware chain to continue")
	}
}

func TestJwtMiddlewareAcceptsQueryToken(t *testing.T) {
	token := signJWTForTests(t, 7)
	router := gin.New()
	router.Use(JwtMiddleware)
	router.GET("/notifications/7", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists || userID.(uint) != 7 {
			t.Fatalf("expected userID 7 in context, got %+v (exists=%v)", userID, exists)
		}
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications/7?token="+token, nil)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestJwtMiddlewareRejectsExpiredTokenList(t *testing.T) {
	token := signJWTForTests(t, 9)
	resetRevokedTokensForTests()
	t.Cleanup(resetRevokedTokensForTests)
	revokeToken(token)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	ctx.Request = req

	JwtMiddleware(ctx)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestLogout(t *testing.T) {
	t.Run("missing token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/logout", nil)

		Logout(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", recorder.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/logout", nil)
		req.Header.Set("Authorization", "not-a-token")
		ctx.Request = req

		Logout(ctx)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", recorder.Code)
		}
	})

	t.Run("valid token is expired", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": 5,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(signingKey())
		if err != nil {
			t.Fatalf("sign token: %v", err)
		}

		resetRevokedTokensForTests()
		t.Cleanup(resetRevokedTokensForTests)

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/logout", nil)
		req.Header.Set("Authorization", tokenString)
		ctx.Request = req

		Logout(ctx)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Code)
		}
		if !isTokenRevoked(tokenString) || revokedTokenCountForTests() != 1 {
			t.Fatalf("expected token to be revoked")
		}
	})
}

func TestLoadChatHistoryFromMongoWithoutClient(t *testing.T) {
	previous := chatService
	chatService = nil
	t.Cleanup(func() {
		chatService = previous
	})

	_, err := loadChatHistoryFromMongo(1)
	if err == nil {
		t.Fatal("expected error when mongo client is nil")
	}
}

func TestDumpChatHistoryToMongoWithoutClient(t *testing.T) {
	previous := chatService
	chatService = nil
	t.Cleanup(func() {
		chatService = previous
	})

	err := dumpChatHistoryToMongo(Conversation{MatchID: 1})
	if err == nil {
		t.Fatal("expected error when mongo client is nil")
	}
}

func TestSignupRejectsInvalidJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	Signup(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestLoginRejectsInvalidJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	Login(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestLogoutResponseBodyIsJSON(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 8,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(signingKey())
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	resetRevokedTokensForTests()
	t.Cleanup(resetRevokedTokensForTests)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req.Header.Set("Authorization", tokenString)
	ctx.Request = req

	Logout(ctx)

	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response body: %v", err)
	}
	if body["message"] == "" {
		t.Fatalf("expected success message in response body, got %+v", body)
	}
}

func TestNewRouterRegistersAuthAndProtectedRoutes(t *testing.T) {
	router := newRouter()

	t.Run("logout route is registered", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/logout", nil)

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected registered logout route to return 400 for missing token, got %d", recorder.Code)
		}
	})

	t.Run("notifications route is protected", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/notifications/5", nil)

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected protected notifications route to return 401, got %d", recorder.Code)
		}
	})
}
