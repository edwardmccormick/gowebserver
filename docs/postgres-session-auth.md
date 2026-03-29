# Postgres Session and Token Revocation Plan

## Summary
JWT expiration should remain in the token itself via the `exp` claim. Postgres should be used for session persistence and revocation, not as a replacement for token expiry semantics.

This gives the app a durable, multi-session authentication model without relying on in-memory token blacklists or Redis.

## Recommended Model
Use short-lived JWTs that include a `jti` claim. Persist one session row per successful login/signup.

### `user_sessions`
- `id`
- `user_id`
- `jti`
- `expires_at`
- `revoked_at`
- `created_at`
- `updated_at`

The JWT should include:
- `sub`
- `exp`
- `email`
- `jti`

## Request Validation
On each authenticated request:

1. validate JWT signature and `exp`
2. extract `jti`
3. load the matching session from Postgres
4. reject if the session is missing, expired, or revoked

## Logout Behavior
Logout should revoke the session identified by the token's `jti`, not just store the raw token string in memory.

This supports:
- multiple concurrent logins for the same user
- revoking one session without invalidating every device
- future "log out all sessions" behavior
- stateless application instances with durable auth state

## Why This Fits This Repo
The frontend intentionally does not use `localStorage`, and multiple simultaneous account logins are part of the current workflow. A Postgres-backed session model supports that well:

- each login creates a distinct session
- each token maps to one row
- revocation is durable across restarts
- no cache layer is required for current scale

## Implementation Direction
1. add `user_sessions` to the schema
2. create sessions on signup/login
3. include `jti` in signed JWTs
4. make middleware validate session state from Postgres
5. make logout revoke the session row
6. keep in-memory token revocation only as a test/local fallback during migration
