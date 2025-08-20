package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("supersecretkey") // Use a secure random key in production!

func JwtMiddleware(c *gin.Context) {
    // First try to get token from Authorization header
    token := c.GetHeader("Authorization")
    
    // If not found in header, check if it's in the query parameters (for SSE)
    if token == "" {
        token = c.Query("token")
    }
    
    // If still not found, return error
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

    // Set the token in the context for other middlewares to use
    c.Set("token", parsedToken)
    
    // Extract user ID from claims
    if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
        if userID, exists := claims["sub"].(float64); exists {
            c.Set("userID", uint(userID))
        }
    }

    c.Next()
}

// AdminMiddleware checks if the user is an admin after JWT authentication
func AdminMiddleware(c *gin.Context) {
    // First ensure the user is authenticated with a valid JWT
    JwtMiddleware(c)
    
    // Check if the request was aborted by the JwtMiddleware
    if c.IsAborted() {
        return
    }
    
    // Get the user ID from the context
    userIDValue, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to identify user"})
        c.Abort()
        return
    }
    
    userID := userIDValue.(uint)
    
    // Find the user in the database
    var user User
    if result := db.First(&user, userID); result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        c.Abort()
        return
    }
    
    // Check if the user is an admin
    if !user.IsAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
        c.Abort()
        return
    }
    
    // User is an admin, proceed
    c.Next()
}

var ExpiredTokens []string // Doing it the stupid way....for now. Wild that there's not a better solution out of the box than either generating a ton of api keys or keeping them in memory - neither of which work for a serverless application. Redis is a good answer but fiddly at POC stage.