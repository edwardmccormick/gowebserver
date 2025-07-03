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
	"go.mongodb.org/mongo-driver/mongo"
)

// Map to store connections by room ID
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomsLock sync.Mutex // Protect access to the rooms map

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

	matchIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid room ID"))
		return
	}
	matchID := uint(matchIDInt)

	// Chat message array for this session
	var sessionChat Conversation
	sessionChat.MatchID = matchID
	if mongoClient == nil {
        fmt.Println("Error: mongoClient is not initialized from websocket function")
	}
	sessionChat.Messages = loadChatHistoryFromMongo(matchID)

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
		if mongoClient == nil {
        fmt.Println("Error: mongoClient is not initialized from websocket function")
	}
		dumpChatHistoryToMongo(sessionChat)
	}()


	
	// // Load chat history from MongoDB
	// start := 0
	// if len(result.History) > 30 {
	// 	start = len(result.History) - 30
	// }
	// for _, msg := range result.History[start:] {
	// 	msgBytes, _ := json.Marshal(msg)
	// 	conn.WriteMessage(websocket.TextMessage, msgBytes)
	// }


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

		receivedMessage.Time = time.Now()
		// Prepare the message to chat history
		sessionChat.Messages = append(sessionChat.Messages, receivedMessage)

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

	}

}

// func deliverUndeliveredMessages(conn *websocket.Conn, matchID int, userID uint) {
// 	var undeliveredMessages []ChatMessage
// 	results := db.Where("match_id = ? AND who != ?", matchID, userID).Find(&undeliveredMessages)
// 	if results.RowsAffected != 0 {
// 		for _, msg := range undeliveredMessages {
// 			msgBytes, _ := json.Marshal(msg)
// 			conn.WriteMessage(websocket.TextMessage, msgBytes)
// 		}

// 		// Delete undelivered messages after sending
// 		db.Where("match_id = ? AND who != ?", matchID, userID).Delete(&ChatMessage{})
// 	}
// }

// func handleUndeliveredMessages(matchID int, message ChatMessage) {
// 	for _, match := range Matches {
// 		if match.ID == uint(matchID) {
// 			activeConnectionsLock.Lock()
// 			_, connected := activeConnections[match.Accepted]
// 			activeConnectionsLock.Unlock()

// 			if !connected {
// 				// Save the message to MySQL
// 				message.MatchID = matchID
// 				db.Create(&message)
// 			}
// 			break
// 		}
// 	}
// }

func dumpChatHistoryToMongo(sessionChat Conversation) {
	if mongoClient == nil {
        fmt.Println("Error: mongoClient is not initialized")
		return
	}

	collection := mongoClient.Database("urmid").Collection("chathistory")

	var currentHistory Conversation

	err := collection.FindOne(context.TODO(), bson.M{"match_id": sessionChat.MatchID}).Decode(&currentHistory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No chat history found in MongoDB, creating one")
			_, err := collection.InsertOne(context.TODO(), sessionChat)
			if err != nil {
				panic(err)
			}
			return
		}
		panic(err)
	}
	_, err = collection.ReplaceOne(context.TODO(), bson.M{"match_id": sessionChat.MatchID}, sessionChat)
	if err != nil {
		fmt.Printf("Error dumping chat history to MongoDB: %v\n", err)
		return
	} else {
		fmt.Printf("Successfully dumped chat history to MongoDB: %v\n", sessionChat)
	}
}

func loadChatHistoryFromMongo( matchID uint) []ChatMessage {
	var History Conversation
	if mongoClient == nil {
        fmt.Println("Error: mongoClient is not initialized")
		return History.Messages
	}

	collection := mongoClient.Database("urmid").Collection("chathistory")
	
	err := collection.FindOne(context.TODO(), bson.M{"match_id": matchID}).Decode(&History)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No chat history found in MongoDB")
			return History.Messages
		}
		panic(err)
	}
	
	return History.Messages

}
