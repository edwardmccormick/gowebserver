package main

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

// Map to store connections by room ID
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomsLock sync.Mutex // Protect access to the rooms map

func WebsocketListener(c *gin.Context) {
    conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // Parse room ID from query parameters
    roomID := c.Query("id")
    if roomID == "" {
        conn.WriteMessage(websocket.TextMessage, []byte("Room ID is required"))
        fmt.Println("Room ID is missing in WebSocket request")
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
            delete(rooms, roomID)
        }
        roomsLock.Unlock()
    }()

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
        var receivedMessage struct {
            ID      int64  `json:"id"`
            Who     string `json:"who"`
            Message string `json:"message"`
        }
        err = json.Unmarshal(msg, &receivedMessage)
        if err != nil {
            fmt.Printf("Error unmarshaling message: %v\n", err)
            conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
            continue
        }

        // Prepare the response
        response := ChatMessage{
            ID:      int64(i),
            Time:    time.Now().Format("03:04:05 PM"),
            Who:     receivedMessage.Who,
            Message: receivedMessage.Message,
        }
        fmt.Printf("Broadcasting response: %+v\n", response)
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
    }
}