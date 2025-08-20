# Go Dating Web Application

A dating web application with a Go backend and React frontend that allows users to find matches, chat in real-time, and get notifications for new messages.

## Recent Updates

### Testing Infrastructure

We've implemented comprehensive testing infrastructure for both frontend and backend:

- **Backend Testing**:
  - Added unit tests for authentication middleware
  - Implemented WebSocket and SSE testing utilities
  - Created handler unit tests for key endpoints
  - Added database mocking utilities for easier test setup
  - Developed structured test files for core components

- **Frontend Testing**:
  - Added React Testing Library setup
  - Created component tests for NotificationService and ChatSelect
  - Implemented mocking utilities for WebSockets and SSE
  - Added Jest configuration for the React application

- **Documentation**:
  - Created comprehensive testing documentation in TODOTesting.md
  - Added detailed backend testing setup guide in BACKEND_TEST_SETUP.md
  - Added frontend testing setup guide in FRONTEND_TEST_SETUP.md
  - Developed test implementation tracking with coverage targets

### Real-time Notification System

We've implemented a comprehensive real-time notification system using Server-Sent Events (SSE):

- **Backend Implementation**:
  - Added SSE endpoint for push notifications (`/notifications/:id`)
  - Created notification tracking in the Match model (UnreadOffered, UnreadAccepted fields)
  - Implemented in-memory message sharing for users joining active chats
  - Added message read status tracking and unread count management

- **Frontend Implementation**:
  - Added NotificationService component to establish SSE connections
  - Implemented notification badges for unread messages
  - Added individual chat entry unread message indicators
  - Implemented automatic notification clearing when opening chats

### Chat Improvements

- Added "Perfect Date" feature using AI to generate personalized date suggestions
- Fixed message duplication issues when reopening chat windows
- Added VibeChat feature to generate conversation starters
- Implemented chat history persistence with MongoDB

### Authentication and Security

- Improved JWT authentication with support for both header and query parameter tokens
- Implemented thread-safe access to shared resources

## Project Structure

The application follows a client-server architecture:

- **Backend**: Go-based web server with Gin, GORM, MongoDB, and WebSockets
- **Frontend**: React application with Bootstrap, React hooks, and real-time features

## Getting Started

See the CLAUDE.md file for detailed information about:
- Project overview and architecture
- Development commands
- Configuration requirements
- Environment variables
- Full project structure

## Testing

For comprehensive testing documentation, see:

- **[TODOTesting.md](TODOTesting.md)** - Overall testing strategy and priorities
- **[BACKEND_TEST_SETUP.md](webserverpoc/BACKEND_TEST_SETUP.md)** - Go backend testing setup
- **[FRONTEND_TEST_SETUP.md](react-vite-poc/FRONTEND_TEST_SETUP.md)** - React frontend testing setup

### Running Backend Tests

```bash
cd webserverpoc
go test -v ./...
```

### Running Frontend Tests

```bash
cd react-vite-poc
npm test
```