package domain

import (
	"time"

	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type Person struct {
	ID           uint           `json:"id" db:"id" gorm:"primaryKey"`
	Name         string         `json:"name" db:"name" gorm:"type:varchar(255) not null"`
	Age          int            `json:"age" db:"age" gorm:"type:int not null"`
	Motto        string         `json:"motto" db:"motto" gorm:"type:varchar(255)"`
	LatLocation  float64        `json:"lat" db:"lat" gorm:"type:float not null"`
	LongLocation float64        `json:"long" db:"long" gorm:"type:float not null"`
	Profile      ProfilePhoto   `json:"profile" gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Details      Details        `json:"details" db:"details" gorm:"embedded"`
	Photos       []ProfilePhoto `json:"photos" gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description  string         `json:"description" db:"description" gorm:"type:text"`
	CreatedAt    time.Time      `json:"create_time" db:"create_time" gorm:"<-:create"`
	UpdatedAt    time.Time      `json:"update_time" db:"update_time" gorm:"<-:update"`
}

type User struct {
	ID           uint      `json:"id" db:"id" gorm:"uniqueIndex;not null"`
	Person       Person    `json:"user" db:"user" gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email        string    `json:"email" db:"email" gorm:"type:varchar(255);not null"`
	PasswordHash string    `json:"-" db:"password_hash" gorm:"type:varchar(255);not null"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin" gorm:"default:false"`
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
	ID        uint      `json:"id" gorm:"primaryKey"`
	PersonID  uint      `json:"person_id" gorm:"not null"`
	Url       string    `json:"url" gorm:"-"`
	Upload    string    `json:"upload,omitempty" gorm:"-"`
	Delete    string    `json:"delete,omitempty" gorm:"-"`
	S3Key     string    `json:"s3key" gorm:"not null,unique"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"create_time" db:"create_time" gorm:"<-:create"`
	UpdatedAt time.Time `json:"update_time" db:"update_time" gorm:"<-:update"`
}

type ChatMessage struct {
	ID        int64     `json:"id" gorm:"primaryKey,AutoIncrement not null,Unique" bson:"id"`
	MatchID   int       `json:"match_id" gorm:"not null" bson:"match_id"`
	Time      time.Time `json:"time" bson:"time"`
	Who       uint      `json:"who" gorm:"not null" bson:"who"`
	Message   string    `json:"message" gorm:"type:text;not null" bson:"message"`
	CreatedAt time.Time `json:"create_time" db:"create_time" gorm:"<-:create"`
	UpdatedAt time.Time `json:"update_time" db:"update_time" gorm:"<-:update"`
}

type Conversation struct {
	MatchID  uint          `json:"match_id" bson:"match_id"`
	Messages []ChatMessage `json:"messages" bson:"messages"`
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
