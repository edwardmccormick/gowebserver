package main

import (
	"os"
	"time"

	"github.com/asmarques/geodist"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

var people = []Person{
	// {ID: 0, Name: "Bobby", Age: 25, Motto: "Always Ready", LatLocation: 29.534261019806404, LongLocation: 98.47049550692051, Profile: "./man_square1.jpg"},
	{ID: 1, Name: "Joe", Age: 31, Motto: "Always Faithful", LatLocation: 29.52016959410149, LongLocation: 98.49401109752402, Profile: "./man_square2.jpg"},
	{ID: 2, Name: "Fred", Age: 28, Motto: "Always Prepared", LatLocation: 29.453596593823395, LongLocation: 98.47166788793534, Profile: "/pexels-kelvin809-810775.jpg"},
	{ID: 3, Name: "Turd Furguson", Age: 41, Motto: "It's a funny name", LatLocation: 29.419922273698763, LongLocation: 98.48366872664229, Profile: "./turdfurguson.jpg"},
	{ID: 4, Name: "Don", Age: 79, Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: 77.037, Profile: "https://www.whitehouse.gov/wp-content/uploads/2025/06/President-Donald-Trump-Official-Presidential-Portrait.png"},
	{ID: 5, Name: "The Founder", Age: 42, Motto: "I just want this fucking thing to work", LatLocation: 48.858, LongLocation: -2.294, Profile: "https://ted.mccormickhub.com/img/tedProfilePicture.jpg"},
	{ID: 6, Name: "Training Data", Age: 32, Motto: "Somebody has to carry this fucking site", LatLocation: 48.858, LongLocation: -2.859, Profile: "https://wallpaperheart.com/wp-content/uploads/2018/04/image-Scarlett-Johansson-Images.jpg"},
	{ID: 7, Name: "Scale is Hard", Age: 32, Motto: "I'm famous! Yay!", LatLocation: 50, LongLocation: 0.123, Profile: "./Lenna.jpg"},
	{ID: 8, Name: "Natalie", Age: 39, Motto: "Yes, that one", LatLocation: 49, LongLocation: -1, Profile: "./natalie.jpg"},
}

var users = []User{
	// {ID: 0, Email: "bobby@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Now()},
	{Email: "joe@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "fred@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "turd@furguson.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "don@furguson.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "ted@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "training@data.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "scale@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "natalie@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
}

var albums = []ProfileAlbum{
	{ID: 1, Person: people[0], Photos: PhotoArray},
	{ID: 2, Person: people[1], Photos: PhotoArray},
	{ID: 3, Person: people[2], Photos: PhotoArray},
	{ID: 4, Person: people[3], Photos: PhotoArray},
	{ID: 5, Person: people[4], Photos: PhotoArray},
	{ID: 6, Person: people[5], Photos: PhotoArray2},
	{ID: 7, Person: people[6], Photos: PhotoArray2},
	{ID: 8, Person: people[7], Photos: PhotoArray2},
}

var locationOrigin = geodist.Point{Lat: 29.42618, Long: 98.48618} // The Stupid Alamo
var PictureArray = []string{"/man_square5.jpg", "./man_portrait3.jpg", "./man_landscape5.jpg", "./man_portrait6.jpg", "./man_landscape9.jpg"}
var PhotoArray = []ProfilePhoto{
	{Url: "/man_square5.jpg", Caption: "Just me and the boys"},
	{Url: "./man_portrait3.jpg", Caption: "haha look at their faces"},
	{Url: "./man_landscape5.jpg", Caption: "omg I can't believe we got away with this"},
	{Url: "./man_portrait6.jpg", Caption: "life is good man"},
	{Url: "./man_landscape9.jpg", Caption: "idk haha"},
}

var PhotoArray2 = []ProfilePhoto{
	{Url: "./woman_square5.jpg", Caption: "Just me and the girls"},
	{Url: "./woman_portrait3.jpg", Caption: "haha look at their faces"},
	{Url: "./woman_landscape8.jpg", Caption: "omg I can't believe we got away with this"},
	{Url: "./woman_portrait6.jpg", Caption: "life is good! ❤️"},
	{Url: "./woman_square8.jpg", Caption: "idk haha"},
	{Url: "./woman_landscape4.jpg", Caption: "Just chilling"},
	{Url: "./woman_portrait10.jpg", Caption: "Having a great time"},
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

var Matches = []Match{
	{Offered: 3, OfferedProfile: people[3], OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 5, AcceptedProfile: people[5], AcceptedTime: mustParseTime("2025-06-16T21:04:59.5225862-05:00"), VibeChat: true},
	{Offered: 4, OfferedProfile: people[4], OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 5, AcceptedProfile: people[5], AcceptedTime: mustParseTime("2025-06-16T21:04:59.5225862-05:00"), VibeChat: true},
	{Offered: 3, OfferedProfile: people[3], OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 1, AcceptedProfile: people[1], AcceptedTime: mustParseTime("2025-06-16T21:04:59.5225862-05:00"), VibeChat: true},
	// {MatchID: 1003, Offered: 3, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 5,  AcceptedTime: mustParseTime("0000-00-1T00:00:0.0000001-05:00"), VibeChat: true},
	{Offered: 5, OfferedProfile: people[5], OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 1, AcceptedProfile: people[1], AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), VibeChat: true},
	{Offered: 6, OfferedProfile: people[6], OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 5, AcceptedProfile: people[5], AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), VibeChat: true},
	{Offered: 5, OfferedProfile: people[5], OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 7, AcceptedProfile: people[7], AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), VibeChat: true},
}

// mustParseTime is a helper to parse time or panic if invalid
func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		panic(err)
	}
	return t
}

var jwtSecret = []byte("supersecretkey") // Use a secure random key in production!

func isRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

func intPtr(i int) *int {
	return &i
}
