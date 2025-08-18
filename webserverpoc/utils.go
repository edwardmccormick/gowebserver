package main

import (
	"math"
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
	{ID: 1, Name: "Marco", Age: 32, Motto: "Civil Engineer & Aspiring Pitmaster", LatLocation: 29.52016959410149, LongLocation: -98.49401109752402, Profile: ProfilePhoto{S3Key: "1/profile"}, Photos: deepCopyPhotoArray(PhotoArray1), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Hey there, I’m Marco. I spend my weekdays designing bridges and roadways to connect our city, and I'd like to think I’m pretty good at building things that last. I'm a logical person by trade, but I have a huge creative side that comes out on the weekends—usually in the form of a massive brisket on my smoker. There's a certain kind of science to getting the bark just right, a challenge I take very seriously.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"When I’m not working or pretending to be a top chef, I'm usually outdoors. I love hiking the trails at Government Canyon State Natural Area or just taking a long walk with my dog, a goofy golden retriever named Churro. I’m equally happy getting dressed up for a nice dinner and cocktails at The Pearl, or keeping it casual with a beer and some live music at a place like Floore's Country Store. I appreciate routine and stability, but I’m always pushing myself to try new things and see new places.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I’m looking for a partner who is kind, intelligent, and has a passion of their own. Someone who is as comfortable planning a weekend getaway to the Hill Country as they are having a quiet night in. Communication and mutual respect are huge for me. I’m looking for my best friend and my biggest supporter, all wrapped into one.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I believe the best relationships, much like the best bridges, are built on a solid foundation of trust, shared experiences, and a willingness to weather any storm together. I want to build something real and lasting with the right person.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"If you’re a woman who knows what she wants, isn’t afraid of a little BBQ smoke, and thinks you could put up with my engineering puns, I’d love to hear from you. Tell me, <name>, what's the best thing you've eaten in San Antonio lately?\"},{\"insert\":\"\\n\"}]}"},
	{ID: 2, Name: "Gabriel", Age: 38, Motto: "Single Dad & Project Manager", LatLocation: 29.453596593823395, LongLocation: -98.47166788793534, Profile: ProfilePhoto{S3Key: "2/profile"}, Photos: deepCopyPhotoArray(PhotoArray2), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Alright, let’s be upfront. I'm a dad to an amazing 8-year-old son who is my entire world. My life revolves around school pickups, little league games, and an impressive amount of LEGOs. I wouldn't have it any other way.\"},{\"insert\":\"\\n\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"When I'm not on dad duty, I work as a project manager for a local construction firm. On my own time, I enjoy the simple things: grilling in the backyard, listening to classic rock, and finding a good patio to enjoy a cold Shiner. I'm a pretty low-key, drama-free guy who is honest, loyal, and looking for the same.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I'm ready to find a partner for myself again. Someone who is patient, kind, and understands that my son will always be my first priority. If you’re a woman who is independent, has her life together, and doesn't mind the occasional conversation about Minecraft, I'd love to talk.\"},{\"insert\":\"\\n\"}]}"},
	{ID: 3, Name: "Turd Furguson", Age: 26, Motto: "Just trying to survive the Texas heat", LatLocation: 29.419922273698763, LongLocation: -98.48366872664229, Profile: ProfilePhoto{S3Key: "3/profile"}, Photos: deepCopyPhotoArray(PhotoArray3), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Recently stationed at JBSA-Lackland. Originally from Oregon, so I'm still getting used to the fact that there are only two seasons here: hot and not-as-hot.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Need a tour guide to show me the best spots that aren't the River Walk. Teach me about Fiesta. Show me where the best puffy tacos are. In return, I can reach things on the top shelf and tell you all about evergreen trees.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Let’s go on an adventure.\"},{\"insert\":\"\\n\"}]}"},
	{ID: 4, Name: "Don", Age: 79, Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: -77.037, Profile: ProfilePhoto{S3Key: "4/profile"}, Photos: deepCopyPhotoArray(PhotoArray4), Description: "{\"ops\":[{\"insert\":\"The Radical Left Democrats have hit pay dirt, again! Just like with the FAKE and fully discredited Steele Dossier, the lying 51 “Intelligence” Agents, the Laptop from Hell, which the Dems swore had come from Russia (No, it came from Hunter Biden's bathroom!), and even the Russia, Russia, Russia Scam itself, a totally fake and made up story used in order to hide Crooked Hillary Clinton's big loss in the 2016 Presidential Election, these Scams and Hoaxes are all the Democrats are good at - It's all they have - They are no good at governing, no good at policy, and no good at picking winning candidates. Also, unlike Republicans, they stick together like glue. Their new SCAM is what we will forever call the Jeffrey Epstein Hoax, and my PAST supporters have bought into this “bullshit,” hook, line, and sinker. They haven't learned their lesson, and probably never will, even after being conned by the Lunatic Left for 8 long years. I have had more success in 6 months than perhaps any President in our Country's history, and all these people want to talk about, with strong prodding by the Fake News and the success starved Dems, is the Jeffrey Epstein Hoax. Let these weaklings continue forward and do the Democrats work, don't even think about talking of our incredible and unprecedented success, because I don't want their support anymore! Thank you for your attention to this matter. MAKE AMERICA GREAT AGAIN!\\n\"}]}"},
	{ID: 5, Name: "The Founder", Age: 42, Motto: "I just want this fucking thing to work", LatLocation: 38.858, LongLocation: -98.294, Profile: ProfilePhoto{S3Key: "5/profile"}, Photos: deepCopyPhotoArray(PhotoArray5), Description: "{\"ops\":[{\"insert\":\"All silliness aside, I used to be in the Navy, and now I'm a software developer. I'm the proud dad to a pair of awesome tweens, so usually if I'm not at work I'm chasing them - or our two dogs. But I dream of the day of being able to head over to Habibi Cafe and smoking some hookah over a cup of coffee with someone cool.\\n\"}]}"},
	{ID: 6, Name: "Cassandra", Age: 34, Motto: "Professor at Trinity", LatLocation: 38.858, LongLocation: -92.859, Profile: ProfilePhoto{S3Key: "6/profile"}, Photos: deepCopyPhotoArray(PhotoArray6), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Professor at Trinity\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Fluent in sarcasm, APA style, and the art of finding parking near campus. I teach literature, which means I get paid to overthink things and argue about what authors really meant. It's a pretty sweet gig.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Looking for someone who reads, can hold a conversation that goes deeper than the weather, and won't be intimidated if I use a semicolon in a text message. Let's debate whether breakfast tacos are better than regular tacos over a bottle of wine. So...what's your favorite book?\"},{\"insert\":\"\\n\"}]}"},
	{ID: 7, Name: "Dr. Aisha", Age: 32, Motto: "Fixing Babies! Pediatric Resident", LatLocation: 50, LongLocation: -100, Profile: ProfilePhoto{S3Key: "7/profile"}, Photos: deepCopyPhotoArray(PhotoArray7), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"By day, I'm navigating the controlled chaos of the pediatric wing at University Hospital, powered by coffee and the incredible resilience of my tiny patients. It’s a demanding job that I absolutely love, but it definitely makes me appreciate my downtime. My life isn't a TV show—it's mostly paperwork and trying to make kids laugh—but it's incredibly rewarding.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"When I finally shed the scrubs, I'm all about decompressing. I love jogging along the Mission Reach trail, losing myself in a good book at a coffee shop in Southtown, or exploring the city's food scene. I'm on a personal quest to find the absolute best margarita in town (current front-runner is at Rosario's, but I'm open to suggestions). I’m looking for someone who understands a busy schedule but is excited to make the time we do have together count.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I’m seeking a genuine connection with someone who is compassionate, has a great sense of humor (a requirement for my line of work), and is emotionally mature. Bonus points if you can make me laugh after a 12-hour shift.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"Let's grab a drink and see if we click.\"},{\"insert\":\"\\n\"}]}"},
	{ID: 8, Name: "Sofia", Age: 31, Motto: "Food Blogger & Marketing Guru", LatLocation: 49, LongLocation: -98.5, Profile: ProfilePhoto{S3Key: "8/profile"}, Photos: deepCopyPhotoArray(PhotoArray8), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"My love language is sharing a plate of perfectly crispy nachos. I'm a marketer by day and a San Antonio foodie by night (and weekend). My camera roll is 90% food pictures, and I have a running list of restaurants I need to try.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I'm looking for a partner-in-crime for culinary adventures. Someone who gets excited about trying a new hole-in-the-wall spot, who understands the profound importance of a good salsa, and who will never say \\\"no\\\" to splitting a dessert.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"If your idea of a perfect date involves a food truck park or a farmer's market, we'll get along just fine.\"},{\"insert\":\"\\n\"}]}"},
	{ID: 9, Name: "David", Age: 29, Motto: "Just a guy.", LatLocation: 29.523, LongLocation: -98.254, Profile: ProfilePhoto{S3Key: "9/profile"}, Photos: deepCopyPhotoArray(PhotoArray2), Description: "{\"ops\":[{\"insert\":\"Spurs. Fishing. Grilling. Not complicated. Let’s grab a beer at The Friendly Spot.}]},{\"insert\":\"\\n\"}]}"},
	{ID: 10, Name: "Elena", Age: 28, Motto: "Graphic Designer", LatLocation: 29.583, LongLocation: -98.454, Profile: ProfilePhoto{S3Key: "10/profile"}, Photos: deepCopyPhotoArray(PhotoArray6), Description: "{\"ops\":[{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I see the world in color palettes and fonts. As a graphic designer, I love finding beauty in unexpected places, whether it’s the street art in the St. Mary's Strip or the architecture downtown. My ideal Sunday involves a latte from Merit, a walk through the Japanese Tea Garden, and maybe a stop at the McNay.\"},{\"insert\":\"\\n\"},{\"attributes\":{\"color\":\"#000000\"},\"insert\":\"I'm looking for a creative soul who is kind, curious, and open-minded. Someone who can appreciate a First Friday art walk and also isn't afraid to be silly. Let's create something beautiful together.\"},{\"insert\":\"\\n\"}]}"},
	{ID: 11, Name: "Admin", Age: 199, Motto: "You shouldn't ever see this", LatLocation: 90, LongLocation: -90, Profile: ProfilePhoto{S3Key: "11/profile"}, Photos: deepCopyPhotoArray(PhotoArray2), Description: "{\"ops\":[{\"insert\":\"Spurs. Fishing. Grilling. Not complicated. Let’s grab a beer at The Friendly Spot.}]},{\"insert\":\"\\n\"}]}"},
}

var users = []User{
	// {ID: 0, Email: "bobby@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Now()},
	{Email: "marco@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "gabriel@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "turd@furguson.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "don@furguson.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "ted@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "cassandra@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "aisha@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "sofia@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "david@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "elena@urmid.com", PasswordHash: HashPassword("password123"), LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
	{Email: "admin@urmid.com", PasswordHash: HashPassword("password123"), IsAdmin: true, LastLogin: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
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

// CalculateHaversineDistance calculates the distance between two points (lat/long) using the Haversine formula.
// Returns distance in miles.
func CalculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth radius in miles
	const earthRadius = 3958.8

	// Convert degrees to radians
	rlat1 := lat1 * math.Pi / 180
	rlat2 := lat2 * math.Pi / 180
	rlon1 := lon1 * math.Pi / 180
	rlon2 := lon2 * math.Pi / 180

	// Haversine formula
	dlon := rlon2 - rlon1
	dlat := rlat2 - rlat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(rlat1)*math.Cos(rlat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// CalculateBoundingBox calculates a bounding box around a point with the given radius in miles.
// Returns minLat, maxLat, minLong, maxLong
func CalculateBoundingBox(lat, lon, radiusMiles float64) (float64, float64, float64, float64) {
	// Earth's radius in miles
	const earthRadius = 3958.8

	// Convert to radians
	latRad := lat * math.Pi / 180

	// Calculate the latitude range in degrees
	latRange := (radiusMiles / earthRadius) * (180 / math.Pi)

	// Calculate the longitude range, which varies with latitude
	longRange := (radiusMiles / (earthRadius * math.Cos(latRad))) * (180 / math.Pi)

	return lat - latRange, lat + latRange, lon - longRange, lon + longRange
}
