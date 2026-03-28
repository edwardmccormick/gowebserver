package auth

import "testing"

func TestTokenManagerSigningKey(t *testing.T) {
	manager := NewTokenManager([]byte("secret"))
	if string(manager.SigningKey()) != "secret" {
		t.Fatalf("expected signing key to round-trip, got %q", string(manager.SigningKey()))
	}
}

func TestTokenManagerRevokeAndReset(t *testing.T) {
	manager := NewTokenManager([]byte("secret"))

	if manager.IsRevoked("token-1") {
		t.Fatal("did not expect token to start revoked")
	}

	manager.Revoke("token-1")
	manager.Revoke("token-2")

	if !manager.IsRevoked("token-1") || !manager.IsRevoked("token-2") {
		t.Fatal("expected revoked tokens to be tracked")
	}
	if manager.Count() != 2 {
		t.Fatalf("expected 2 revoked tokens, got %d", manager.Count())
	}

	manager.Reset()

	if manager.IsRevoked("token-1") || manager.IsRevoked("token-2") {
		t.Fatal("expected reset to clear revoked tokens")
	}
	if manager.Count() != 0 {
		t.Fatalf("expected count to reset to 0, got %d", manager.Count())
	}
}
