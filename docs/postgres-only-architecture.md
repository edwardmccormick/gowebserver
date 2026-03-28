# Postgres-Only Architecture Direction

## Summary
For this codebase, a Postgres-only design is the cleanest next target. It simplifies the current `MySQL + Mongo + in-memory websocket state` split, removes cross-store consistency problems, and gives enough headroom to delay Redis until realtime scale or operational complexity actually justifies it.

The current product model is mostly relational already: users, profiles, matches, unread counts, admin views, moderation, and reporting all fit naturally in SQL. Chat also fits well enough in Postgres that Mongo is not buying much here.

## Why Postgres Fits Better
- Strong relational support for the existing core domain.
- JSONB support for AI/system-message metadata, prompt payloads, or flexible moderation data.
- Better long-term query flexibility for admin, safety, and analytics work.
- Easier consistency than splitting durable state across MySQL and Mongo.
- Cleaner local/dev/serverless story with one durable backing store.

## Recommended Data Shape
Keep Postgres as the system of record for all durable product data:

- `users`
- `profiles` or `people`
- `profile_photos`
- `matches`
- `chat_messages`
- `message_reads` or per-match read markers
- `moderation_events`
- optional `audit_events`

For chat, prefer a relational design:

- `chat_messages(id, match_id, sender_id, body, message_type, metadata_jsonb, created_at)`
- `match_participants` only if the model expands beyond 1:1 matching
- `match_read_state(match_id, user_id, last_read_message_id, updated_at)` for unread counts

This is simpler than storing history in Mongo while also keeping unread state in SQL.

## What This Replaces
Postgres-only should replace:

- MySQL as the primary relational store
- Mongo for chat history persistence

It does **not** fully replace:

- in-memory websocket connection state during a single process lifetime
- eventual distributed coordination concerns if the app scales horizontally

That means Redis can stay out of scope for now, but still be added later for presence, fanout, token revocation, or rate limiting if needed.

## Migration Direction
1. Refactor `internal/store` behind interfaces first.
2. Model chat as relational storage, not a document store.
3. Implement a Postgres-backed store layer.
4. Migrate MySQL schema and seed data.
5. Retire Mongo-specific chat code and simplify websocket persistence flows.

## Immediate Architecture Implications
The next refactor should assume:

- `internal/store/postgres` becomes the primary target
- chat history should move out of Mongo-specific code paths
- websocket/SSE logic should depend on store interfaces, not database clients
- unread-count logic should be tied to relational read-state, not mixed in-memory/document behavior

## Recommendation
Proceed as if Postgres is the destination architecture. Keep Redis as a future operational optimization, not a present dependency. Retire Mongo from the planned target state.
