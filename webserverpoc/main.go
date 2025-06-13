package main

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
	"github.com/asmarques/geodist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type person struct {
	ID           int         `json:"id"`
	Name         string         `json:"name"`
	Motto        string         `json:"motto"`
	LatLocation  float64        `json:"lat"`
	LongLocation float64        `json:"long"`
	Profile      string         `json:"profile"`
	Verbiage     map[string]int `json:"verbiage"` // Add this to store rankings for each category
}

type user struct {
	ID       int `json:"id"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
}

type processedProfile struct {
	ID       int         `json:"id"`
	Name     string         `json:"name"`
	Motto    string         `json:"motto"`
	Distance float64        `json:"distance"`
	Profile  string         `json:"profile"`
	Verbiage map[string]int `json:"verbiage"`
}

type profilePhoto struct {
	Url     string `json:"url"`
	Caption string `json:"caption"`
}

type chatMessage struct {
	ID      int64  `json:"id"`
	Who     string `json:"who"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default (for development)
		return true
	},
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
	router.POST("/peoplelocation", getPeopleByLocation)
	router.GET("/people/:id", getPeopleByID)
	router.POST("/people", postPeople)
	router.GET("/", greetUser)
	router.GET("/greet/:name", greetUserByName)
	router.GET("/login", login)
	router.POST("login", login)
	router.GET("/photos/:id", getPhotosByID)
	router.GET("/ws", websocketDummy)
	router.GET("/logout", signout)
	router.Run("localhost:8080")
}

// albums slice to seed record album data.
var people = []person{
	{ID: 0, Name: "Bobby", Motto: "Always Ready", LatLocation: 29.534261019806404, LongLocation: 98.47049550692051, Profile: "https://picsum.photos/200/200"},
	{ID: 1, Name: "Joe", Motto: "Always Faithful", LatLocation: 29.52016959410149, LongLocation: 98.49401109752402, Profile: "https://picsum.photos/150/300"},
	{ID: 2, Name: "Fred", Motto: "Always Prepared", LatLocation: 29.453596593823395, LongLocation: 98.47166788793534, Profile: "https://picsum.photos/250/250"},
	{ID: 3, Name: "Turd Furguson", Motto: "I'm supre close let's party!", LatLocation: 29.419922273698763, LongLocation: 98.48366872664229, Profile: "https://picsum.photos/250/250"},
	{ID: 4, Name: "Don", Motto: "I want you Bigly", LatLocation: 38.898, LongLocation: 77.037, Profile: "https://www.whitehouse.gov/wp-content/uploads/2025/06/President-Donald-Trump-Official-Presidential-Portrait.png"},
	{ID: 5, Name: "The Founder", Motto: "I just want this fucking thing to work", LatLocation: 48.858, LongLocation: -2.294, Profile: "https://ted.mccormickhub.com/img/tedProfilePicture.jpg"},
}

var users = []user{
	{ID: 5, Email: "ted@urmid.com", PasswordHash: hashPassword("password123")},
}

var locationOrigin = geodist.Point{Lat: 29.42618, Long: 98.48618} // The Stupid Alamo
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

func dummyMessages(num int) chatMessage {
	var message = chatMessage{ID: int64(num), Who: isme[num%len(isme)], Message: messages[num%len(messages)]}
	return message
}

var jwtSecret = []byte("supersecretkey") // Use a secure random key in production!
var refreshTokens = make(map[string]string) // Map user ID to refresh token

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
			Profile:  p.Profile,
		})

		sort.Slice(processedPeople, func(i, j int) bool {
			return processedPeople[i].Distance < processedPeople[j].Distance
		})
	}

	// Return the processedPeople array as JSON
	c.IndentedJSON(http.StatusOK, processedPeople)
}

func getPeopleByLocation(c *gin.Context) {
	var processedPeople []processedProfile
	jwtMiddleware(c)
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
		processedPeople = append(processedPeople, processedProfile{
			ID:       p.ID,
			Name:     p.Name,
			Motto:    p.Motto,
			Distance: roundedDistance,
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

func getPhotosByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(id)
	if id == 0 || id == 1 || id == 2 || id == 3 || id == 4 || id == 5 || id == 6 || id == 7 || id == 8 || id == 9 || id == 10 {
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
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Find the user by email
    var user *user
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
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
        "email": user.Email,
    })
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
        return
    }

    // Return token and user info (excluding password hash)
    resp := struct {
        Token  string  `json:"token"`
        Person person  `json:"person"`
    }{
        Token:  tokenString,
        Person: people[user.ID],
    }

    c.JSON(http.StatusOK, resp)
}

func hashPassword(pw string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    return string(hash)
}

func signout(c *gin.Context) {
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
    delete(refreshTokens, userID)
    c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func jwtMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
        c.Abort()
        return
    }

    // Validate the token
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !parsedToken.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        c.Abort()
        return
    }

    c.Next()
}


func websocketDummy(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		return
	}
	defer conn.Close()

	i := 0
	for {
		i++
		// Marshal messages to JSON
		msgBytes, err := json.Marshal(dummyMessages(i))
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Error marshaling messages"))
			return
		}
		conn.WriteMessage(websocket.TextMessage, msgBytes)
		time.Sleep(time.Second * 10)
	}
}

func websocketListener(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	i := 0
	for {
		// Read message from the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}

		// Log the received message
		fmt.Printf("Received message: %s\n", string(msg))

		// Echo the message back to the client
		response := chatMessage{
			ID:      int64(i),
			Who:     "Me", // You can customize this based on the sender
			Message: string(msg),
		}
		i++

		// Marshal the response to JSON
		msgBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("Error marshaling response: %v\n", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error processing message"))
			continue
		}

		// Send the response back to the client
		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			fmt.Printf("Error writing message: %v\n", err)
			break
		}
	}
}
