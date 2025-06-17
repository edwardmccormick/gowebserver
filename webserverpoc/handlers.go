package main

import (
	// "encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/asmarques/geodist"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Add your route handlers here
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Find the user by email
	var user *User
	for i := range users {
		if users[i].Email == req.Email {
			user = &users[i]
			break
		}
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"email": user.Email,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// Return token and user info (excluding password hash)
	resp := struct {
		Token  string `json:"token"`
		Person Person `json:"person"`
	}{
		Token:  tokenString,
		Person: people[user.ID],
	}

	c.JSON(http.StatusOK, resp)
}

func Signout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	// Parse and validate the token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Remove the refresh token
	delete(RefreshTokens, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func GetPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, people)
}

func PostPeople(c *gin.Context) {
	var newPerson Person
	// var newDetails Details

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	// // Ensure the Details field is initialized if it's missing
	// if newPerson.Details == nil {
	//     newPerson.Details = c.BindJSON(&newDetails);
	// 	err != nil {
	// 		return
	// 	}
	// }

	// Add the new album to the slice.
	people = append(people, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func GetPeopleByLocation(c *gin.Context) {
	var processedPeople []ProcessedProfile
	JwtMiddleware(c)
	var req struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
	}
	fmt.Println(c)
	fmt.Println(req)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	lat, err := strconv.ParseFloat(req.Lat, 64)
	if err != nil {
		fmt.Println("Error parsing latitude:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}
	fmt.Println(lat)

	long, err := strconv.ParseFloat(req.Long, 64)
	if err != nil {
		fmt.Println("Error parsing longitude:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}
	fmt.Println(long)
	for _, p := range people {

		// Calculate the distance from locationOrigin
		distance, err := geodist.VincentyDistance(geodist.Point{Lat: p.LatLocation, Long: p.LongLocation}, geodist.Point{Lat: lat, Long: long})
		if err != nil {
			fmt.Printf("Error calculating distance for person %s: %v\n", p.ID, err)
			continue
		}

		// Round the distance to two decimal places
		roundedDistance := math.Round(distance*100) / 100

		// Create a processedProfile and append it to the processedPeople slice
		processedPeople = append(processedPeople, ProcessedProfile{
			ID:       p.ID,
			Name:     p.Name,
			Age:      p.Age,
			Motto:    p.Motto,
			Distance: roundedDistance,
			Profile:  p.Profile,
			Details:  p.Details,
		})

		sort.Slice(processedPeople, func(i, j int) bool {
			return processedPeople[i].Distance < processedPeople[j].Distance
		})
	}

	// Return the processedPeople array as JSON
	c.IndentedJSON(http.StatusOK, processedPeople)
}

func GetProcessedPeople(c *gin.Context) {
	var processedPeople []ProcessedProfile

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
		processedPeople = append(processedPeople, ProcessedProfile{
			ID:       p.ID,
			Name:     p.Name,
			Age:      p.Age,
			Motto:    p.Motto,
			Distance: roundedDistance,
			Profile:  p.Profile,
			Details:  p.Details,
		})

		sort.Slice(processedPeople, func(i, j int) bool {
			return processedPeople[i].Distance < processedPeople[j].Distance
		})
	}

	// Return the processedPeople array as JSON
	c.IndentedJSON(http.StatusOK, processedPeople)
}

func GetPeopleByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(id)

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range people {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
}

func GetPhotosByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(id)
	if id == 0 || id == 1 || id == 2 || id == 3 {
		c.IndentedJSON(http.StatusOK, PhotoArray2)
		return
	}
	if id == 4 || id == 5 || id == 6 || id == 7 || id == 8 || id == 9 || id == 10 {
		c.IndentedJSON(http.StatusOK, PhotoArray)
		return
	}
	// TODO: Implement a Document DB instance and associated call
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found for that id"})
}

func GreetUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

func GreetUserByName(c *gin.Context) {
	name := c.Param("name")

	c.String(http.StatusOK, "Hello %s", name)
}

func GetFaviconIco(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Stop hitting endpoints with your stupid browser ya goof!🤣🤷‍♂️")
}

func GetMatches(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Matches)
}

func GetMatchByID(c *gin.Context) {
	// Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Find the match with the given ID
	for _, match := range Matches {
		if match.MatchID == id {
			c.IndentedJSON(http.StatusOK, match)
			return
		}
	}

	// If no match is found, return a 404 Not Found response
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Match not found"})
}

func GetMatchByPersonID(c *gin.Context) {
	// Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	matchesForPerson := []Match{}

	// Find the match with the given ID

	// Find the match with the given ID
	for _, match := range Matches {
		if match.Offered == id || match.Accepted == id {
			if !match.AcceptedTime.IsZero() { // Check if AcceptedTime is not null
				// Determine the other person's ID
				otherPersonID := match.Offered
				if match.Offered == id {
					otherPersonID = match.Accepted
				}

				// Find the person in the people array
				for _, person := range people {
					if person.ID == otherPersonID {
						match.Person = person // Add the person to match.Profile
						break
					}
				}
			}
			matchesForPerson = append(matchesForPerson, match)
		}
	}

	if len(matchesForPerson) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No matches found for this person"})
		return
	}

	c.IndentedJSON(http.StatusOK, matchesForPerson)

}
