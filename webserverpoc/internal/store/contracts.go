package store

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("store: not found")

type User struct {
	ID      uint
	Email   string
	IsAdmin bool
}

type Session struct {
	UserID    uint
	JTI       string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type Match struct {
	ID              uint
	Offered         uint
	Accepted        uint
	UnreadOffered   int
	UnreadAccepted  int
	LastMessageTime time.Time
}

type ChatMessage struct {
	ID          int64     `bson:"id"`
	MatchID     int       `bson:"match_id"`
	Time        time.Time `bson:"time"`
	Who         uint      `bson:"who"`
	MessageType string    `bson:"message_type"`
	Message     string    `bson:"message"`
	CreatedAt   time.Time `bson:"create_time"`
	UpdatedAt   time.Time `bson:"update_time"`
}

type Conversation struct {
	MatchID  uint          `bson:"match_id"`
	Messages []ChatMessage `bson:"messages"`
}

type MatchReadState struct {
	MatchID           uint
	UserID            uint
	LastReadMessageID int64
	UpdatedAt         time.Time
}

type UserStore interface {
	GetByID(ctx context.Context, id uint) (User, error)
}

type SessionStore interface {
	Create(ctx context.Context, session Session) error
	GetByJTI(ctx context.Context, jti string) (Session, error)
	RevokeByJTI(ctx context.Context, jti string, revokedAt time.Time) error
}

type MatchStore interface {
	GetByID(ctx context.Context, id uint) (Match, error)
	Update(ctx context.Context, match Match) error
}

type ChatHistoryStore interface {
	LoadByMatchID(ctx context.Context, matchID uint) (Conversation, error)
	SaveConversation(ctx context.Context, conversation Conversation) error
	AppendMessage(ctx context.Context, matchID uint, message ChatMessage) error
	LatestMessageID(ctx context.Context, matchID uint) (int64, error)
	SaveReadState(ctx context.Context, state MatchReadState) error
}
