package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type person struct {
    ID     string  `json:"id"`
    Name  string  `json:"name"`
    Motto string  `json:"motto"`
    LatLocation  float64 `json:"lat"`
	LongLocation  float64 `json:"long"`
}

func main() {
    router := gin.Default()
    router.GET("/people", getPeople)
	router.GET("/people/:id", getPeopleByID)
	router.POST("/people", postPeople)
	router.GET("/", greetUser)
	router.GET("/greet/:name", greetUserByName)

    router.Run("localhost:8080")
}

// albums slice to seed record album data.
var people = []person{
    {ID: "1", Name: "Bobby", Motto: "Always Ready", LatLocation: 0, LongLocation: 0},
    {ID: "2", Name: "Joe", Motto: "Always Faithful", LatLocation: 0, LongLocation: 0},
    {ID: "3", Name: "Fred", Motto: "Always Prepared", LatLocation: 0, LongLocation: 0},
}

func getPeople(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, people)
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