# Codebase Assessment and Implementation Path

## Executive Summary
The repository is no longer just a concept POC. The core dating flow exists across the Go API and the Vite frontend: account signup/login, profile creation/editing, photo upload via S3 presigned URLs, matching, chat history, unread counts, SSE notifications, and Gemini-powered "vibe chat" and date suggestions are all implemented in some form.

The main issue is not missing surface area. It is that several core paths are still prototype-grade and need hardening before additional feature expansion.

## Highest-Risk Findings
1. Core write routes are not consistently protected by auth middleware. `POST /people`, `POST /matches`, and several read endpoints are exposed at the router level in `webserverpoc/main.go`, while handlers like `PostPeople` trust caller-supplied IDs. This is the first thing to fix before building more social features.
2. There is effectively no automated test coverage. The repo has implementation breadth but little regression protection, especially around auth, matching, unread counts, and chat persistence.
3. The mobile app is still a stub. `react-native-expo/urmid/App.tsx` is not aligned with the web app feature set, so "React scales to mobile" is still only partially answered.
4. Product logic is still thin in several places. Matching defaults, stats, moderation, password recovery, and trust-and-safety controls remain mostly unimplemented.

## What Appears Complete or Close
- Signup/login with JWT issuance
- Profile create/update flow
- Photo management and profile photo cropping
- Match creation and match review UI
- WebSocket chat with Mongo-backed history
- Unread counters plus SSE notifications
- Mapbox match visualization
- Gemini-powered introduction, vibe chat, and date suggestions
- Docker Compose for app + MySQL + MongoDB

## TODO Status Readout
Implemented or partially implemented:
- Frontend signup/profile split, profile updates, photo upload/crop
- Backend persistence, chat history, AI features
- Architectural containerization and initial Expo exploration

Still materially open:
- OAuth sign-in
- Password reset/recovery
- Likes on prompts/photos with message context
- Better match ranking defaults
- User stats
- Geocoding and location-sharing UX
- Moderation and anti-scam controls
- Legal/trust-and-safety documents

## Recommended Path Forward
### Phase 1: Hardening
- Put JWT auth on all profile, match, chat-history, and admin-sensitive routes.
- Derive acting user IDs from the token, not request payloads.
- Add table-driven Go tests for auth, profile writes, match creation, unread count updates, and mark-read flows.
- Move secrets and local credentials out of tracked config where possible.

### Phase 2: Product Completion
- Implement password reset first; it unlocks a baseline production account lifecycle.
- Add explicit "like with context" data structures for photo/prompt likes instead of overloading `Match`.
- Formalize match ranking inputs and introduce a simple stats model.

### Phase 3: Safety and Scale
- Add first-message and link/contact moderation in chat.
- Introduce audit logging and rate limiting for signup, login, matching, and chat endpoints.
- Decide whether SSE remains acceptable or whether notifications should move to a queue/push-friendly design.

### Phase 4: Platform Expansion
- Either port the working web flows into Expo or pause mobile until backend/auth/testing are stable.
- Rework the frontend visual system only after the data contracts and safety controls stop moving.

## Suggested Next Implementation Slice
If only one slice should be tackled next, make it: authenticated profile and match writes plus a small Go test suite around those flows. That reduces the largest security risk and creates a base for every remaining TODO item.
