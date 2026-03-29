package store

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormChatStore struct {
	db *gorm.DB
}

func NewGormChatStore(db *gorm.DB) *GormChatStore {
	return &GormChatStore{db: db}
}

type chatMessageRecord struct {
	ID          int64 `gorm:"primaryKey"`
	MatchID     int
	Time        time.Time
	Who         uint
	MessageType string
	Message     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (chatMessageRecord) TableName() string {
	return "chat_messages"
}

type matchReadStateRecord struct {
	ID                uint `gorm:"primaryKey"`
	MatchID           uint
	UserID            uint
	LastReadMessageID int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (matchReadStateRecord) TableName() string {
	return "match_read_states"
}

func (s *GormChatStore) LoadByMatchID(ctx context.Context, matchID uint) (Conversation, error) {
	var records []chatMessageRecord
	if err := s.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Order("time asc").
		Order("id asc").
		Find(&records).Error; err != nil {
		return Conversation{}, err
	}
	if len(records) == 0 {
		return Conversation{MatchID: matchID, Messages: []ChatMessage{}}, ErrNotFound
	}

	conversation := Conversation{
		MatchID:  matchID,
		Messages: make([]ChatMessage, 0, len(records)),
	}
	for _, record := range records {
		conversation.Messages = append(conversation.Messages, ChatMessage{
			ID:          record.ID,
			MatchID:     record.MatchID,
			Time:        record.Time,
			Who:         record.Who,
			MessageType: record.MessageType,
			Message:     record.Message,
			CreatedAt:   record.CreatedAt,
			UpdatedAt:   record.UpdatedAt,
		})
	}
	return conversation, nil
}

func (s *GormChatStore) SaveConversation(ctx context.Context, conversation Conversation) error {
	for _, message := range conversation.Messages {
		if _, err := s.appendMessageRecord(ctx, message); err != nil {
			return err
		}
	}
	return nil
}

func (s *GormChatStore) AppendMessage(ctx context.Context, matchID uint, message ChatMessage) error {
	message.MatchID = int(matchID)
	_, err := s.appendMessageRecord(ctx, message)
	return err
}

func (s *GormChatStore) LatestMessageID(ctx context.Context, matchID uint) (int64, error) {
	var record chatMessageRecord
	if err := s.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Order("id desc").
		First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return record.ID, nil
}

func (s *GormChatStore) SaveReadState(ctx context.Context, state MatchReadState) error {
	record := matchReadStateRecord{
		MatchID:           state.MatchID,
		UserID:            state.UserID,
		LastReadMessageID: state.LastReadMessageID,
		UpdatedAt:         state.UpdatedAt,
	}

	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "match_id"},
			{Name: "user_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"last_read_message_id": record.LastReadMessageID,
			"updated_at":           record.UpdatedAt,
		}),
	}).Create(&record).Error
}

func (s *GormChatStore) appendMessageRecord(ctx context.Context, message ChatMessage) (ChatMessage, error) {
	record := chatMessageRecord{
		ID:          message.ID,
		MatchID:     message.MatchID,
		Time:        message.Time,
		Who:         message.Who,
		MessageType: message.MessageType,
		Message:     message.Message,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	}

	if record.MessageType == "" {
		record.MessageType = "user"
	}
	if record.Time.IsZero() {
		record.Time = time.Now()
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = record.Time
	}
	if record.UpdatedAt.IsZero() {
		record.UpdatedAt = record.Time
	}

	if record.ID == 0 {
		if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
			return ChatMessage{}, err
		}
	} else {
		if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&record).Error; err != nil {
			return ChatMessage{}, err
		}
	}

	return ChatMessage{
		ID:          record.ID,
		MatchID:     record.MatchID,
		Time:        record.Time,
		Who:         record.Who,
		MessageType: record.MessageType,
		Message:     record.Message,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}, nil
}
