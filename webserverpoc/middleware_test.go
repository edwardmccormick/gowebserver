package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func createTestToken(userID uint, isAdmin bool) string {
	claims := jwt.MapClaims{
		"sub": float64(userID),
		"admin": isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JwtSecret)
	return tokenString
}

func TestJwtMiddlewareSimple(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add a test endpoint protected by JwtMiddleware
	router.GET("/protected", JwtMiddleware, func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	
	// Test with no token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)
	
	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Test with valid token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	token := createTestToken(1, false)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)
	
	// Should return OK
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["userID"])
}