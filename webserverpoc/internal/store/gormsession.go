package store

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GormSessionStore struct {
	db *gorm.DB
}

func NewGormSessionStore(db *gorm.DB) *GormSessionStore {
	return &GormSessionStore{db: db}
}

type sessionRecord struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	JTI       string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (sessionRecord) TableName() string {
	return "user_sessions"
}

func (s *GormSessionStore) Create(ctx context.Context, session Session) error {
	record := sessionRecord{
		UserID:    session.UserID,
		JTI:       session.JTI,
		ExpiresAt: session.ExpiresAt,
		RevokedAt: session.RevokedAt,
	}
	return s.db.WithContext(ctx).Create(&record).Error
}

func (s *GormSessionStore) GetByJTI(ctx context.Context, jti string) (Session, error) {
	var record sessionRecord
	if err := s.db.WithContext(ctx).Where("jti = ?", jti).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Session{}, ErrNotFound
		}
		return Session{}, err
	}
	return Session{
		UserID:    record.UserID,
		JTI:       record.JTI,
		ExpiresAt: record.ExpiresAt,
		RevokedAt: record.RevokedAt,
	}, nil
}

func (s *GormSessionStore) RevokeByJTI(ctx context.Context, jti string, revokedAt time.Time) error {
	result := s.db.WithContext(ctx).Model(&sessionRecord{}).
		Where("jti = ?", jti).
		Update("revoked_at", revokedAt)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
