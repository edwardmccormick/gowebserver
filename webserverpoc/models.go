package main

import (
	"time"

	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type Person struct {
	ID           uint      `json:"id" db:"id" gorm:"primaryKey"`
	Name         string    `json:"name" db:"name" gorm:"type:varchar(255) not null"`
	Age          int       `json:"age" db:"age" gorm:"type:int not null"`
	Motto        string    `json:"motto" db:"motto" gorm:"type:varchar(255)"`
	LatLocation  float64   `json:"lat" db:"lat" gorm:"type:float not null"`
	LongLocation float64   `json:"long" db:"long" gorm:"type:float not null"`
	Profile      string    `json:"profile" db:"profile" gorm:"type:varchar(255)"`
	Details      Details   `json:"details" db:"details" gorm:"embedded"`
	Description  string    `json:"description" db:"description" gorm:"type:text"`
	CreatedAt    time.Time `json:"create_time" db:"create_time" gorm:"<-:create"`
	UpdatedAt    time.Time `json:"update_time" db:"update_time" gorm:"<-:update"`
}

type User struct {
	ID           uint      `json:"id" db:"id" gorm:"uniqueIndex;not null"`
	Person       Person    `json:"user" db:"user" gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email        string    `json:"email" db:"email" gorm:"type:varchar(255);not null"`
	PasswordHash string    `json:"-" db:"password_hash" gorm:"type:varchar(255);not null"`
	LastLogin    time.Time `json:"last_login" db:"last_login" gorm:"<-:update"`
	CreatedAt    time.Time `json:"create_time" db:"create_time" gorm:"<-:create"`
	UpdatedAt    time.Time `json:"update_time" db:"update_time" gorm:"<-:update"`
}

type Config struct {
	MySQL struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"mysql"`
	Mongo struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"mongodb"`
}

type ProfilePhoto struct {
	Url     string `json:"url"`
	Caption string `json:"caption"`
}

type ProfileAlbum struct {
	ID     uint           `json:"id" db:"id" gorm:"uniqueIndex;not null"`
	Person Person         `json:"-" db:"user" gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Photos []ProfilePhoto `json:"photos"` // Using JSONB for PostgreSQL or JSON for MySQL
}

type ChatMessage struct {
	ID      int64 `json:"id" gorm:"primaryKey,AutoIncrement not null,Unique" bson:"id"`
	MatchID int   `json:"match_id" gorm:"not null" bson:"match_id"` // Foreign key to Match.MatchID
	// Match      Match     `gorm:"foreignKey:MatchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key relationship
	Time time.Time `json:"time" bson:"time"`
	Who  uint      `json:"who" gorm:"not null" bson:"who"` // Foreign key to Person.ID
	// WhoProfile Person    `gorm:"foreignKey:Who;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key relationship
	Message string `json:"message" gorm:"type:text;not null" bson:"message"`
}

type Conversation struct {
	MatchID  uint          `json:"match_id" bson:"match_id"` //
	Messages []ChatMessage `json:"messages" bson:"messages"` // Array of chat messages
}

type Details struct {
	Bouginess               int `json:"bouginess"`
	Cats                    int `json:"cats"`
	Dogs                    int `json:"dogs"`
	Drinking                int `json:"drinking"`
	EnergyLevel             int `json:"energy_levels"`
	Food                    int `json:"food"`
	ImportantanceOfPolitics int `json:"importantance_of_politics"`
	Kids                    int `json:"kids"`
	Outdoorsyness           int `json:"outdoorsy_ness"`
	Religion                int `json:"religion"`
	Smoking                 int `json:"smoking"`
	Travel                  int `json:"travel"`
}

type Match struct {
	gorm.Model
	Offered         uint      `json:"offered"`                                                           // Foreign key to User.ID
	OfferedProfile  Person    `gorm:"foreignKey:Offered;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Foreign key relationship
	OfferedTime     time.Time `json:"offered_time"`
	OfferedLiked    int       `json:"offered_liked"`
	Accepted        uint      `json:"accepted"`                                                           // Foreign key to User.ID
	AcceptedProfile Person    `gorm:"foreignKey:Accepted;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Foreign key relationship
	AcceptedTime    time.Time `json:"accepted_time"`
	VibeChat        bool      `json:"vibe_chat"`

	// OfferedChat  []ChatMessage `json:"offered_chat"`
	// AcceptedChat []ChatMessage `json:"accepted_chat"`
	// ChatHistory  []ChatMessage `json:"chat_history"`
}
