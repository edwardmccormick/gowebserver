package main

import (
	"time"
)

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (Person) TableName() string {
	return "people"
}

type Person struct {
	ID           uint      `json:"id" db:"id" gorm:"primaryKey,AutoIncrement,not null,Unique"`
	Name         string    `json:"name" db:"name" gorm:"type:varchar(255),not null"`
	Age          int       `json:"age" db:"age" gorm:"type:int,not null"`
	Motto        string    `json:"motto" db:"motto" gorm:"type:varchar(255)"`
	LatLocation  float64   `json:"lat" db:"lat" gorm:"type:float,not null"`
	LongLocation float64   `json:"long" db:"long" gorm:"type:float,not null"`
	Profile      string    `json:"profile" db:"profile" gorm:"type:varchar(255)"`
	Details      Details   `json:"details" db:"details" gorm:"embedded"`
	Description  string    `json:"description" db:"description" gorm:"type:text"`
	CreateTime   time.Time `json:"create_time" db:"create_time" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"update_time" db:"update_time" gorm:"autoUpdateTime"`
}

type User struct {
	ID           uint   `json:"id" db:"id" gorm:"uniqueIndex;not null"`
	Email        string `json:"email" db:"email" gorm:"type:varchar(255);not null"`
	PasswordHash string `json:"-" db:"password_hash" gorm:"type:varchar(255);not null"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ProfilePhoto struct {
	Url     string `json:"url"`
	Caption string `json:"caption"`
}

type ChatMessage struct {
	ID      int64     `json:"id"`
	Time    time.Time `json:"time"`
	Who     string    `json:"who"`
	Message string    `json:"message"`
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
	MatchID      int           `json:"id"`
	MatchesIDs   []int         `json:"matches_ids"`
	Offered      uint           `json:"offered"`
	OfferedTime  time.Time     `json:"offered_time"`
	Accepted     uint           `json:"accepted"`
	AcceptedTime time.Time     `json:"accepted_time"`
	VibeChat     bool          `json:"vibe_chat"`
	OfferedChat  []ChatMessage `json:"offered_chat"`
	AcceptedChat []ChatMessage `json:"accepted_chat"`
	Person       Person        `json:"person"`
}
