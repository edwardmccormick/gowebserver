package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateSSEToken validates a JWT token for SSE connections
// Returns the token claims if valid, or an error if invalid
func ValidateSSEToken(tokenString string) (jwt.MapClaims, error) {
	// Remove 'Bearer ' prefix if it exists
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	
	// Validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	// Verify that the token is not expired
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return nil, fmt.Errorf("token expired")
		}
	}
	
	return claims, nil
}