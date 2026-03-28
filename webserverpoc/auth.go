package gowebserver

import internalauth "github.com/edwardmccormick/gowebserver/internal/auth"

var tokenManager = internalauth.NewTokenManager([]byte("supersecretkey"))

func signingKey() []byte {
	return tokenManager.SigningKey()
}

func revokeToken(token string) {
	tokenManager.Revoke(token)
}

func isTokenRevoked(token string) bool {
	return tokenManager.IsRevoked(token)
}

func resetRevokedTokensForTests() {
	tokenManager.Reset()
}

func revokedTokenCountForTests() int {
	return tokenManager.Count()
}
