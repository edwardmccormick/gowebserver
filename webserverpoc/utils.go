package main

import (
	"time"

	"github.com/asmarques/geodist"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

var people = []Person{
	{ID: 0, Name: "Bobby", Age: 25, Motto: "Always Ready", LatLocation: 29.534261019806404, LongLocation: 98.47049550692051, Profile: "https://picsum.photos/200/200"},
	{ID: 1, Name: "Joe", Age: 31, Motto: "Always Faithful", LatLocation: 29.52016959410149, LongLocation: 98.49401109752402, Profile: "https://picsum.photos/150/300"},
	{ID: 2, Name: "Fred", Age: 28, Motto: "Always Prepared", LatLocation: 29.453596593823395, LongLocation: 98.47166788793534, Profile: "https://picsum.photos/250/250"},
	{ID: 3, Name: "Turd Furguson", Age: 41, Motto: "I'm supre close let's party!", LatLocation: 29.419922273698763, LongLocation: 98.48366872664229, Profile: "./pexels-kelvin809-810775.jpg"},
	{ID: 4, Name: "Don", Age: 79, Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: 77.037, Profile: "https://www.whitehouse.gov/wp-content/uploads/2025/06/President-Donald-Trump-Official-Presidential-Portrait.png"},
	{ID: 5, Name: "The Founder", Age: 42, Motto: "I just want this fucking thing to work", LatLocation: 48.858, LongLocation: -2.294, Profile: "https://ted.mccormickhub.com/img/tedProfilePicture.jpg"},
	{ID: 6, Name: "Training Data", Age: 32, Motto: "Somebody has to carry this fucking site", LatLocation: 48.858, LongLocation: -2.859, Profile: "https://wallpaperheart.com/wp-content/uploads/2018/04/image-Scarlett-Johansson-Images.jpg"},
	{ID: 7, Name: "Scale is Hard", Age: 32, Motto: "I'm famous! Yay!", LatLocation: 50, LongLocation: 0.123, Profile: "./Lenna.jpg"},
	{ID: 8, Name: "Natalie", Age: 39, Motto: "Yes, that one", LatLocation: 49, LongLocation: -1, Profile: "./natalie.jpg"},
}

var users = []User{
	{ID: 0, Email: "bobby@urmid.com", PasswordHash: HashPassword("password123")},
	{ID: 1, Email: "joe@urmid.com", PasswordHash: HashPassword("password123")},
	{ID: 2, Email: "fred@urmid.com", PasswordHash: HashPassword("password123")},
	{ID: 3, Email: "turd@furguson.com", PasswordHash: HashPassword("password123")},
	{ID: 4, Email: "don@furguson.com", PasswordHash: HashPassword("password123")},
	{ID: 5, Email: "ted@urmid.com", PasswordHash: HashPassword("password123")},
	{ID: 6, Email: "training@data.com", PasswordHash: HashPassword("password123")},
	{ID: 7, Email: "scale@urmid.com", PasswordHash: HashPassword("password123")},
	{ID: 8, Email: "natalie@urmid.com", PasswordHash: HashPassword("password123")},
}

var locationOrigin = geodist.Point{Lat: 29.42618, Long: 98.48618} // The Stupid Alamo
var PictureArray = []string{"https://picsum.photos/250/250", "https://picsum.photos/300/300", "https://picsum.photos/450/300", "https://picsum.photos/450/450", "https://picsum.photos/500/500"}
var PhotoArray = []ProfilePhoto{
	{Url: "https://picsum.photos/250/250", Caption: "Just me and the boys"},
	{Url: "https://picsum.photos/300/300", Caption: "haha look at their faces"},
	{Url: "https://picsum.photos/450/300", Caption: "omg I can't believe we got away with this"},
	{Url: "https://picsum.photos/450/450", Caption: "life is good man"},
	{Url: "https://picsum.photos/500/500", Caption: "idk haha"},
}

var PhotoArray2 = []ProfilePhoto{
	{Url: "https://picsum.photos/610/610", Caption: "Just me and the boys"},
	{Url: "https://picsum.photos/300/300", Caption: "haha look at their faces"},
	{Url: "https://picsum.photos/450/300", Caption: "omg I can't believe we got away with this"},
	{Url: "https://picsum.photos/450/450", Caption: "life is good man"},
	{Url: "https://picsum.photos/500/500", Caption: "idk haha"},
	{Url: "https://picsum.photos/600/600", Caption: "Just chilling"},
	{Url: "https://picsum.photos/700/700", Caption: "Having a great time"},
}

var locations = []geodist.Point{
	{Lat: 29.534261019806404, Long: 98.47049550692051}, // SATX
	{Lat: 29.52016959410149, Long: 98.49401109752402},  // NorthStar Mall
	{Lat: 29.453596593823395, Long: 98.47166788793534}, // Doseum
	{Lat: 29.419922273698763, Long: 98.48366872664229}, // Tower of the Americas
	{Lat: 48.858, Long: -2.294},                        // Eiffel Tower
	{Lat: 38.898, Long: 77.037},                        // White House
}

var messages = []string{
	"Hey look at this? Is it working?",
	"Yes, it is working! I can see your message.",
	"Great! I was worried it wasn't working.",
	"Don't worry, it is working. I can see everything you type.",
	"Is it really working? I don't think it is.",
	"Yes, it is working. I can see your messages.",
	"I don't think it is. Can you see what I'm typing? Try again?",
}
var isme = []string{
	"Me",
	"Them",
	"Admin",
}

var Matches = []Match{
	{MatchID: 1000, MatchesIDs: []int{3, 5}, Offered: 3, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 5, AcceptedTime: mustParseTime("2025-06-16T21:03:59.5225862-05:00"), VibeChat: true},
	{MatchID: 1001, MatchesIDs: []int{4, 5}, Offered: 4, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 5, AcceptedTime: mustParseTime("2025-06-16T21:03:59.5225862-05:00"), VibeChat: true},
	{MatchID: 1002, MatchesIDs: []int{1, 3}, Offered: 3, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 1, AcceptedTime: mustParseTime("2025-06-16T21:03:59.5225862-05:00"), VibeChat: true},
	{MatchID: 1003, MatchesIDs: []int{0, 5}, Offered: 0, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 5, VibeChat: true},
	{MatchID: 1004, MatchesIDs: []int{1, 5}, Offered: 5, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 1, VibeChat: true},
	{MatchID: 1005, MatchesIDs: []int{6, 5}, Offered: 6, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 5, VibeChat: true},
	{MatchID: 1006, MatchesIDs: []int{5, 7}, Offered: 5, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 7, VibeChat: true},
}

// mustParseTime is a helper to parse time or panic if invalid
func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		panic(err)
	}
	return t
}

// varOfferedChat: []ChatMessage{
// 	{ID: 0, Time: time.Now(), Who: "Me", Message: "Hey, how's it going?"},

var jwtSecret = []byte("supersecretkey") // Use a secure random key in production!
