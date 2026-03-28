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

type Match struct {
	ID              uint
	Offered         uint
	Accepted        uint
	UnreadOffered   int
	UnreadAccepted  int
	LastMessageTime time.Time
}

type ChatMessage struct {
	ID        int64     `bson:"id"`
	MatchID   int       `bson:"match_id"`
	Time      time.Time `bson:"time"`
	Who       uint      `bson:"who"`
	Message   string    `bson:"message"`
	CreatedAt time.Time `bson:"create_time"`
	UpdatedAt time.Time `bson:"update_time"`
}

type Conversation struct {
	MatchID  uint          `bson:"match_id"`
	Messages []ChatMessage `bson:"messages"`
}

type UserStore interface {
	GetByID(ctx context.Context, id uint) (User, error)
}

type MatchStore interface {
	GetByID(ctx context.Context, id uint) (Match, error)
	Update(ctx context.Context, match Match) error
}

type ChatHistoryStore interface {
	LoadByMatchID(ctx context.Context, matchID uint) (Conversation, error)
	SaveConversation(ctx context.Context, conversation Conversation) error
	AppendMessage(ctx context.Context, matchID uint, message ChatMessage) error
}
