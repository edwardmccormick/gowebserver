package gowebserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

	if isTokenRevoked(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		c.Abort()
		return
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return signingKey(), nil
	})
	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("token", parsedToken)

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if jti, exists := claims["jti"].(string); exists && jti != "" && sessionStore != nil {
			session, err := sessionStore.GetByJTI(c.Request.Context(), jti)
			if err != nil || session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session has expired"})
				c.Abort()
				return
			}
		}
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if userID, exists := claims["sub"].(float64); exists {
			c.Set("userID", uint(userID))
		}
	}

	c.Next()
}

// AdminMiddleware checks if the user is an admin after JWT authentication
func AdminMiddleware(c *gin.Context) {
	JwtMiddleware(c)
	if c.IsAborted() {
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to identify user"})
		c.Abort()
		return
	}

	userID := userIDValue.(uint)

	if authService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not initialized"})
		c.Abort()
		return
	}

	isAdmin, err := authService.IsAdmin(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		c.Abort()
		return
	}

	c.Next()
}
