package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminGetAllUsers returns all users in the system (for admin use only)
func AdminGetAllUsers(c *gin.Context) {
	var users []User
	if result := db.Preload("Person").Preload("Person.Photos").Preload("Person.Profile").Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// AdminGetAllChats returns all chat messages in the system (for admin use only)
func AdminGetAllChats(c *gin.Context) {
	var messages []ChatMessage
	if result := db.Order("time desc").Find(&messages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// AdminGetRecentUpdates returns recently updated profiles and photos (for admin use only)
func AdminGetRecentUpdates(c *gin.Context) {
	// Get recently updated profiles
	var recentProfiles []Person
	if result := db.Order("updated_at desc").Limit(20).Preload("Photos").Preload("Profile").Find(&recentProfiles); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent profiles: " + result.Error.Error()})
		return
	}

	// Get recently uploaded photos
	var recentPhotos []ProfilePhoto
	if result := db.Order("created_at desc").Limit(20).Find(&recentPhotos); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent photos: " + result.Error.Error()})
		return
	}

	// Return both sets of data
	response := gin.H{
		"recent_profiles": recentProfiles,
		"recent_photos":   recentPhotos,
	}

	c.JSON(http.StatusOK, response)
}

// AdminSetUserAdmin allows an admin to promote another user to admin status
func AdminSetUserAdmin(c *gin.Context) {
	var req struct {
		UserID  uint `json:"user_id"`
		IsAdmin bool `json:"is_admin"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Find the user
	var user User
	if result := db.First(&user, req.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update admin status
	user.IsAdmin = req.IsAdmin
	if result := db.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User admin status updated successfully", "user": user})
}