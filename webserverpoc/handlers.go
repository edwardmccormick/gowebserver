package main

import (
	// "encoding/json"
	"fmt"
	"math"
	"net/http"
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
	db.Where("email = ?", req.Email).First(&user)
	// for i := range users {
	// 	if users[i].Email == req.Email {
	// 		// users[i].LastLogin = time.Now()
	// 		user = &users[i]
	// 		break
	// 	}
	// }
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


	var person = Person{ID: user.ID}
	db.First(&person)
	// Return token and user info (excluding password hash)
	resp := struct {
		Token  string `json:"token"`
		Person Person `json:"person"`
	}{
		Token:  tokenString,
		Person: person,
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

	// claims := parsedToken.Claims.(jwt.MapClaims)

	// Remove the refresh token
	ExpiredTokens = append(ExpiredTokens, token)
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func GetPeople(c *gin.Context) {
	var people []Person
	if result := db.Find(&people); result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, people)
}

func GetUsers(c *gin.Context) {
	var users []User
	if result := db.Find(&users); result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, users)
}

func PostPeople(c *gin.Context) {
	var newPerson Person

	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	
	result := db.Create(newPerson) // pass a slice to insert multiple row
	fmt.Println("Created rows: ", result.RowsAffected)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func GetPeopleByLocation(c *gin.Context) {
	var processedPeople []Person
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
		fmt.Println(roundedDistance)
		// Create a processedProfile and append it to the processedPeople slice
		processedPeople = append(processedPeople, p)

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

	var person = Person{ID: uint(id)}
	results := db.First(&person)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, person)
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
	c.IndentedJSON(http.StatusOK, "localhost:3306/urmid.svg")
}

func GetMatches(c *gin.Context) {
	var matches []Match
	if result := db.Find(&matches); result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, matches)
}

func GetMatchByID(c *gin.Context) {
	// Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var matches []Match
	results := db.Where("offered = ?", uint(id)).Or("accepted = ?", uint(id)).Preload("OfferedProfile").Preload("AcceptedProfile").Find(&matches)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, matches)
}

func GetMatchByPersonID(c *gin.Context) {
	// Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var matches []Match
	results := db.Where("offered = ?", uint(id)).Or("accepted = ?", uint(id)).Preload("OfferedProfile").Preload("AcceptedProfile").Find(&matches)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, matches)

}

func PostMatch(c *gin.Context) {
	var newMatch Match

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newMatch); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newMatch.MatchID != 0 {
		for i, match := range Matches {
			if match.MatchID == newMatch.MatchID {
				// Update the existing match
				Matches[i].AcceptedTime = time.Now() // Update the AcceptedTime
				Matches[i].Offered = newMatch.Offered
				Matches[i].Accepted = newMatch.Accepted
				c.IndentedJSON(http.StatusOK, Matches)
				return
			}
		}
	}

	newMatch.MatchID = len(Matches) + 1000 // Assign a new ID based on the length of the slice
	newMatch.OfferedTime = time.Now()      // Set the OfferedTime to the current time
	Matches = append(Matches, newMatch)
	c.IndentedJSON(http.StatusCreated, newMatch)

}

func Signup(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind the JSON payload to the struct
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate the password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new user
	newUser := User{
		ID:           uint(len(users)),
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}
	users = append(users, newUser)

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   newUser.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"email": newUser.Email,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// Return token and user info (excluding password hash)
	resp := struct {
		Token string `json:"token"`
		ID    uint   `json:"id"`
	}{
		Token: tokenString,
		ID:    newUser.ID,
	}
	c.IndentedJSON(http.StatusCreated, resp)

}
