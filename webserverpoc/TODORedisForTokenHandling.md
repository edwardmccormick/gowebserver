# Scalable JWT Revocation with Redis Blacklist

## Why Use an In-Memory Cache for Token Blacklisting?

1. **Low Latency:**
   - In-memory caches like Redis are extremely fast and can handle millions of requests per second with minimal latency.

2. **Automatic Expiration:**
   - You can set an expiration time for blacklisted tokens in Redis, so they are automatically removed when the token expires.

3. **Scalability:**
   - Redis is highly scalable and can be deployed in a distributed manner, making it suitable for cloud-native applications.

4. **Reduced API Calls:**
   - Instead of querying a database or making multiple API calls, you can check the blacklist in Redis with a single lightweight operation.

---

## Implementation Strategy

### 1. Add Redis to Your Project
Install a Redis client for Go, such as [go-redis](https://github.com/redis/go-redis).

```bash
go get github.com/redis/go-redis/v9
```

### 2. Initialize Redis in main.go
Set up a Redis client in your application.

```golang
package main

import (
    "context"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

func init() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Replace with your Redis server address
    })
}
```

### 3. Update the Signout Handler to Blacklist Tokens
When a user logs out, add their JWT to the Redis blacklist with an expiration time.

```golang
func Signout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
        return
    }

    // Parse and validate the token
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return JwtSecret, nil
    })
    if err != nil || !parsedToken.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    claims := parsedToken.Claims.(jwt.MapClaims)
    exp := int64(claims["exp"].(float64)) // Extract the token expiration time

    // Add the token to the Redis blacklist with an expiration time
    err = redisClient.Set(ctx, token, "blacklisted", time.Until(time.Unix(exp, 0))).Err()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not blacklist token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}
```

### 4. Update Middleware to Check the Blacklist
Before processing a request, check if the token is in the Redis blacklist.
```go
func JwtMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
        c.Abort()
        return
    }

    // Check if the token is blacklisted
    isBlacklisted, err := redisClient.Get(ctx, token).Result()
    if err == nil && isBlacklisted == "blacklisted" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
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
```


---

## Advantages of This Approach

1. **Efficient Blacklist Checks:**
   - Redis provides fast lookups, so checking if a token is blacklisted is efficient.

2. **Automatic Cleanup:**
   - Blacklisted tokens are automatically removed from Redis when they expire, reducing manual maintenance.

3. **Scalable:**
   - Redis can handle high traffic and is suitable for distributed cloud-native applications.

4. **Reduced API Calls:**
   - The blacklist check is a single Redis operation, avoiding the need for additional API calls or database queries.

---

## Example Workflow

### 1. User Logs Out
- The `Signout` handler adds the user's JWT to the Redis blacklist with an expiration time.

### 2. User Makes a Request with the Blacklisted Token
- The `JwtMiddleware` checks the Redis blacklist and rejects the request if the token is blacklisted.

### 3. Token Expires
- Redis automatically removes the token from the blacklist when it expires.

---

## Next Steps

1. **Deploy Redis:**
   - Set up a Redis instance (e.g., using AWS ElastiCache, Azure Cache for Redis, or a self-hosted Redis server).

2. **Test the Implementation:**
   - Verify that blacklisted tokens are rejected and that expired tokens are automatically removed from Redis.

3. **Optimize for Production:**
   - Use a secure Redis setup with authentication and encryption for production environments.

This approach provides a scalable and efficient way to handle token revocation without introducing excessive API calls or database queries.