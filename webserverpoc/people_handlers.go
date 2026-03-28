package gowebserver

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func PostPeople(c *gin.Context) {
	var newPerson Person
	if err := c.BindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile service not initialized"})
		return
	}

	savedPerson, err := profileService.SaveProfile(newPerson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uploadURLs, profileUploadURLs, err := profileService.GenerateUploadTargets(savedPerson.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profileService.AttachPhotoViewURLs(savedPerson.Photos, 180*time.Minute)

	resp := struct {
		Person            Person         `json:"person"`
		UploadUrls        []ProfilePhoto `json:"upload_urls,omitempty"`
		ProfileUploadUrls []ProfilePhoto `json:"profile_upload_urls,omitempty"`
	}{
		Person:            savedPerson,
		UploadUrls:        uploadURLs,
		ProfileUploadUrls: profileUploadURLs,
	}

	c.IndentedJSON(http.StatusCreated, resp)
}

func GetUsers(c *gin.Context) {
	if profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile service not initialized"})
		return
	}

	users, err := profileService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

type PersonWithDistance struct {
	Person
	Distance float64 `json:"distance"`
}

func GetPeople(c *gin.Context) {
	var people []Person
	var userLat, userLong float64
	var hasLocation bool
	var err error

	token := c.GetHeader("Authorization")
	if token == "" {
		if profileService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile service not initialized"})
			return
		}
		people, err = profileService.ListPeople()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hasLocation = false
	} else {
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return signingKey(), nil
		})
		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
			return
		}

		currentUser, err := profileService.GetPersonByID(uint(userID))
		if err != nil {
			people, err = profileService.ListPeople()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			hasLocation = false
		} else {
			userLat = currentUser.LatLocation
			userLong = currentUser.LongLocation
			hasLocation = userLat != 0 && userLong != 0

			people, err = profileService.ListPeopleNearLocation(userLat, userLong, uint(userID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if hasLocation {
		var peopleWithDistance []PersonWithDistance
		for _, person := range people {
			if person.ID != 0 && person.LatLocation != 0 && person.LongLocation != 0 {
				distance := CalculateHaversineDistance(userLat, userLong, person.LatLocation, person.LongLocation)
				peopleWithDistance = append(peopleWithDistance, PersonWithDistance{
					Person:   person,
					Distance: math.Round(distance*10) / 10,
				})
			}
		}
		c.IndentedJSON(http.StatusOK, peopleWithDistance)
	} else {
		c.IndentedJSON(http.StatusOK, people)
	}
}

func GetPeopleByLocation(c *gin.Context) {
	JwtMiddleware(c)
	var req struct {
		Lat   float64 `json:"lat"`
		Long  float64 `json:"long"`
		Range float64 `json:"range"`
	}
	fmt.Println(c)
	fmt.Println(req)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	const earthRadius = 6371.0
	latRange := req.Range / earthRadius * (180 / math.Pi)
	longRange := req.Range / (earthRadius * math.Cos(req.Lat*math.Pi/180)) * (180 / math.Pi)

	minLat := req.Lat - latRange
	maxLat := req.Lat + latRange
	minLong := req.Long - longRange
	maxLong := req.Long + longRange

	if profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile service not initialized"})
		return
	}

	people, err := profileService.ListPeopleWithinBounds(minLat, maxLat, minLong, maxLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}

	c.IndentedJSON(http.StatusOK, people)
}

func GetPeopleByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile service not initialized"})
		return
	}

	person, err := profileService.GetPersonByID(uint(id))
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, person)
}

func GetPhotosByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if profileService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile service not initialized"})
		return
	}

	profilePhotos, err := profileService.ListPhotosByPersonID(uint(id))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(profilePhotos) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found for that id"})
		return
	}

	profileService.AttachPhotoManagementURLs(profilePhotos, 60*time.Minute)
	c.IndentedJSON(http.StatusOK, profilePhotos)
}
