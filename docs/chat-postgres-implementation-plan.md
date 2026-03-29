# Chat Postgres Implementation Plan

## Summary
Chat should move to a single durable Postgres-backed design. The current in-memory plus Mongo arrangement has failed in the ways that matter most:

- message history is not reliably durable
- reconnect behavior is inconsistent
- websocket state is doing too much work
- the server cannot honestly be treated as stateless

The target design is intentionally boring: write messages to Postgres first, then fan them out over websocket connections if recipients are currently connected.

## Target Data Model
Use Postgres as the source of truth for chat history and read state.

### `chat_messages`
- `id`
- `match_id`
- `sender_id`
- `message_type`
- `body`
- `created_at`
- `updated_at`

Recommended message types:
- `user`
- `system_intro`
- `vibe_chat`
- `date_suggestion`

System/AI messages should not increment unread counts by default.

### `match_read_states`
- `id`
- `match_id`
- `user_id`
- `last_read_message_id`
- `updated_at`

For this app, `last_read_message_id` is a better unread mechanism than updating `read_at` on every message row. It keeps "mark thread read" cheap and makes unread calculations predictable.

## Runtime Shape
The websocket layer should only do:

1. authenticate and validate the user against the match
2. persist the incoming message
3. update unread/read state
4. broadcast to active connections

Process memory should keep only live connection maps. It should not be authoritative for:

- pending messages
- message history
- unread state

## Implementation Order
1. Add relational chat store and read-state store in `internal/store`
2. Expand the chat domain/store models to include message type and read-state data
3. Update `ChatService` to use Postgres-backed history and relational read tracking
4. Simplify `WebsocketListener` so it loads history from Postgres and writes through to Postgres on every message
5. Remove in-memory message buffering and Mongo chat persistence paths
6. Keep unread counters on `matches` temporarily if helpful, but derive their updates from relational chat events only

## Transitional Decisions
- Keep websocket fanout in memory for now
- Keep existing unread counters on `matches` during the first cut to avoid changing too much at once
- Retire Mongo from chat persistence immediately once the Postgres store is in place

## Expected Outcome
This change will make chat durable, reconnect-safe, and much closer to stateless server behavior. It also clears the path for future scaling because application memory will no longer be part of the durable chat model.
