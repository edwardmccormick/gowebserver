package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
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

// sortMessagesByTime sorts messages by their timestamp in ascending order
func sortMessagesByTime(messages []ChatMessage) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})
}

// Map to store active connections by user ID
var activeConnections = make(map[uint]*websocket.Conn)
var activeConnectionsLock sync.Mutex // Protect access to the activeConnections map

// broadcastMessageToMatch sends a message to all connected clients in a match room
func broadcastMessageToMatch(matchID uint, message *ChatMessage) {
	// Convert matchID to string for the room map
	roomID := strconv.FormatUint(uint64(matchID), 10)
	
	// Serialize the message to JSON
	msgBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshaling message for broadcast: %v\n", err)
		return
	}
	
	// Lock the rooms map to safely access it
	roomsLock.Lock()
	defer roomsLock.Unlock()
	
	// Get all connections in this room
	connections := rooms[roomID]
	if connections == nil {
		fmt.Printf("No active connections for match %d\n", matchID)
		return
	}
	
	// Send the message to all connected clients in this room
	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			fmt.Printf("Error sending message to client: %v\n", err)
			// Note: We don't remove the connection here as it might be a temporary issue
			// The connection will be properly cleaned up when it's closed
		}
	}
	
	fmt.Printf("Broadcast message to %d clients in match %d\n", len(connections), matchID)
}

// addMessageToMongo adds a single message to an existing conversation in MongoDB
func addMessageToMongo(matchID uint, message ChatMessage) error {
	// Check if MongoDB client is initialized
	if mongoClient == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}
	
	// Use a consistent database name across the application
	collection := mongoClient.Database("urmid").Collection("chathistory")
	
	// Create a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// First, check if a conversation for this match exists
	var conversation Conversation
	err := collection.FindOne(ctx, bson.M{"match_id": matchID}).Decode(&conversation)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no conversation exists, create a new one with this message
			newConversation := Conversation{
				MatchID:  matchID,
				Messages: []ChatMessage{message},
			}
			
			_, err = collection.InsertOne(ctx, newConversation)
			if err != nil {
				return fmt.Errorf("failed to create new conversation: %v", err)
			}
			
			fmt.Printf("Created new conversation for match %d with message\n", matchID)
			return nil
		}
		
		return fmt.Errorf("error finding conversation: %v", err)
	}
	
	// If conversation exists, add the new message to it
	conversation.Messages = append(conversation.Messages, message)
	
	// Update the conversation in MongoDB
	_, err = collection.ReplaceOne(
		ctx,
		bson.M{"match_id": matchID},
		conversation,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update conversation: %v", err)
	}
	
	fmt.Printf("Added message to conversation for match %d\n", matchID)
	return nil
}

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
	
	// Load chat history from MongoDB
	messages, err := loadChatHistoryFromMongo(matchID)
	if err != nil {
		// Log error but continue without history
		fmt.Printf("Error loading chat history in WebsocketListener: %v\n", err)
		sessionChat.Messages = []ChatMessage{}
	} else {
		sessionChat.Messages = messages
		
		// Check if we have any messages
		if len(messages) == 0 {
			// No messages yet, try to generate an AI introduction
			fmt.Printf("No chat history for match %d, generating AI introduction\n", matchID)
			
			// Generate AI introduction
			introMessage, err := GenerateMatchIntroduction(matchID)
			if err != nil {
				fmt.Printf("Failed to generate introduction: %v\n", err)
			} else {
				// Set the 'who' field to 0, which represents the AI system in our application
				introMessage.Who = 0 // Use 0 to represent the AI/system
				
				// Add to session chat
				sessionChat.Messages = append(sessionChat.Messages, *introMessage)
				
				// Save to MongoDB
				tempChat := Conversation{
					MatchID:  matchID,
					Messages: []ChatMessage{*introMessage},
				}
				err = dumpChatHistoryToMongo(tempChat)
				if err != nil {
					fmt.Printf("Failed to save introduction message: %v\n", err)
				}
				
				// Send to client
				msgBytes, err := json.Marshal(introMessage)
				if err == nil {
					conn.WriteMessage(websocket.TextMessage, msgBytes)
				}
			}
		} else {
			// We have existing messages, send them to the client
			// Limit to the most recent messages if there are too many
			start := 0
			if len(messages) > 50 {
				start = len(messages) - 50
			}
			
			// Send history to the client
			for _, msg := range messages[start:] {
				msgBytes, err := json.Marshal(msg)
				if err == nil {
					conn.WriteMessage(websocket.TextMessage, msgBytes)
				}
			}
		}
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



	defer func() {
		fmt.Printf("WebSocket connection closing for user %d in room %s\n", userID, roomID)
		
		// Remove connection from the room when it disconnects
		roomsLock.Lock()
		delete(rooms[roomID], conn)
		roomsLock.Unlock()

		// Remove the user's active connection
		activeConnectionsLock.Lock()
		delete(activeConnections, uint(userID))
		activeConnectionsLock.Unlock()

		// Only save chat history if there are messages
		if len(sessionChat.Messages) > 0 {
			// Dump chat history to MongoDB
			fmt.Printf("Saving %d chat messages for match %d\n", len(sessionChat.Messages), sessionChat.MatchID)
			err := dumpChatHistoryToMongo(sessionChat)
			if err != nil {
				fmt.Printf("Failed to save chat history: %v\n", err)
			}
		}
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

		// Set the current time for this message
		now := time.Now()
		receivedMessage.Time = now
		receivedMessage.CreatedAt = now
		receivedMessage.UpdatedAt = now
		
		// Add match ID if not set
		if receivedMessage.MatchID == 0 {
			receivedMessage.MatchID = int(matchID)
		}
		
		// Add message to chat history
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

// dumpChatHistoryToMongo saves the chat history to MongoDB
// This function handles updating existing conversations or creating new ones
func dumpChatHistoryToMongo(sessionChat Conversation) error {
	// Check if MongoDB client is initialized
	if mongoClient == nil {
		error := "Error: mongoClient is not initialized"
		fmt.Println(error)
		return fmt.Errorf(error)
	}

	// Use a consistent database name across the application
	collection := mongoClient.Database("urmid").Collection("chathistory")
	
	// Create a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if there is an existing conversation for this match
	var currentHistory Conversation
	err := collection.FindOne(ctx, bson.M{"match_id": sessionChat.MatchID}).Decode(&currentHistory)

	if err != nil {
		// If no documents found, create a new conversation
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No chat history found for match %d, creating a new conversation\n", sessionChat.MatchID)
			
			// Add timestamp to each message if not already present
			now := time.Now()
			for i := range sessionChat.Messages {
				if sessionChat.Messages[i].CreatedAt.IsZero() {
					sessionChat.Messages[i].CreatedAt = now
				}
			}
			
			// Insert the new conversation
			_, err := collection.InsertOne(ctx, sessionChat)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to create new chat history: %v", err)
				fmt.Println(errMsg)
				return fmt.Errorf(errMsg)
			}
			fmt.Printf("Successfully created new chat history for match %d with %d messages\n", 
				sessionChat.MatchID, len(sessionChat.Messages))
			return nil
		}

		// Other database errors
		errMsg := fmt.Sprintf("Error querying MongoDB for match %d: %v", sessionChat.MatchID, err)
		fmt.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Update existing conversation
	_, err = collection.ReplaceOne(ctx, bson.M{"match_id": sessionChat.MatchID}, sessionChat)
	if err != nil {
		errMsg := fmt.Sprintf("Error updating chat history for match %d: %v", sessionChat.MatchID, err)
		fmt.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	fmt.Printf("Successfully updated chat history for match %d with %d messages\n", 
		sessionChat.MatchID, len(sessionChat.Messages))
	return nil
}

// loadChatHistoryFromMongo retrieves chat message history for a specific match from MongoDB
// Returns the messages and any error that occurred during the operation
func loadChatHistoryFromMongo(matchID uint) ([]ChatMessage, error) {
	var history Conversation

	// Check if MongoDB client is initialized
	if mongoClient == nil {
		errMsg := "Error: mongoClient is not initialized"
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Use a consistent database name across the application
	collection := mongoClient.Database("urmid").Collection("chathistory")
	
	// Create a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Find the conversation document for this match ID
	err := collection.FindOne(ctx, bson.M{"match_id": matchID}).Decode(&history)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No history found is a normal situation for new conversations
			fmt.Printf("No chat history found for match %d\n", matchID)
			return []ChatMessage{}, nil
		}
		
		// Other database errors should be reported
		errMsg := fmt.Sprintf("Error loading chat history for match %d: %v", matchID, err)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	
	// Sort messages by time to ensure they appear in chronological order
	if len(history.Messages) > 0 {
		// Simple in-memory sort by time
		sortMessagesByTime(history.Messages)
		fmt.Printf("Loaded %d messages for match %d from MongoDB\n", len(history.Messages), matchID)
	}
	
	return history.Messages, nil
}
