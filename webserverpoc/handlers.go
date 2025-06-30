package main

import (
	// "encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Sort by general functionality - signup, auth, login, logout
func Signup(c *gin.Context) {
	// Signup both creates the user and gives them a JWT for the session so they can create their person/profile
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

func Logout(c *gin.Context) {
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

func PostPeople(c *gin.Context) {
	var newPerson Person

	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	result := db.Create(newPerson) // pass a slice to insert multiple row
	fmt.Println("Created rows: ", result.RowsAffected)
	c.IndentedJSON(http.StatusCreated, newPerson)
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

// Search and find people and profile information

func GetPeople(c *gin.Context) {
	var people []Person
	if result := db.Find(&people); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, people)
}

func GetPeopleByLocation(c *gin.Context) {
	JwtMiddleware(c)
	var req struct {
		Lat   float64 `json:"lat"`
		Long  float64 `json:"long"`
		Range float64 `json:"range"`
	}
	fmt.Println(c)
	fmt.Println(req)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Earth's radius in kilometers
	const earthRadius = 6371.0

	// Calculate the bounding box
	latRange := req.Range / earthRadius * (180 / math.Pi)
	longRange := req.Range / (earthRadius * math.Cos(req.Lat*math.Pi/180)) * (180 / math.Pi)

	minLat := req.Lat - latRange
	maxLat := req.Lat + latRange
	minLong := req.Long - longRange
	maxLong := req.Long + longRange

	fmt.Printf("Bounding box: minLat=%f, maxLat=%f, minLong=%f, maxLong=%f\n", minLat, maxLat, minLong, maxLong)

	// Query the database for people within the bounding box
	var people []Person
	if err := db.Where("lat_location BETWEEN ? AND ? AND long_location BETWEEN ? AND ?", minLat, maxLat, minLong, maxLong).Find(&people).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}

	// Return the filtered people
	c.IndentedJSON(http.StatusOK, people)
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

// Matchmaking and chat functionality

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

	// Call BindJSON to bind the received JSON to Match struct
	if err := c.BindJSON(&newMatch); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newMatch.MatchID != 0 {
		results := db.Where("offered = ?", newMatch.MatchID).Preload("OfferedProfile").Preload("AcceptedProfile").First(&newMatch)
		if results.Error != nil {
			fmt.Println(results.Error)
		}
		if results.RowsAffected == 0 {
			fmt.Printf("No match found with ID %d, checking if match exists between these two people", newMatch.MatchID)
			results := db.Where("offered = ? AND accepted = ?", newMatch.Offered, newMatch.Accepted).Or("offered = ? AND accepted = ?", newMatch.Accepted, newMatch.Offered).Preload("OfferedProfile").Preload("AcceptedProfile").First(&newMatch)
			if results.Error != nil {
				fmt.Println(results.Error)
			}
			if results.RowsAffected == 0 {
				fmt.Printf("No match found between these two people, creating a new match")
				results := db.Create(newMatch) // pass a slice to insert multiple rows
				fmt.Println("Created rows: ", results.RowsAffected)
				if results.Error != nil {
					fmt.Println(results.Error)
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating match, existing match found and unable to create a new match"})
					return
				}
				c.IndentedJSON(http.StatusCreated, newMatch)
				return
			} else {
				fmt.Printf("Match found between these two people, although the match ID was wrong, updating the existing match")

			}
		}
		// Update the existing match
		newMatch.AcceptedTime = time.Now() // Update the AcceptedTime

		c.IndentedJSON(http.StatusOK, newMatch)
		return
	}

	// If MatchID is 0, create a new match
	results := db.Create(newMatch) // pass a slice to insert multiple row
	fmt.Println("Created rows: ", results.RowsAffected)
	if results.Error != nil {
		fmt.Println(results.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating match"})
		return
	}
	c.IndentedJSON(http.StatusCreated, newMatch)

}

// Troubleshooting and utility endpoints
func GetFaviconIco(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "https://urmid.com/favicon.ico")
}

func GreetUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

func GreetUserByName(c *gin.Context) {
	name := c.Param("name")

	c.String(http.StatusOK, "Hello %s", name)
}
