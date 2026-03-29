package gowebserver

import (
	"errors"

	"github.com/edwardmccormick/gowebserver/internal/store"
)

var (
	userStore        store.UserStore
	sessionStore     store.SessionStore
	matchStore       store.MatchStore
	chatHistoryStore store.ChatHistoryStore
)

func initializeStores() error {
	if db == nil {
		return errors.New("sql database is not initialized")
	}

	userStore = store.NewGormUserStore(db)
	sessionStore = store.NewGormSessionStore(db)
	matchStore = store.NewGormMatchStore(db)
	chatHistoryStore = store.NewGormChatStore(db)

	return nil
}
