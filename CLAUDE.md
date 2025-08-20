# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This project is a dating web application with a Go backend and React frontend. The application allows users to:
- Sign up, log in, and manage profiles
- Browse and search for potential matches based on location and preferences
- Confirm/reject potential matches
- Chat with matched users in real-time using WebSockets
- Receive real-time notifications for new messages using Server-Sent Events (SSE)
- Upload and manage profile photos using AWS S3

## Backend Architecture

The backend is a Go web server built with:
- [Gin](https://github.com/gin-gonic/gin) for the HTTP router and middleware
- [GORM](https://gorm.io/) for MySQL database interactions
- [MongoDB](https://go.mongodb.org/mongo-driver) for chat message storage
- [JWT](https://github.com/golang-jwt/jwt) for authentication
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for real-time chat
- Server-Sent Events (SSE) for real-time notifications

### Key Backend Components:

1. **Database Layer**: 
   - MySQL for user accounts, profiles, matches
   - MongoDB for chat message history

2. **Authentication**:
   - JWT-based authentication with token validation in middleware

3. **API Endpoints**:
   - User authentication (/signup, /login, /logout)
   - Profile management (/people, /people/:id)
   - Match handling (/matches, /matches/:id)
   - WebSocket connections (/ws)
   - SSE notifications (/notifications/:id)
   - Chat management (/chat/markread/:id)

4. **File Storage**:
   - AWS S3 for profile and photo storage
   - Presigned URLs for secure direct uploads

## Frontend Architecture

The frontend is a React application built with Vite featuring:
- Bootstrap for styling
- React hooks for state management
- WebSockets for real-time chat functionality
- Server-Sent Events for real-time notifications
- MapBox for location-based features
- Custom components for profile editing, match management, and chat

### Key Frontend Components:
- Authentication flow (login/signup)
- Profile creation and editing
- Match browsing and selection
- Real-time chat interface with unread message indicators
- Real-time notification system for new messages
- Advanced search with preference matching

## Development Commands

### Backend (Go)

```bash
# Change to the backend directory
cd webserverpoc

# Run the backend server
go run .

# Build the backend
go build -o server .

# Run tests
go test -v ./...

# Run tests with coverage
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Frontend (React)

```bash
# Change to the frontend directory
cd react-vite-poc

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Run linter
npm run lint

# Preview production build
npm run preview

# Run tests
npm test

# Run tests with coverage
npm test -- --coverage
```

## Configuration

The application requires configuration files for database connections:

- `configlocal.json` - Local development configuration
- `config.json` - Production/Docker configuration

Example configuration format:
```json
{
  "mysql": {
    "host": "localhost",
    "port": 3306,
    "user": "username",
    "password": "password",
    "database": "gowebserver"
  },
  "mongodb": {
    "host": "localhost",
    "port": 27017,
    "user": "username",
    "password": "password",
    "database": "urmid"
  }
}
```

## Environment Variables

The application requires the following environment variables:
- `AWS_REGION` - AWS region for S3 bucket
- `AWS_S3_BUCKET` - S3 bucket name for photo storage

## Project Structure

```
gowebserver/
├── webserverpoc/         # Go backend
│   ├── main.go           # Entry point
│   ├── handlers.go       # HTTP handlers
│   ├── models.go         # Data models
│   ├── database.go       # Database connections
│   ├── websocket.go      # WebSocket implementation and SSE notifications
│   ├── middleware.go     # Auth middleware
│   ├── artificialintelligence.go # AI integration
│   ├── middleware_test.go # Authentication tests
│   ├── websocket_test.go # WebSocket/SSE tests 
│   ├── handlers_test.go  # API handler tests
│   └── testutils/        # Testing utilities
├── react-vite-poc/       # React frontend
│   ├── src/              # Source code
│   │   ├── App.jsx       # Main application component
│   │   ├── components/   # UI components
│   │   │   └── __tests__/ # Component tests
│   │   └── __mocks__/    # Test mocks
│   ├── public/           # Static assets
│   ├── jest.config.js    # Jest configuration
│   ├── FRONTEND_TEST_SETUP.md # Frontend test documentation
│   └── package.json      # Dependencies
├── webserverpoc/BACKEND_TEST_SETUP.md # Backend test documentation
├── TODOTesting.md        # Testing strategy documentation
└── go-ai/                # AI integration module
    └── main.go           # AI service
```

## Testing

The application has comprehensive testing setup for both backend and frontend:

### Backend Testing

- Unit tests for authentication, handlers, and WebSocket/SSE
- Mock implementations for database and external services
- Test utilities for common testing needs
- See `BACKEND_TEST_SETUP.md` for detailed instructions

### Frontend Testing

- React Testing Library for component testing
- Jest for test running and assertions
- Mock implementations for WebSockets and Server-Sent Events
- See `FRONTEND_TEST_SETUP.md` for detailed instructions

### Testing Strategy

The project follows a comprehensive testing strategy outlined in `TODOTesting.md` including:
- Test coverage targets for different components
- Mocking strategies for external dependencies
- Test prioritization matrix
- Implementation tracking elements