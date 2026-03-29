package gowebserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/realtime"
	"github.com/edwardmccormick/gowebserver/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// sortMessagesByTime sorts messages by their timestamp in ascending order
func sortMessagesByTime(messages []ChatMessage) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})
}

var realtimeHub = realtime.NewHub()
var notificationCenter = realtime.NewNotificationCenter()

// updateUnreadCounts increments unread message counters and sends notifications
func updateUnreadCounts(matchID uint, senderID uint) {
	fmt.Printf("updateUnreadCounts: Processing unread count for match %d, sender %d\n", matchID, senderID)
	if chatService == nil {
		fmt.Printf("Error finding match %d: chat service is not initialized\n", matchID)
		return
	}
	if err := chatService.UpdateUnreadCounts(context.Background(), matchID, senderID); err != nil {
		fmt.Printf("Error updating match %d unread counts: %v\n", matchID, err)
	}
}

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

	connections := realtimeHub.RoomConnections(roomID)
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

	// Update unread counters if the message is from a user (not system)
	if message.Who != 0 {
		updateUnreadCounts(matchID, message.Who)
	}

	fmt.Printf("Broadcast message to %d clients in match %d\n", len(connections), matchID)
}

func saveChatMessage(matchID uint, message ChatMessage) error {
	if chatService == nil {
		return fmt.Errorf("chat service is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := chatService.AppendMessage(ctx, matchID, ChatMessage{
		ID:          message.ID,
		MatchID:     message.MatchID,
		Time:        message.Time,
		Who:         message.Who,
		MessageType: message.MessageType,
		Message:     message.Message,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to persist message: %v", err)
	}

	fmt.Printf("Saved message for match %d\n", matchID)
	return nil
}

// SSEHandler creates a Server-Sent Events stream for real-time notifications
func SSEHandler(c *gin.Context) {
	fmt.Println("SSEHandler: Starting SSE connection handler")
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify JWT authentication
	tokenUserIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Authentication required"})
		fmt.Println("SSEHandler: Authentication failed - no userID in context")
		return
	}

	tokenUserID := tokenUserIDValue.(uint)
	if uint(userID) != tokenUserID {
		c.JSON(403, gin.H{"error": "Not authorized to subscribe to this user's notifications"})
		fmt.Printf("SSEHandler: Authorization failed - token userID %d doesn't match requested userID %d\n", tokenUserID, userID)
		return
	}

	fmt.Printf("SSEHandler: Starting SSE stream for user %d\n", userID)

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Create notification channel for this client
	notificationChan := make(chan realtime.NotificationEvent, 10) // Buffer for 10 notifications

	// Register client in the notification center
	notificationCenter.Register(userIDStr, notificationChan)

	// Make sure we unregister when client disconnects
	defer notificationCenter.Unregister(userIDStr, notificationChan)

	// Create a context that will be canceled when the client disconnects
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// Monitor the connection status
	go func() {
		<-ctx.Done()
		// Connection closed
		fmt.Printf("SSE connection closed for user %s\n", userIDStr)
	}()

	// Stream events until client disconnects
	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-notificationChan:
			if !ok {
				// Channel was closed
				fmt.Printf("SSEHandler: Channel closed for user %d\n", userID)
				return false
			}

			// Convert notification to JSON
			eventData, err := json.Marshal(event)
			if err != nil {
				fmt.Printf("Error marshaling SSE event: %v\n", err)
				return true // Continue streaming despite error
			}

			// Send the event
			c.SSEvent("message", string(eventData))
			fmt.Printf("SSEHandler: Sent event to user %d: %s\n", userID, string(eventData))
			return true

		case <-time.After(30 * time.Second):
			// Send a heartbeat/keep-alive message
			c.SSEvent("ping", time.Now().Unix())
			fmt.Printf("SSEHandler: Sent ping to user %d\n", userID)
			return true
		}
	})
}

// resetUnreadCount resets the unread message count for a user in a match
func resetUnreadCount(matchID uint, userID uint) {
	if chatService == nil {
		fmt.Printf("Error resetting unread count for match %d: chat service is not initialized\n", matchID)
		return
	}

	if err := chatService.MarkConversationRead(context.Background(), matchID, userID); err != nil {
		fmt.Printf("Error resetting unread count for match %d: %v\n", matchID, err)
	}
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

	messages, err := loadChatHistory(matchID)
	if err != nil {
		fmt.Printf("Error loading chat history in WebsocketListener: %v\n", err)
		messages = []ChatMessage{}
	}

	if len(messages) == 0 && chatMessageService != nil {
		introMessage, err := chatMessageService.EnsureIntroduction(c.Request.Context(), matchID)
		if err != nil {
			fmt.Printf("Failed to generate introduction: %v\n", err)
		} else if introMessage != nil {
			messages = append(messages, ChatMessage(*introMessage))
		}
	}

	start := 0
	if len(messages) > 50 {
		start = len(messages) - 50
	}
	for _, existing := range messages[start:] {
		msgBytes, err := json.Marshal(existing)
		if err == nil {
			conn.WriteMessage(websocket.TextMessage, msgBytes)
		}
	}

	realtimeHub.AddRoomConnection(roomID, conn)
	realtimeHub.SetActiveConnection(uint(userID), conn)

	// Reset unread message count when user opens the chat
	resetUnreadCount(matchID, uint(userID))

	defer func() {
		fmt.Printf("WebSocket connection closing for user %d in room %s\n", userID, roomID)
		realtimeHub.RemoveRoomConnection(roomID, conn)
		realtimeHub.RemoveActiveConnection(uint(userID))
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
		if receivedMessage.MessageType == "" {
			receivedMessage.MessageType = "user"
		}

		// Add match ID if not set
		if receivedMessage.MatchID == 0 {
			receivedMessage.MatchID = int(matchID)
		}

		if err := saveChatMessage(matchID, receivedMessage); err != nil {
			fmt.Printf("Error saving chat message: %v\n", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to save message"))
			continue
		}

		broadcastMessageToMatch(matchID, &receivedMessage)
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
//
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
//
// 			if !connected {
// 				// Save the message to MySQL
// 				message.MatchID = matchID
// 				db.Create(&message)
// 			}
// 			break
// 		}
// 	}
// }

func dumpChatHistoryToMongo(sessionChat Conversation) error {
	if chatService == nil {
		error := "Error: chat service is not initialized"
		fmt.Println(error)
		return fmt.Errorf("%s", error)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conversation := Conversation{
		MatchID:  sessionChat.MatchID,
		Messages: make([]ChatMessage, 0, len(sessionChat.Messages)),
	}
	for _, msg := range sessionChat.Messages {
		conversation.Messages = append(conversation.Messages, ChatMessage{
			ID:          msg.ID,
			MatchID:     msg.MatchID,
			Time:        msg.Time,
			Who:         msg.Who,
			MessageType: msg.MessageType,
			Message:     msg.Message,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		})
	}

	if err := chatService.SaveConversation(ctx, conversation); err != nil {
		errMsg := fmt.Sprintf("Error updating chat history for match %d: %v", sessionChat.MatchID, err)
		fmt.Println(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	fmt.Printf("Successfully updated chat history for match %d with %d messages\n",
		sessionChat.MatchID, len(sessionChat.Messages))
	return nil
}

func loadChatHistory(matchID uint) ([]ChatMessage, error) {
	var history Conversation

	if chatService == nil {
		errMsg := "Error: chat service is not initialized"
		fmt.Println(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conversation, err := chatService.LoadConversation(ctx, matchID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			fmt.Printf("No chat history found for match %d\n", matchID)
			return []ChatMessage{}, nil
		}

		errMsg := fmt.Sprintf("Error loading chat history for match %d: %v", matchID, err)
		fmt.Println(errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}

	history.MatchID = conversation.MatchID
	history.Messages = make([]ChatMessage, 0, len(conversation.Messages))
	for _, msg := range conversation.Messages {
		history.Messages = append(history.Messages, ChatMessage{
			ID:          msg.ID,
			MatchID:     msg.MatchID,
			Time:        msg.Time,
			Who:         msg.Who,
			MessageType: msg.MessageType,
			Message:     msg.Message,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		})
	}

	// Sort messages by time to ensure they appear in chronological order
	if len(history.Messages) > 0 {
		// Simple in-memory sort by time
		sortMessagesByTime(history.Messages)
		fmt.Printf("Loaded %d messages for match %d from store\n", len(history.Messages), matchID)
	}

	return history.Messages, nil
}
