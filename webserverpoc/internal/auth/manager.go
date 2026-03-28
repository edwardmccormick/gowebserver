package auth

import "sync"

type TokenManager struct {
	secret []byte
	mu     sync.RWMutex
	tokens map[string]struct{}
}

func NewTokenManager(secret []byte) *TokenManager {
	return &TokenManager{
		secret: secret,
		tokens: make(map[string]struct{}),
	}
}

func (m *TokenManager) SigningKey() []byte {
	return m.secret
}

func (m *TokenManager) Revoke(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[token] = struct{}{}
}

func (m *TokenManager) IsRevoked(token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.tokens[token]
	return exists
}

func (m *TokenManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens = make(map[string]struct{})
}

func (m *TokenManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tokens)
}
