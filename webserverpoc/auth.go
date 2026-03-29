package gowebserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	internalauth "github.com/edwardmccormick/gowebserver/internal/auth"
	"github.com/edwardmccormick/gowebserver/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

var tokenManager = internalauth.NewTokenManager([]byte("supersecretkey"))

func signingKey() []byte {
	return tokenManager.SigningKey()
}

func revokeToken(token string) {
	if jti, ok := tokenJTI(token); ok && sessionStore != nil {
		_ = sessionStore.RevokeByJTI(context.Background(), jti, time.Now())
		return
	}
	tokenManager.Revoke(token)
}

func isTokenRevoked(token string) bool {
	if jti, ok := tokenJTI(token); ok && sessionStore != nil {
		session, err := sessionStore.GetByJTI(context.Background(), jti)
		if err != nil {
			return true
		}
		if session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
			return true
		}
		return false
	}
	return tokenManager.IsRevoked(token)
}

func resetRevokedTokensForTests() {
	tokenManager.Reset()
}

func revokedTokenCountForTests() int {
	return tokenManager.Count()
}

func createSessionToken(userID uint, email string, expiresAt time.Time) (string, error) {
	jti, err := newTokenID()
	if err != nil {
		return "", err
	}

	if sessionStore != nil {
		if err := sessionStore.Create(context.Background(), store.Session{
			UserID:    userID,
			JTI:       jti,
			ExpiresAt: expiresAt,
		}); err != nil {
			return "", err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"exp":   expiresAt.Unix(),
		"email": email,
		"jti":   jti,
	})

	return token.SignedString(signingKey())
}

func tokenJTI(token string) (string, bool) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", false
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	jti, ok := claims["jti"].(string)
	return jti, ok && jti != ""
}

func newTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
