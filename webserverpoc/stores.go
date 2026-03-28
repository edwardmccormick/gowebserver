package gowebserver

import (
	"errors"

	"github.com/edwardmccormick/gowebserver/internal/store"
)

var (
	userStore        store.UserStore
	matchStore       store.MatchStore
	chatHistoryStore store.ChatHistoryStore
)

func initializeStores() error {
	if db == nil {
		return errors.New("sql database is not initialized")
	}

	userStore = store.NewGormUserStore(db)
	matchStore = store.NewGormMatchStore(db)

	if mongoClient != nil {
		chatHistoryStore = store.NewMongoChatStore(mongoClient, "urmid")
	}

	return nil
}
