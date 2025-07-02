package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

// Map to store connections by room ID
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomsLock sync.Mutex // Protect access to the rooms map
// Map to store chat history by MatchID
var chatHistory = make(map[int][]ChatMessage)

// Map to store active connections by user ID
var activeConnections = make(map[uint]*websocket.Conn)
var activeConnectionsLock sync.Mutex // Protect access to the activeConnections map

func WebsocketListener(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Parse room ID and user ID from query parameters
	roomID := c.Query("id")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	if roomID == "" || userID == 0 {
		conn.WriteMessage(websocket.TextMessage, []byte("Room ID and User ID are required"))
		return
	}

	matchID, err := strconv.Atoi(roomID)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid room ID"))
		return
	}

	// Add connection to the room
	roomsLock.Lock()
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rooms[roomID][conn] = true
	roomsLock.Unlock()

	// Track the user's active connection
	activeConnectionsLock.Lock()
	activeConnections[uint(userID)] = conn
	activeConnectionsLock.Unlock()

	// var config *Config

	// if isRunningInDockerContainer() {
	// 	config, err = LoadConfig("./config.json") // Adjust the path as needed
	// 	if err != nil {
	// 		fmt.Println("Error loading config:", err)
	// 		return
	// 	}
	// } else {
	// 	config, err = LoadConfig("./configlocal.json") // Adjust the path as needed
	// 	if err != nil {
	// 		fmt.Println("Error loading config:", err)
	// 		return
	// 	}
	// }

	// // Connect to MongoDB
	// mongoClient, err := ConnectToMongoDBWithConfig(config)
	// if err != nil {
	// 	fmt.Println("Error connecting to MongoDB:", err)
	// 	return
	// }
	// fmt.Println("Connected to MongoDB.")

	defer func() {
		// Remove connection from the room when it disconnects
		roomsLock.Lock()
		delete(rooms[roomID], conn)
		roomsLock.Unlock()

		// Remove the user's active connection
		activeConnectionsLock.Lock()
		delete(activeConnections, uint(userID))
		activeConnectionsLock.Unlock()

		// Dump chat history to MongoDB
		dumpChatHistoryToMongo(matchID)
	}()

	// Deliver undelivered messages
	deliverUndeliveredMessages(conn, matchID, uint(userID))

	// Load chat history from MongoDB
	loadChatHistoryFromMongo(conn, matchID)

	// Chat message array for this session
	var sessionChat []ChatMessage

	i := 0
	for {
		// Read message from the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}

		// Parse the received message as JSON
		var receivedMessage ChatMessage
		err = json.Unmarshal(msg, &receivedMessage)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		// Prepare the response
		response := ChatMessage{
			ID:      int64(i),
			Time:    time.Now(),
			Who:     uint(userID),
			Message: receivedMessage.Message,
		}
		sessionChat = append(sessionChat, response)
		i++

		// Broadcast the message to all connections in the room except the sender
		roomsLock.Lock()
		for client := range rooms[roomID] {
			if client != conn {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					client.Close()
					delete(rooms[roomID], client)
				}
			}
		}
		roomsLock.Unlock()

		// Handle undelivered messages
		handleUndeliveredMessages(matchID, response)
	}

	// Save chat history for the session
	roomsLock.Lock()
	chatHistory[matchID] = append(chatHistory[matchID], sessionChat...)
	roomsLock.Unlock()
}

func deliverUndeliveredMessages(conn *websocket.Conn, matchID int, userID uint) {
	var undeliveredMessages []ChatMessage
	results := db.Where("match_id = ? AND who != ?", matchID, userID).Find(&undeliveredMessages)
	if results.RowsAffected != 0 {
		for _, msg := range undeliveredMessages {
			msgBytes, _ := json.Marshal(msg)
			conn.WriteMessage(websocket.TextMessage, msgBytes)
		}

		// Delete undelivered messages after sending
		db.Where("match_id = ? AND who != ?", matchID, userID).Delete(&ChatMessage{})
	}
}

func handleUndeliveredMessages(matchID int, message ChatMessage) {
	for _, match := range Matches {
		if match.ID == uint(matchID) {
			activeConnectionsLock.Lock()
			_, connected := activeConnections[match.Accepted]
			activeConnectionsLock.Unlock()

			if !connected {
				// Save the message to MySQL
				message.MatchID = matchID
				db.Create(&message)
			}
			break
		}
	}
}

func dumpChatHistoryToMongo(matchID int) {
	history := chatHistory[matchID]
	if len(history) == 0 {
		return
	}

	collection := mongoClient.Database("gowebserver").Collection("chathistory")
	_, err := collection.InsertOne(context.TODO(), bson.M{
		"match_id":  matchID,
		"history":   history,
		"timestamp": time.Now(),
	})
	if err != nil {
		fmt.Printf("Error dumping chat history to MongoDB: %v\n", err)
		return
	}

	// Clear in-memory chat history
	roomsLock.Lock()
	delete(chatHistory, matchID)
	roomsLock.Unlock()
}

func loadChatHistoryFromMongo(conn *websocket.Conn, matchID int) {
	collection := mongoClient.Database("gowebserver").Collection("chathistory")

	var result struct {
		History []ChatMessage `bson:"history"`
	}
	err := collection.FindOne(context.TODO(), bson.M{"match_id": matchID}).Decode(&result)
	if err != nil {
		fmt.Printf("Error loading chat history from MongoDB: %v\n", err)
		return
	}

	// Send the last 30 messages to the user
	start := 0
	if len(result.History) > 30 {
		start = len(result.History) - 30
	}
	for _, msg := range result.History[start:] {
		msgBytes, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}
