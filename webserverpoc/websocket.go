package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Map to store connections by room ID
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomsLock sync.Mutex // Protect access to the rooms map
// Map to store chat history by MatchID
var chatHistory = make(map[int][]ChatMessage)

func WebsocketListener(c *gin.Context) {
    conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // Parse room ID from query parameters
    roomID := c.Query("id")
	userID, _ := strconv.Atoi(c.Query("user_id")) // Assume user_id is passed as a query parameter
    if roomID == "" {
        conn.WriteMessage(websocket.TextMessage, []byte("Room ID is required"))
        fmt.Println("Room ID is missing in WebSocket request")
        return
    }

    // Parse MatchID from room ID
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

    defer func() {
        // Remove connection from the room when it disconnects
        roomsLock.Lock()
        delete(rooms[roomID], conn)
        if len(rooms[roomID]) == 0 {
            // Save chat history when both users disconnect
            if history, exists := chatHistory[matchID]; exists {
                fmt.Printf("Saving chat history for MatchID %d: %+v\n", matchID, history)
            }
        }
        roomsLock.Unlock()
    }()

    // Deliver undelivered messages
    deliverUndeliveredMessages(conn, matchID, uint(userID))

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

        // Log the received message
        fmt.Printf("Received message: %s\n", string(msg))

        // Parse the received message as JSON
        var receivedMessage ChatMessage
        err = json.Unmarshal(msg, &receivedMessage)
        if err != nil {
            fmt.Printf("Error unmarshaling message: %v\n", err)
            conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
            continue
        }

        // Prepare the response
        response := ChatMessage{
            ID:      int64(i),
            Time:    time.Now(),
            Who:     receivedMessage.Who,
            Message: receivedMessage.Message,
        }
        sessionChat = append(sessionChat, response)
        i++

        // Marshal the response to JSON
        msgBytes, err := json.Marshal(response)
        if err != nil {
            fmt.Printf("Error marshaling response: %v\n", err)
            conn.WriteMessage(websocket.TextMessage, []byte("Error processing message"))
            continue
        }

        // Broadcast the message to all connections in the room except the sender
        roomsLock.Lock()
        for client := range rooms[roomID] {
            if client != conn {
                err := client.WriteMessage(websocket.TextMessage, msgBytes)
                if err != nil {
                    fmt.Printf("Error broadcasting message: %v\n", err)
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
    for _, match := range Matches {
        if match.MatchID == matchID {
            // Check if the user is the Offered or Accepted user
            var undeliveredMessages []ChatMessage
            if match.Offered == userID {
                undeliveredMessages = match.OfferedChat
            } else if match.Accepted == userID {
                undeliveredMessages = match.AcceptedChat
            }

            // Deliver only undelivered messages
            for _, msg := range undeliveredMessages {
                msgBytes, _ := json.Marshal(msg)
                conn.WriteMessage(websocket.TextMessage, msgBytes)
            }

            // Clear undelivered messages after sending
            if match.Offered == userID {
                match.OfferedChat = nil
            } else if match.Accepted == userID {
                match.AcceptedChat = nil
            }
            break
        }
    }
}

func handleUndeliveredMessages(matchID int, message ChatMessage) {
    for i, match := range Matches {
        if match.MatchID == matchID {
            // Check if the sender is the Offered or Accepted user
            if match.Offered == uint(message.ID) {
                Matches[i].OfferedChat = append(Matches[i].OfferedChat, message)
            } else if match.Accepted == uint(message.ID) {
                Matches[i].AcceptedChat = append(Matches[i].AcceptedChat, message)
            }
            break
        }
    }
}