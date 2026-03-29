package domain

import (
	"time"

	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type Person struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"not null"`
	Age              int            `json:"age" gorm:"not null"`
	Motto            string         `json:"motto"`
	GenderIdentity   string         `json:"gender_identity"`
	InterestedIn     string         `json:"interested_in"`
	RelationshipGoal string         `json:"relationship_goal"`
	LatLocation      float64        `json:"lat" gorm:"not null"`
	LongLocation     float64        `json:"long" gorm:"not null"`
	Profile          ProfilePhoto   `json:"profile" gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Details          Details        `json:"details" gorm:"embedded"`
	Photos           []ProfilePhoto `json:"photos" gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description      string         `json:"description"`
	CreatedAt        time.Time      `json:"create_time"`
	UpdatedAt        time.Time      `json:"update_time"`
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Person       Person    `json:"user" gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"not null"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	LastLogin    time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"create_time"`
	UpdatedAt    time.Time `json:"update_time"`
}

type UserSession struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	JTI       string     `json:"jti" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
}

type Config struct {
	Postgres struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		SSLMode  string `json:"sslmode,omitempty"`
	} `json:"postgres"`
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
	ID        uint      `json:"id" gorm:"primaryKey"`
	PersonID  uint      `json:"person_id" gorm:"not null"`
	Url       string    `json:"url" gorm:"-"`
	Upload    string    `json:"upload,omitempty" gorm:"-"`
	Delete    string    `json:"delete,omitempty" gorm:"-"`
	S3Key     string    `json:"s3key" gorm:"not null;index"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"create_time"`
	UpdatedAt time.Time `json:"update_time"`
}

type ChatMessage struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement" bson:"id"`
	MatchID     int        `json:"match_id" gorm:"not null;index" bson:"match_id"`
	Time        time.Time  `json:"time" bson:"time"`
	Who         uint       `json:"who" gorm:"not null" bson:"who"`
	MessageType string     `json:"message_type" gorm:"not null;default:user" bson:"message_type"`
	Message     string     `json:"message" gorm:"not null" bson:"message"`
	ReadAt      *time.Time `json:"read_at,omitempty" bson:"read_at,omitempty" gorm:"-"`
	CreatedAt   time.Time  `json:"create_time"`
	UpdatedAt   time.Time  `json:"update_time"`
}

type Conversation struct {
	MatchID  uint          `json:"match_id" bson:"match_id"`
	Messages []ChatMessage `json:"messages" bson:"messages"`
}

type MatchReadState struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	MatchID           uint      `json:"match_id" gorm:"not null;index:idx_match_user_read_state,unique"`
	UserID            uint      `json:"user_id" gorm:"not null;index:idx_match_user_read_state,unique"`
	LastReadMessageID int64     `json:"last_read_message_id" gorm:"not null;default:0"`
	CreatedAt         time.Time `json:"create_time"`
	UpdatedAt         time.Time `json:"update_time"`
}

type Details struct {
	Bouginess            int `json:"bouginess,omitempty"`
	Cats                 int `json:"cats,omitempty"`
	Dogs                 int `json:"dogs,omitempty"`
	Drinking             int `json:"drinking,omitempty"`
	EnergyLevel          int `json:"energy_levels,omitempty"`
	Food                 int `json:"food,omitempty"`
	ImportanceOfPolitics int `json:"importance_of_politics,omitempty"`
	Kids                 int `json:"kids,omitempty"`
	Outdoorsyness        int `json:"outdoorsy_ness,omitempty"`
	Religion             int `json:"religion,omitempty"`
	Smoking              int `json:"smoking,omitempty"`
	Travel               int `json:"travel,omitempty"`
}

type Match struct {
	gorm.Model
	Offered         uint      `json:"offered"`
	OfferedProfile  Person    `gorm:"foreignKey:Offered;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OfferedTime     time.Time `json:"offered_time"`
	OfferedLiked    int       `json:"offered_liked"`
	Accepted        uint      `json:"accepted"`
	AcceptedProfile Person    `gorm:"foreignKey:Accepted;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AcceptedTime    time.Time `json:"accepted_time"`
	VibeChat        bool      `json:"vibe_chat"`
	UnreadOffered   int       `json:"unread_offered" gorm:"default:0"`
	UnreadAccepted  int       `json:"unread_accepted" gorm:"default:0"`
	LastMessageTime time.Time `json:"last_message_time"`
}
