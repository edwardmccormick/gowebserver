package main

import (
	"time"
)

type Person struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Motto        string  `json:"motto"`
	LatLocation  float64 `json:"lat"`
	LongLocation float64 `json:"long"`
	Profile      string  `json:"profile"`
	Details      Details `json:"details"`
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
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
	Offered      int           `json:"offered"`
	OfferedTime  time.Time     `json:"offered_time"`
	Accepted     int           `json:"accepted"`
	AcceptedTime time.Time     `json:"accepted_time"`
	VibeChat     bool          `json:"vibe_chat"`
	OfferedChat  []ChatMessage `json:"offered_chat"`
	AcceptedChat []ChatMessage `json:"accepted_chat"`
	Person       Person        `json:"person"`
}
