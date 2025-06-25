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

    //Ensure the token isn't in our expired list
    for _, expiredToken := range ExpiredTokens {
        if token == expiredToken {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
            c.Abort()
            return
        }
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

var ExpiredTokens []string // Doing it the stupid way....for now. Wild that there's not a better solution out of the box than either generating a ton of api keys or keeping them in memory - neither of which work for a serverless application. Redis is a good answer but fiddly at POC stage.