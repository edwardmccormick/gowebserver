package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *gorm.DB
var mongoClient *mongo.Client

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default (for development)
	},
}

func main() {
	router := gin.Default()

	// Custom CORS configuration
	configCors := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5172"}, // Replace with your frontend's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router.Use(cors.New(configCors))

	// Load the configuration
	var config *Config
	var err error

	if isRunningInDockerContainer() {
		config, err = LoadConfig("./config.json") // Adjust the path as needed
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
	} else {
		config, err = LoadConfig("./configlocal.json") // Adjust the path as needed
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
	}

	// Connect to MySQL
	db, err = ConnectToMySQLWithConfig(config)
	if err != nil {
		fmt.Println("Error connecting to MySQL:", err)
		return
	}
	fmt.Println("Connected to MySQL.")
	db = db.Debug()

	// Perform GORM automigration
	if err := db.AutoMigrate(
		&User{},
		&Person{},
		&Details{},
		&Match{},
		&ChatMessage{},
	); err != nil {
		fmt.Errorf("Error during automigration: %v", err)
	}
	fmt.Println("Database schema migrated successfully.")

	// Connect to MongoDB
	mongoClient, err := ConnectToMongoDBWithConfig(config)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	fmt.Println("Connected to MongoDB.")

	// Populate the database
	if err := PopulateDatabase(db, mongoClient); err != nil {
		fmt.Println("Error populating database:", err)
		return
	}
	fmt.Println(mongoClient)

	fmt.Println("Database check and population complete.")

	fmt.Println("Let's do the thing")
	// Sort by general functionality - signup, auth, login, logout
	router.POST("/signup", Signup)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.POST("/people", PostPeople) // Really create profile for yourself but this logic made sense to me
	router.GET("/users", GetUsers)     // for troubleshooting or admin purposes

	// Search and find people and profile information
	router.GET("/people", GetPeople)
	router.POST("/peoplelocation", GetPeopleByLocation)
	router.GET("/people/:id", GetPeopleByID)
	router.GET("/photos/:id", GetPhotosByID)

	// Matchmaking and chat functionality
	router.GET("/matches", GetMatches)
	router.GET("/matches/:id", GetMatchByPersonID)
	router.POST("/matches", PostMatch)
	router.GET("/ws", WebsocketListener)

	// Troubleshooting and utility endpoints
	router.GET("/favicon.ico", GetFaviconIco)
	router.GET("/", GreetUser)
	router.GET("/greet/:name", GreetUserByName)
	router.GET("/chat", ChatMessagesFromSQL)
	router.GET("/chat/:id", ChatMessagesFromMongo)
	if isRunningInDockerContainer() {
		router.Run("0.0.0.0:8080")
	} else {
		router.Run("localhost:8080")
	}

}
