package store

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GormUserStore struct {
	db *gorm.DB
}

type GormMatchStore struct {
	db *gorm.DB
}

func NewGormUserStore(db *gorm.DB) *GormUserStore {
	return &GormUserStore{db: db}
}

func NewGormMatchStore(db *gorm.DB) *GormMatchStore {
	return &GormMatchStore{db: db}
}

type userRecord struct {
	ID      uint
	Email   string
	IsAdmin bool
}

func (userRecord) TableName() string {
	return "users"
}

type matchRecord struct {
	ID              uint `gorm:"primaryKey"`
	Offered         uint
	Accepted        uint
	UnreadOffered   int
	UnreadAccepted  int
	LastMessageTime time.Time
}

func (matchRecord) TableName() string {
	return "matches"
}

func (s *GormUserStore) GetByID(ctx context.Context, id uint) (User, error) {
	var record userRecord
	if err := s.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return User{
		ID:      record.ID,
		Email:   record.Email,
		IsAdmin: record.IsAdmin,
	}, nil
}

func (s *GormMatchStore) GetByID(ctx context.Context, id uint) (Match, error) {
	var record matchRecord
	if err := s.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Match{}, ErrNotFound
		}
		return Match{}, err
	}

	return Match{
		ID:              record.ID,
		Offered:         record.Offered,
		Accepted:        record.Accepted,
		UnreadOffered:   record.UnreadOffered,
		UnreadAccepted:  record.UnreadAccepted,
		LastMessageTime: record.LastMessageTime,
	}, nil
}

func (s *GormMatchStore) Update(ctx context.Context, match Match) error {
	record := matchRecord{
		ID:              match.ID,
		Offered:         match.Offered,
		Accepted:        match.Accepted,
		UnreadOffered:   match.UnreadOffered,
		UnreadAccepted:  match.UnreadAccepted,
		LastMessageTime: match.LastMessageTime,
	}
	return s.db.WithContext(ctx).Save(&record).Error
}
