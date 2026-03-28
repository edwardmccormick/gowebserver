package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"file:%s?mode=memory&cache=shared",
		strings.NewReplacer("/", "_", "\\", "_", " ", "_").Replace(t.Name()),
	)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(
		&domain.Person{},
		&domain.ProfilePhoto{},
		&domain.User{},
		&domain.Match{},
	); err != nil {
		t.Fatalf("auto migrate sqlite db: %v", err)
	}

	return db
}
