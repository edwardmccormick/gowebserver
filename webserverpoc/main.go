package main

import (
	"fmt"
	"math"
	"net/http"
	"sort"

	"github.com/asmarques/geodist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type person struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Motto        string  `json:"motto"`
	LatLocation  float64 `json:"lat"`
	LongLocation float64 `json:"long"`
	Dogs         int     `json:"dogs"`
	Profile      string  `json:"profile"`
}

type processedProfile struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Motto    string  `json:"motto"`
	Distance float64 `json:"distance"`
	Dogs     int     `json:"dogs"`
	Profile  string  `json:"profile"`
}

type profilePhoto struct {
	Url     string `json:"url"`
	Caption string `json:"caption"`
}

func main() {
	router := gin.Default()
	router.Use(cors.Default()) // add this line
	fmt.Println(geodist.VincentyDistance(locationOrigin, locations[0]))
	fmt.Println(geodist.VincentyDistance(locationOrigin, locations[1]))
	fmt.Println(geodist.VincentyDistance(locationOrigin, locations[2]))
	fmt.Println(geodist.VincentyDistance(locationOrigin, locations[3]))
	fmt.Println(geodist.VincentyDistance(locationOrigin, locations[4]))
	fmt.Println(geodist.VincentyDistance(locationOrigin, locations[5]))
	fmt.Println(geodist.VincentyDistance(locations[4], locations[5])) // 6177.45 WH -> Eiffel Tower
	router.GET("/people", getProcessedPeople)
	router.GET("/people/:id", getPeopleByID)
	router.POST("/people", postPeople)
	router.GET("/", greetUser)
	router.GET("/greet/:name", greetUserByName)
	router.GET("/login", login)
	router.GET("/photos/:id", getPhotosByID)

	router.Run("localhost:8080")
}

// albums slice to seed record album data.
var people = []person{
	{ID: "0", Name: "Bobby", Motto: "Always Ready", LatLocation: 29.534261019806404, LongLocation: 98.47049550692051, Dogs: 4, Profile: "https://picsum.photos/200/200"},
	{ID: "1", Name: "Joe", Motto: "Always Faithful", LatLocation: 29.52016959410149, LongLocation: 98.49401109752402, Dogs: 10, Profile: "https://picsum.photos/150/300"},
	{ID: "2", Name: "Fred", Motto: "Always Prepared", LatLocation: 29.453596593823395, LongLocation: 98.47166788793534, Dogs: 0, Profile: "https://picsum.photos/250/250"},
	{ID: "3", Name: "Turd Furguson", Motto: "I'm supre close let's party!", LatLocation: 29.419922273698763, LongLocation: 98.48366872664229, Dogs: 2, Profile: "https://picsum.photos/250/250"},
	{ID: "4", Name: "Don", Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: 77.037, Dogs: 0, Profile: "https://www.whitehouse.gov/wp-content/uploads/2025/06/President-Donald-Trump-Official-Presidential-Portrait.png"},
}

var locationOrigin = geodist.Point{Lat: 29.42618, Long: 98.48618}
var pictureArray = []string{"https://picsum.photos/250/250", "https://picsum.photos/300/300", "https://picsum.photos/450/300", "https://picsum.photos/450/450", "https://picsum.photos/500/500"}
var photoArray = []profilePhoto{
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

// 29.456001687343456, -98.471976423337 DoSeum?
// 29.53502981348254, -98.47082162761036 SATX
// 29.52016959410149, -98.49401109752402 NorthStar Cowboy Boots
// 29.45666986455716, -98.70002021204594 Seaworld
// 29.426057444276093, -98.4861282632775 Alamo

func getPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, people)
}

func getProcessedPeople(c *gin.Context) {
	var processedPeople []processedProfile

	for _, p := range people {

		// Calculate the distance from locationOrigin
		distance, err := geodist.VincentyDistance(locationOrigin, geodist.Point{Lat: p.LatLocation, Long: p.LongLocation})
		if err != nil {
			fmt.Printf("Error calculating distance for person %s: %v\n", p.ID, err)
			continue
		}

		// Round the distance to two decimal places
		roundedDistance := math.Round(distance*100) / 100

		// Create a processedProfile and append it to the processedPeople slice
		processedPeople = append(processedPeople, processedProfile{
			ID:       p.ID,
			Name:     p.Name,
			Motto:    p.Motto,
			Distance: roundedDistance,
			Dogs:     p.Dogs,
			Profile:  p.Profile,
		})

		sort.Slice(processedPeople, func(i, j int) bool {
			return processedPeople[i].Distance < processedPeople[j].Distance
		})
	}

	// Return the processedPeople array as JSON
	c.IndentedJSON(http.StatusOK, processedPeople)
}

func greetUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

func greetUserByName(c *gin.Context) {
	name := c.Param("name")

	c.String(http.StatusOK, "Hello %s", name)
}

func getPeopleByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range people {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getPhotosByID(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	if id == "0" || id == "1" || id == "2" || id == "3" || id == "4" || id == "5" || id == "6" || id == "7" || id == "8" || id == "9" || id == "10" {
		c.IndentedJSON(http.StatusOK, photoArray)
		return
	}
	// TODO: Implement a Document DB instance and associated call
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found for that id"})
}

func postPeople(c *gin.Context) {
	var newPerson person

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	// Add the new album to the slice.
	people = append(people, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func login(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "TODO: Implement Auth!")
	return
}
