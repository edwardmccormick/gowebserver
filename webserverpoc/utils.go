package main

import (
        "golang.org/x/crypto/bcrypt"
        "github.com/asmarques/geodist"
)

func HashPassword(pw string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    return string(hash)
}

var people = []Person{
	{ID: 0, Name: "Bobby", Motto: "Always Ready", LatLocation: 29.534261019806404, LongLocation: 98.47049550692051, Profile: "https://picsum.photos/200/200"},
	{ID: 1, Name: "Joe", Motto: "Always Faithful", LatLocation: 29.52016959410149, LongLocation: 98.49401109752402, Profile: "https://picsum.photos/150/300"},
	{ID: 2, Name: "Fred", Motto: "Always Prepared", LatLocation: 29.453596593823395, LongLocation: 98.47166788793534, Profile: "https://picsum.photos/250/250"},
	{ID: 3, Name: "Turd Furguson", Motto: "I'm supre close let's party!", LatLocation: 29.419922273698763, LongLocation: 98.48366872664229, Profile: "https://picsum.photos/250/250"},
	{ID: 4, Name: "Don", Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: 77.037, Profile: "https://www.whitehouse.gov/wp-content/uploads/2025/06/President-Donald-Trump-Official-Presidential-Portrait.png"},
	{ID: 5, Name: "The Founder", Motto: "I just want this fucking thing to work", LatLocation: 48.858, LongLocation: -2.294, Profile: "https://ted.mccormickhub.com/img/tedProfilePicture.jpg"},
}

var users = []User{
	{ID: 3, Email: "turd@furguson.com", PasswordHash: HashPassword("password123")},
	{ID: 5, Email: "ted@urmid.com", PasswordHash: HashPassword("password123")},
}

var locationOrigin = geodist.Point{Lat: 29.42618, Long: 98.48618} // The Stupid Alamo
var PictureArray = []string{"https://picsum.photos/250/250", "https://picsum.photos/300/300", "https://picsum.photos/450/300", "https://picsum.photos/450/450", "https://picsum.photos/500/500"}
var photoArray = []ProfilePhoto{
	{Url: "https://picsum.photos/250/250", Caption: "Just me and the boys"},
	{Url: "https://picsum.photos/300/300", Caption: "haha look at their faces"},
	{Url: "https://picsum.photos/450/300", Caption: "omg I can't believe we got away with this"},
	{Url: "https://picsum.photos/450/450", Caption: "life is good man"},
	{Url: "https://picsum.photos/500/500", Caption: "idk haha"},
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

var jwtSecret = []byte("supersecretkey")    // Use a secure random key in production!