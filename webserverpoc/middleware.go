package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("supersecretkey") // Use a secure random key in production!

func JwtMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
        c.Abort()
        return
    }

    // Validate the token
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return JwtSecret, nil
    })
    if err != nil || !parsedToken.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        c.Abort()
        return
    }

    c.Next()
}

var RefreshTokens = make(map[string]string) // Map user ID to refresh token