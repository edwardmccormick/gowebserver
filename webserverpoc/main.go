package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default (for development)
	},
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	fmt.Println("Let's do the thing")
	router.GET("/people", GetPeople)
	router.POST("/peoplelocation", GetPeopleByLocation)
	router.GET("/people/:id", GetPeopleByID)
	router.POST("/people", PostPeople)
	router.GET("/", GreetUser)
	router.GET("/greet/:name", GreetUserByName)
	router.POST("/login", Login)
	router.GET("/logout", Signout)
	router.GET("/ws", WebsocketListener)
	router.GET("/photos/:id", GetPhotosByID)
	router.GET("/favicon.ico", GetFaviconIco)
	router.GET("/matches", GetMatches)
	router.GET("/matches/:id", GetMatchByPersonID)
	router.POST("/matches", PostMatch)
	router.POST("/signup", Signup)

	router.Run("localhost:8080")
}
