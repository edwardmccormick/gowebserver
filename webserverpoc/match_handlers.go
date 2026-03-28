package gowebserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMatches(c *gin.Context) {
	if matchService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Match service not initialized"})
		return
	}

	matches, err := matchService.ListMatches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, matches)
}

func GetMatchByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if matchService == nil || profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Match/profile service not initialized"})
		return
	}

	matches, err := matchService.ListMatchesForPerson(uint(id))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(matches) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	profileService.HydrateMatchMedia(matches)
	c.IndentedJSON(http.StatusOK, matches)
}

func GetMatchByPersonID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if matchService == nil || profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Match/profile service not initialized"})
		return
	}

	matches, err := matchService.ListMatchesForPerson(uint(id))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(matches) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	profileService.HydrateMatchMedia(matches)
	c.IndentedJSON(http.StatusOK, matches)
}

func PostMatch(c *gin.Context) {
	var newMatch Match

	if err := c.BindJSON(&newMatch); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if matchService == nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Match service not initialized"})
		return
	}

	savedMatch, err := matchService.UpsertMatch(newMatch)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating match"})
		return
	}

	status := http.StatusCreated
	if newMatch.ID != 0 {
		status = http.StatusOK
	}
	c.IndentedJSON(status, savedMatch)
}
