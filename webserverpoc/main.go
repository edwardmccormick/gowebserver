package gowebserver

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var db *gorm.DB
var mongoClient *mongo.Client

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default (for development)
	},
}

func DefaultBindAddress() string {
	if isRunningInDockerContainer() {
		return "0.0.0.0:8080"
	}
	return "localhost:8080"
}
