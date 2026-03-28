package gowebserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"github.com/gin-gonic/gin"
)

func ChatMessagesFromSQL(c *gin.Context) {
	var messages []ChatMessage
	if result := db.Find(&messages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, messages)
}

func GetChatMessagesForMatch(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	if chatMessageService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chat message service not available"})
		return
	}

	messages, totalCount, err := chatMessageService.GetPagedMessages(c.Request.Context(), uint(matchID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load chat history: %v", err)})
		return
	}

	if len(messages) == 0 {
		introMessage, err := chatMessageService.EnsureIntroduction(c.Request.Context(), uint(matchID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate introduction: %v", err)})
			return
		}
		if introMessage == nil {
			c.IndentedJSON(http.StatusOK, []ChatMessage{})
			return
		}
		c.IndentedJSON(http.StatusOK, []domain.ChatMessage{*introMessage})
		return
	}

	response := struct {
		Messages   []domain.ChatMessage `json:"messages"`
		Pagination struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
			Count  int `json:"count"`
		} `json:"pagination"`
	}{
		Messages: messages,
		Pagination: struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
			Count  int `json:"count"`
		}{
			Total:  totalCount,
			Offset: offset,
			Limit:  limit,
			Count:  len(messages),
		},
	}

	c.IndentedJSON(http.StatusOK, response)
}

func MarkMessagesAsRead(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	var match Match
	result := db.First(&match, matchID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if match.Offered != userID && match.Accepted != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this match"})
		return
	}

	if chatService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chat service not initialized"})
		return
	}

	if err := chatService.ResetUnreadCount(c.Request.Context(), uint(matchID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error resetting unread count: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}

func GenerateDateSuggestionForMatch(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	var match Match
	result := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if match.Offered != userID && match.Accepted != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this match"})
		return
	}

	if chatMessageService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chat message service not initialized"})
		return
	}

	dateMessage, err := chatMessageService.GenerateDateMessage(c.Request.Context(), uint(matchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate date suggestion: %v", err)})
		return
	}

	c.IndentedJSON(http.StatusOK, dateMessage)
}

func GenerateVibeChatForMatch(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	var match Match
	result := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if match.Offered != userID && match.Accepted != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this match"})
		return
	}

	if chatMessageService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chat message service not initialized"})
		return
	}

	vibeMessage, err := chatMessageService.GenerateVibeMessage(c.Request.Context(), uint(matchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate vibe chat: %v", err)})
		return
	}

	c.IndentedJSON(http.StatusOK, vibeMessage)
}

func ChatMessagesFromMongo(c *gin.Context) {
	str := c.Param("id")
	matchID, err := strconv.Atoi(str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	messages, _, err := chatMessageService.GetPagedMessages(c.Request.Context(), uint(matchID), limit, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error loading chat history: %v", err)})
		return
	}

	c.IndentedJSON(http.StatusOK, messages)
}
