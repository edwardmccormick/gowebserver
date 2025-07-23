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
	{ID: 1, Name: "Joe", Age: 31, Motto: "Always Faithful", LatLocation: 29.52016959410149, LongLocation: -98.49401109752402, Profile: ProfilePhoto{S3Key: "1/profile"}, Photos: deepCopyPhotoArray(PhotoArray1)},
	{ID: 2, Name: "Fred", Age: 28, Motto: "Always Prepared", LatLocation: 29.453596593823395, LongLocation: -98.47166788793534, Profile: ProfilePhoto{S3Key: "2/profile"}, Photos: deepCopyPhotoArray(PhotoArray2)},
	{ID: 3, Name: "Turd Furguson", Age: 41, Motto: "It's a funny name", LatLocation: 29.419922273698763, LongLocation: -98.48366872664229, Profile: ProfilePhoto{S3Key: "3/profile"}, Photos: deepCopyPhotoArray(PhotoArray3)},
	{ID: 4, Name: "Don", Age: 79, Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: -77.037, Profile: ProfilePhoto{S3Key: "4/profile"}, Photos: deepCopyPhotoArray(PhotoArray4)},
	{ID: 5, Name: "The Founder", Age: 42, Motto: "I just want this fucking thing to work", LatLocation: 38.858, LongLocation: -92.294, Profile: ProfilePhoto{S3Key: "5/profile"}, Photos: deepCopyPhotoArray(PhotoArray5)},
	{ID: 6, Name: "Training Data", Age: 32, Motto: "Somebody has to carry this fucking site", LatLocation: 38.858, LongLocation: -92.859, Profile: ProfilePhoto{S3Key: "6/profile"}, Photos: deepCopyPhotoArray(PhotoArray6)},
	{ID: 7, Name: "Scale is Hard", Age: 32, Motto: "I'm famous! Yay!", LatLocation: 50, LongLocation: -100, Profile: ProfilePhoto{S3Key: "7/profile"}, Photos: deepCopyPhotoArray(PhotoArray7)},
	{ID: 8, Name: "Natalie", Age: 39, Motto: "Yes, that one", LatLocation: 49, LongLocation: -100, Profile: ProfilePhoto{S3Key: "8/profile"}, Photos: deepCopyPhotoArray(PhotoArray8)},
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

// var albums = []ProfileAlbum{
// 	{ID: 1, Person: people[0], Photos: PhotoArray},
// 	{ID: 2, Person: people[1], Photos: PhotoArray},
// 	{ID: 3, Person: people[2], Photos: PhotoArray},
// 	{ID: 4, Person: people[3], Photos: PhotoArray},
// 	{ID: 5, Person: people[4], Photos: PhotoArray},
// 	{ID: 6, Person: people[5], Photos: PhotoArray2},
// 	{ID: 7, Person: people[6], Photos: PhotoArray2},
// 	{ID: 8, Person: people[7], Photos: PhotoArray2},
// }

// var locationOrigin = geodist.Point{Lat: 29.42618, Long: 98.48618} // The Stupid Alamo
var PictureArray = []string{"/man_square5.jpg", "./man_portrait3.jpg", "./man_landscape5.jpg", "./man_portrait6.jpg", "./man_landscape9.jpg"}
var PhotoArray1 = []ProfilePhoto{
	{Url: "/man_square5.jpg", S3Key: "example/man_square5.jpg", Caption: "Just me and the boys"},
	{Url: "./man_portrait3.jpg", S3Key: "example/man_portrait3.jpg", Caption: "haha look at their faces"},
	{Url: "./man_landscape5.jpg", S3Key: "example/man_landscape5.jpg", Caption: "omg I can't believe we got away with this"},
	{Url: "./man_portrait6.jpg", S3Key: "example/man_portrait6.jpg", Caption: "life is good man"},
	{Url: "./man_landscape7.jpg", S3Key: "example/man_landscape7.jpg", Caption: "idk haha"},
}

var PhotoArray2 = []ProfilePhoto{
	{Url: "/man_square5.jpg", S3Key: "example/man_square5.jpg", Caption: "Just me and the boys"},
	{Url: "./man_portrait3.jpg", S3Key: "example/man_portrait3.jpg", Caption: "haha look at their faces"},
	{Url: "./man_landscape5.jpg", S3Key: "example/man_landscape5.jpg", Caption: "omg I can't believe we got away with this"},
	{Url: "./man_portrait6.jpg", S3Key: "example/man_portrait6.jpg", Caption: "life is good man"},
	{Url: "./man_landscape7.jpg", S3Key: "example/man_landscape7.jpg", Caption: "idk haha"},
}

var PhotoArray3 = []ProfilePhoto{
	{Url: "/man_square5.jpg", S3Key: "example/man_square5.jpg", Caption: "Just me and the boys"},
	{Url: "./man_portrait3.jpg", S3Key: "example/man_portrait3.jpg", Caption: "haha look at their faces"},
	{Url: "./man_landscape5.jpg", S3Key: "example/man_landscape5.jpg", Caption: "omg I can't believe we got away with this"},
	{Url: "./man_portrait6.jpg", S3Key: "example/man_portrait6.jpg", Caption: "life is good man"},
	{Url: "./man_landscape7.jpg", S3Key: "example/man_landscape7.jpg", Caption: "idk haha"},
}

var PhotoArray4 = []ProfilePhoto{
	{Url: "/man_square5.jpg", S3Key: "4/image1", Caption: "I call this one 'Blue Steel'"},
	{Url: "./man_portrait3.jpg", S3Key: "4/image2", Caption: "so hot it almost melted my makeup lol"},
	{Url: "./man_landscape5.jpg", S3Key: "4/image3", Caption: "i'm a little teapot, short and stout"},
}

var PhotoArray5 = []ProfilePhoto{
	{Url: "/man_square5.jpg", S3Key: "5/image1", Caption: "I'm on the left"},
	{Url: "./man_portrait3.jpg", S3Key: "5/image2", Caption: "man, that was so long ago lol"},
	{Url: "./man_landscape5.jpg", S3Key: "example/man_landscape5.jpg", Caption: "omg I can't believe we got away with this"},
	{Url: "./man_portrait6.jpg", S3Key: "example/man_portrait6.jpg", Caption: "life is good man"},
	{Url: "./man_landscape7.jpg", S3Key: "example/man_landscape7.jpg", Caption: "idk haha"},
}

var PhotoArray6 = []ProfilePhoto{
	{Url: "./woman_square5.jpg", S3Key: "6/image1", Caption: "photo booth!"},
	{Url: "./woman_portrait3.jpg", S3Key: "6/image2", Caption: "omg I look high!"},
	{Url: "./woman_landscape8.jpg", S3Key: "6/image3", Caption: "red hair don't care"},
	{Url: "./woman_portrait6.jpg", S3Key: "6/image4", Caption: "so. casual. and. comfy. ❤️"},
	{Url: "./woman_square8.jpg", S3Key: "6/image5", Caption: "not casual, and super uncomfy, but I love this outfit"},
	{Url: "./woman_landscape4.jpg", S3Key: "6/image6", Caption: "this was so much fun"},
	{Url: "./woman_landscape4.jpg", S3Key: "6/image7", Caption: "love my girls!"},
}

var PhotoArray7 = []ProfilePhoto{
	{Url: "./woman_square5.jpg", S3Key: "7/image1", Caption: "no makeup....but maybe a *little* filter lol"},
	{Url: "./woman_portrait3.jpg", S3Key: "7/image2", Caption: "i'm not REALLY blonde I promise"},
	{Url: "./woman_landscape8.jpg", S3Key: "7/image3", Caption: "omg my freckles lol"},
	{Url: "./woman_portrait6.jpg", S3Key: "7/image4", Caption: "just smile, my friend said, you're so cute you can't take a bad picture lolol ❤️"},
	{Url: "./woman_square8.jpg", S3Key: "7/image5", Caption: "idk haha, love that sweater though"},
}

var PhotoArray8 = []ProfilePhoto{
	{Url: "./woman_square5.jpg", S3Key: "8/image1", Caption: "this was at a dinner or something"},
	{Url: "./woman_portrait3.jpg", S3Key: "8/image2", Caption: "haha bedroom eyes!"},
	{Url: "./woman_landscape8.jpg", S3Key: "8/image3", Caption: "roses are red, violets are blue, I'm wearing roses, how about you?"},
	{Url: "./woman_portrait6.jpg", S3Key: "8/image4", Caption: "my ballet phase! ❤️"},
	{Url: "./woman_square8.jpg", S3Key: "8/image5", Caption: "idk haha"},
	{Url: "./woman_landscape4.jpg", S3Key: "8/image6", Caption: "hate the blog, love the look!"},
	{Url: "./woman_portrait10.jpg", S3Key: "8/image7", Caption: "me and my bestie"},
	{Url: "./woman_portrait10.jpg", S3Key: "8/image8", Caption: "just relaxing; might delete later!"},
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
	{Offered: 3, OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 5, AcceptedTime: mustParseTime("2025-06-16T21:04:59.5225862-05:00"), VibeChat: true},
	{Offered: 4, OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 5, AcceptedTime: mustParseTime("2025-06-16T21:04:59.5225862-05:00"), VibeChat: true},
	{Offered: 3, OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 1, AcceptedTime: mustParseTime("2025-06-16T21:04:59.5225862-05:00"), VibeChat: true},
	// {MatchID: 1003, Offered: 3, OfferedTime: mustParseTime("2025-06-16T21:03:56.5225862-05:00"), Accepted: 5,  AcceptedTime: mustParseTime("0000-00-1T00:00:0.0000001-05:00"), VibeChat: true},
	{Offered: 5, OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 1, AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), VibeChat: true},
	{Offered: 6, OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 5, AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), VibeChat: true},
	{Offered: 5, OfferedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Accepted: 7, AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), VibeChat: true},
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

// Helper function to create a deep copy of a slice
func deepCopyPhotoArray(original []ProfilePhoto) []ProfilePhoto {
	copy := make([]ProfilePhoto, len(original))
	for i, photo := range original {
		copy[i] = photo
	}
	return copy
}
