package main

import (
	// "encoding/json"
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// Sort by general functionality - signup, auth, login, logout
func Signup(c *gin.Context) {
	// Signup both creates the user and gives them a JWT for the session so they can create their person/profile
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind the JSON payload to the struct
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate the password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new user
	newUser := User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		LastLogin:    time.Now(),
	}
	results := db.Create(&newUser)
	if results.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": results.Error.Error()})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   newUser.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"email": newUser.Email,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// Generate presigned upload URLs
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")
	var uploadUrls []ProfilePhoto
	var profileUploadUrls []ProfilePhoto

	// Generate dedicated presigned URLs for profile photos (both raw and cropped)
	profileKey := fmt.Sprintf("%d/profile", newUser.ID)

	// Cropped profile photo URL
	profileReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileUrl, err := profileReq.Presign(180 * time.Minute)
	if err == nil {
		profileUploadUrls = append(profileUploadUrls, ProfilePhoto{
			Url:      profileUrl,
			S3Key:    profileKey,
			PersonID: newUser.ID,
		})
	}

	// Generate presigned URLs for regular photos
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d/image%d", newUser.ID, i+1)
		req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		url, err := req.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
		if err != nil {
			fmt.Printf("Failed to generate presigned URL for key %s: %v\n", key, err)
			continue
		}
		uploadUrls = append(uploadUrls, ProfilePhoto{
			Url:      url,
			S3Key:    key,
			PersonID: newUser.ID,
		})
	}

	// Return token and user info (excluding password hash)
	resp := struct {
		Token             string         `json:"token"`
		ID                uint           `json:"id"`
		UploadUrls        []ProfilePhoto `json:"upload_urls,omitempty"`         // Optional field for upload URLs
		ProfileUploadUrls []ProfilePhoto `json:"profile_upload_urls,omitempty"` // Optional field for profile upload URLs
	}{
		Token:             tokenString,
		ID:                newUser.ID,
		UploadUrls:        uploadUrls,
		ProfileUploadUrls: profileUploadUrls,
	}
	c.IndentedJSON(http.StatusCreated, resp)

}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Find the user by email
	var user *User
	db.Where("email = ?", req.Email).First(&user)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"email": user.Email,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	var person = Person{ID: user.ID}
	db.Preload("Photos").Preload("Profile").First(&person)

	// Generate presigned upload URLs
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")
	var uploadUrls []ProfilePhoto
	var profileUploadUrls []ProfilePhoto

	// Generate dedicated presigned URLs for profile photos (both raw and cropped)
	profileKey := fmt.Sprintf("%d/profile", user.ID)

	// Cropped profile photo URL
	profileReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileUrl, err := profileReq.Presign(180 * time.Minute)
	if err == nil {
		profileUploadUrls = append(profileUploadUrls, ProfilePhoto{
			Url:      profileUrl,
			S3Key:    profileKey,
			PersonID: user.ID,
		})
	}

	// Raw profile photo URL
	// rawProfileReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
	// 	Bucket: aws.String(bucketName),
	// 	Key:    aws.String(rawProfileKey),
	// })
	// rawProfileUrl, err := rawProfileReq.Presign(180 * time.Minute)
	// if err == nil {
	// 	uploadUrls = append(uploadUrls, ProfilePhoto{
	// 		Url:      rawProfileUrl,
	// 		S3Key:    rawProfileKey,
	// 		PersonID: user.ID,
	// 	})
	// }

	// Generate presigned URLs for regular photos
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d/image%d", user.ID, i+1)
		req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		url, err := req.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
		if err != nil {
			fmt.Printf("Failed to generate presigned URL for key %s: %v\n", key, err)
			continue
		}
		uploadUrls = append(uploadUrls, ProfilePhoto{
			Url:      url,
			S3Key:    key,
			PersonID: user.ID,
		})
	}

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	// Generate presigned URLs for all photos
	for j := range person.Photos {
		wg.Add(1)
		go func(photo *ProfilePhoto) {
			defer wg.Done()
			// Generate presigned URL
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(photo.S3Key),
			})
			url, err := req.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
			if err != nil {
				fmt.Printf("Failed to generate presigned URL for S3Key %s: %v\n", photo.S3Key, err)
				return
			}
			photo.Url = url // Update the Url field with the presigned URL
		}(&person.Photos[j])
	}

	// // Generate presigned URLs for viewing profile photos
	// profileKey := fmt.Sprintf("%d/profile", user.ID)
	// rawProfileKey := fmt.Sprintf("%d/rawprofile", user.ID)
	// These values already are declared above!

	// Cropped profile photo URL for viewing
	profileGetReq, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileViewUrl, err := profileGetReq.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
	if err == nil {
		// Check if profile photo exists in the database
		if person.Profile.ID == 0 {
			// Create a new profile photo entry
			person.Profile = ProfilePhoto{
				PersonID: user.ID,
				S3Key:    profileKey,
				Url:      profileViewUrl,
			}
		} else {
			// Update existing profile photo URL
			person.Profile.Url = profileViewUrl
		}
	}

	// // Raw profile photo URL for viewing
	// rawProfileGetReq, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
	// 	Bucket: aws.String(bucketName),
	// 	Key:    aws.String(rawProfileKey),
	// })
	// rawProfileViewUrl, err := rawProfileGetReq.Presign(180 * time.Minute)
	// if err == nil {
	// 	// Find or create the raw profile photo entry
	// 	var rawProfileFound bool
	// 	for i := range person.Photos {
	// 		if person.Photos[i].S3Key == rawProfileKey {
	// 			person.Photos[i].Url = rawProfileViewUrl
	// 			rawProfileFound = true
	// 			break
	// 		}
	// 	}

	// 	if !rawProfileFound {
	// 		// Add raw profile photo to photos array
	// 		person.Photos = append(person.Photos, ProfilePhoto{
	// 			PersonID: user.ID,
	// 			S3Key:    rawProfileKey,
	// 			Url:      rawProfileViewUrl,
	// 			Caption:  "Raw Profile Photo",
	// 		})
	// 	}
	// }

	// Wait for all goroutines to finish
	wg.Wait()

	// Return token and user info (excluding password hash)
	resp := struct {
		Token             string         `json:"token"`
		Person            Person         `json:"person"`
		UploadUrls        []ProfilePhoto `json:"upload_urls,omitempty"`         // Optional field for upload URLs
		ProfileUploadUrls []ProfilePhoto `json:"profile_upload_urls,omitempty"` // Optional field for profile upload URLs
	}{
		Token:             tokenString,
		Person:            person,
		UploadUrls:        uploadUrls,
		ProfileUploadUrls: profileUploadUrls,
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	// Parse and validate the token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// claims := parsedToken.Claims.(jwt.MapClaims)

	// Remove the refresh token
	ExpiredTokens = append(ExpiredTokens, token)
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func PostPeople(c *gin.Context) {
	var newPerson Person

	// Bind the JSON payload to the Person struct
	if err := c.BindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Save the Person record first
	result := db.Omit("Photos", "Profile").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&newPerson)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Handle associated photos
	// First, get existing photos for this person
	var existingPhotos []ProfilePhoto
	if err := db.Where("person_id = ?", newPerson.ID).Find(&existingPhotos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch existing photos: %v", err)})
		return
	}

	// Create a map of existing S3 keys for quick lookup
	existingS3Keys := make(map[string]uint)
	for _, photo := range existingPhotos {
		existingS3Keys[photo.S3Key] = photo.ID
	}

	// Track which photos are in the update request
	requestedPhotoKeys := make(map[string]bool)
	
	// Process each photo in the request
	for i := range newPerson.Photos {
		photo := &newPerson.Photos[i]
		photo.PersonID = newPerson.ID // Ensure the foreign key is set correctly
		
		// Mark this photo as requested to keep
		requestedPhotoKeys[photo.S3Key] = true

		// Check if this S3 key already exists
		if existingID, exists := existingS3Keys[photo.S3Key]; exists {
			// Update existing photo instead of creating a new one
			photo.ID = existingID
			if err := db.Model(&ProfilePhoto{}).Where("id = ?", existingID).Updates(map[string]interface{}{
				"caption": photo.Caption,
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update photo: %v", err)})
				return
			}
		} else {
			// This is a new photo, create it
			photo.ID = 0 // Reset the ID to let the database auto-generate it
			if err := db.Create(photo).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save photo: %v", err)})
				return
			}
		}
	}
	
	// Delete photos that are no longer included in the request
	// This ensures we don't keep accumulating photos when updating the profile
	for _, existingPhoto := range existingPhotos {
		if !requestedPhotoKeys[existingPhoto.S3Key] {
			// Skip deleting profile photos (they're handled separately)
			if strings.HasSuffix(existingPhoto.S3Key, "/profile") {
				continue
			}
			
			// This photo is no longer needed, delete it
			if err := db.Delete(&ProfilePhoto{}, existingPhoto.ID).Error; err != nil {
				fmt.Printf("Failed to delete unused photo %d: %v\n", existingPhoto.ID, err)
				// Continue without returning error - this is not critical
			}
		}
	}

	// Generate presigned upload URLs
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")
	var uploadUrls []ProfilePhoto
	var profileUploadUrls []ProfilePhoto

	// Generate dedicated presigned URL for profile photo
	profileKey := fmt.Sprintf("%d/profile", newPerson.ID)

	// Cropped profile photo URL
	profileReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileUrl, err := profileReq.Presign(180 * time.Minute)
	if err == nil {
		profileUploadUrls = append(profileUploadUrls, ProfilePhoto{
			Url:      profileUrl,
			S3Key:    profileKey,
			PersonID: newPerson.ID,
		})
	}

	// Generate presigned URLs for regular photos
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d/image%d", newPerson.ID, i+1)
		req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		url, err := req.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
		if err != nil {
			fmt.Printf("Failed to generate presigned URL for key %s: %v\n", key, err)
			continue
		}
		uploadUrls = append(uploadUrls, ProfilePhoto{
			Url:      url,
			S3Key:    key,
			PersonID: newPerson.ID,
		})
	}

	// Generate presigned URLs for viewing all photos
	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	// Generate presigned URLs for all photos
	for j := range newPerson.Photos {
		wg.Add(1)
		go func(photo *ProfilePhoto) {
			defer wg.Done()
			// Generate presigned URL
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(photo.S3Key),
			})
			url, err := req.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
			if err != nil {
				fmt.Printf("Failed to generate presigned URL for S3Key %s: %v\n", photo.S3Key, err)
				return
			}
			photo.Url = url // Update the Url field with the presigned URL
		}(&newPerson.Photos[j])
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Return the created Person record with presigned URLs
	resp := struct {
		Person            Person         `json:"person"`
		UploadUrls        []ProfilePhoto `json:"upload_urls,omitempty"`         // Optional field for upload URLs
		ProfileUploadUrls []ProfilePhoto `json:"profile_upload_urls,omitempty"` // Optional field for profile upload URLs
	}{
		Person:            newPerson,
		UploadUrls:        uploadUrls,
		ProfileUploadUrls: profileUploadUrls,
	}

	c.IndentedJSON(http.StatusCreated, resp)
}

func GetUsers(c *gin.Context) {
	var users []User
	if result := db.Preload("Person").Preload("Person.Photos").Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, users)
}

// Search and find people and profile information

// PersonWithDistance extends Person to include distance information
type PersonWithDistance struct {
	Person
	Distance float64 `json:"distance"` // Distance in miles
}

func GetPeople(c *gin.Context) {
	// Initialize variables we'll use in both branches
	var people []Person
	var userLat, userLong float64
	var hasLocation bool
	
	// Extract the JWT token from the Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		// If no token is provided, return all people without distance filtering
		if result := db.Preload("Photos").Preload("Profile").Find(&people); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		hasLocation = false // No user location available
	} else {
		// Parse and validate the token
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract user ID from token claims
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
			return
		}

		// Get the user ID from the token
		userID, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in token"})
			return
		}

		// Get the current user's location from the database
		var currentUser Person
		result := db.First(&currentUser, uint(userID))
		if result.Error != nil {
			// If we can't get the user, fall back to returning all people without distance filtering
			if result := db.Preload("Photos").Preload("Profile").Find(&people); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			hasLocation = false
		} else {
			// Now we have the user's location, let's apply the distance filter
			userLat = currentUser.LatLocation
			userLong = currentUser.LongLocation
			hasLocation = true

			// Default radius is 50 miles
			const defaultRadius = 50.0

			// Calculate the bounding box for faster initial filtering
			minLat, maxLat, minLong, maxLong := CalculateBoundingBox(userLat, userLong, defaultRadius)

			// Get people within the bounding box
			if result := db.Where("lat_location BETWEEN ? AND ? AND long_location BETWEEN ? AND ? AND id != ?", 
				minLat, maxLat, minLong, maxLong, uint(userID)).Preload("Photos").Preload("Profile").Find(&people); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
		}
	}

	// Load AWS credentials from environment variables
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup
	for i := range people {
		// Find the profile photo in Photos
		var profilePhoto *ProfilePhoto
		for j := range people[i].Photos {
			if people[i].Photos[j].S3Key == fmt.Sprintf("%d/profile", people[i].ID) {
				profilePhoto = &people[i].Photos[j]
				break // Stop after finding the first match
			}
		}

		// If found, generate the signed URL and update Profile
		if profilePhoto != nil {
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(profilePhoto.S3Key),
			})
			url, err := req.Presign(180 * time.Minute)
			if err == nil {
				people[i].Profile = *profilePhoto // Copy all fields
				people[i].Profile.Url = url       // Update the URL
			}
		}

		// Now filter out profile photos from Photos
		var filteredPhotos []ProfilePhoto
		for _, photo := range people[i].Photos {
			if !strings.Contains(photo.S3Key, "/profile") {
				filteredPhotos = append(filteredPhotos, photo)
			}
		}
		people[i].Photos = filteredPhotos

		// Generate presigned URLs for regular photos
		for j := range people[i].Photos {
			wg.Add(1)
			go func(photo *ProfilePhoto) {
				defer wg.Done()
				// Generate presigned URL
				req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(photo.S3Key),
				})
				url, err := req.Presign(60 * time.Minute) // Presigned URL valid for 60 minutes
				if err != nil {
					fmt.Printf("Failed to generate presigned URL for S3Key %s: %v\n", photo.S3Key, err)
					return
				}
				photo.Url = url // Update the Url field with the presigned URL
			}(&people[i].Photos[j])
		}

		// Generate presigned URL for profile photo
		// if people[i].Profile.S3Key != "" {
		// 	wg.Add(1)
		// 	go func(profile *ProfilePhoto) {
		// 		defer wg.Done()
		// 		// Generate presigned URL
		// 		req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		// 			Bucket: aws.String(bucketName),
		// 			Key:    aws.String(profile.S3Key),
		// 		})
		// 		url, err := req.Presign(180 * time.Minute) // Presigned URL valid for 180 minutes
		// 		if err != nil {
		// 			fmt.Printf("Failed to generate presigned URL for profile S3Key %s: %v\n", profile.S3Key, err)
		// 			return
		// 		}
		// 		profile.Url = url // Update the Url field with the presigned URL
		// 	}(&people[i].Profile)
		// }
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// If we have user location, calculate distances and filter by 50 miles
	if hasLocation {
		// Create a new slice for people with distances
		peopleWithDistance := make([]PersonWithDistance, 0, len(people))
		
		// Default radius is 50 miles
		const defaultRadius = 50.0
		
		// Calculate exact distance for each person and filter
		for _, person := range people {
			// Skip the current user
			distance := CalculateHaversineDistance(userLat, userLong, person.LatLocation, person.LongLocation)
			
			// Only include people within the exact distance limit
			if distance <= defaultRadius {
				peopleWithDistance = append(peopleWithDistance, PersonWithDistance{
					Person:   person,
					Distance: math.Round(distance*10) / 10, // Round to 1 decimal place
				})
			}
		}
		
		c.IndentedJSON(http.StatusOK, peopleWithDistance)
	} else {
		// If no user location, return people without distance
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

	// Earth's radius in kilometers
	const earthRadius = 6371.0

	// Calculate the bounding box
	latRange := req.Range / earthRadius * (180 / math.Pi)
	longRange := req.Range / (earthRadius * math.Cos(req.Lat*math.Pi/180)) * (180 / math.Pi)

	minLat := req.Lat - latRange
	maxLat := req.Lat + latRange
	minLong := req.Long - longRange
	maxLong := req.Long + longRange

	fmt.Printf("Bounding box: minLat=%f, maxLat=%f, minLong=%f, maxLong=%f\n", minLat, maxLat, minLong, maxLong)

	// Query the database for people within the bounding box
	var people []Person
	if err := db.Where("lat_location BETWEEN ? AND ? AND long_location BETWEEN ? AND ?", minLat, maxLat, minLong, maxLong).Find(&people).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}

	// Return the filtered people
	c.IndentedJSON(http.StatusOK, people)
}

func GetPeopleByID(c *gin.Context) {
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(id)

	var person = Person{ID: uint(id)}
	results := db.Preload("Photos").Preload("Profile").First(&person)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
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
	fmt.Println(id)

	var profilePhotos = []ProfilePhoto{}
	results := db.Where("person_id = ?", uint(id)).Find(&profilePhotos)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found for that id"})
		return
	}

		// Load AWS credentials from environment variables
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")

	// Use a WaitGroup to wait for all goroutines to finish
    var wg sync.WaitGroup
    for j := range profilePhotos {
        wg.Add(1)
        go func(photo *ProfilePhoto) {
            defer wg.Done()

            // Generate presigned URL for viewing
            getReq, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
                Bucket: aws.String(bucketName),
                Key:    aws.String(photo.S3Key),
            })
            getUrl, err := getReq.Presign(60 * time.Minute) // Presigned URL valid for 60 minutes
            if err != nil {
                fmt.Printf("Failed to generate presigned GET URL for S3Key %s: %v\n", photo.S3Key, err)
                return
            }
            photo.Url = getUrl

            // Generate presigned URL for uploading
            putReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
                Bucket: aws.String(bucketName),
                Key:    aws.String(photo.S3Key),
            })
            putUrl, err := putReq.Presign(60 * time.Minute) // Presigned URL valid for 60 minutes
            if err != nil {
                fmt.Printf("Failed to generate presigned PUT URL for S3Key %s: %v\n", photo.S3Key, err)
                return
            }
            photo.Upload = putUrl

            // Generate presigned URL for deleting
            deleteReq, _ := s3Client.DeleteObjectRequest(&s3.DeleteObjectInput{
                Bucket: aws.String(bucketName),
                Key:    aws.String(photo.S3Key),
            })
            deleteUrl, err := deleteReq.Presign(60 * time.Minute) // Presigned URL valid for 60 minutes
            if err != nil {
                fmt.Printf("Failed to generate presigned DELETE URL for S3Key %s: %v\n", photo.S3Key, err)
                return
            }
            photo.Delete = deleteUrl
        }(&profilePhotos[j])
    }

    // Wait for all goroutines to finish
    wg.Wait()


	c.IndentedJSON(http.StatusOK, profilePhotos)
}

// Matchmaking and chat functionality

func GetMatches(c *gin.Context) {
	var matches []Match
	if result := db.Preload("OfferedProfile").Preload("AcceptedProfile").Find(&matches); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, matches)
}

func GetMatchByID(c *gin.Context) {
	// Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var matches []Match
	// Preload the profiles with their photos
	results := db.Where("offered = ?", uint(id)).Or("accepted = ?", uint(id)).
		Preload("OfferedProfile.Photos").Preload("OfferedProfile.Profile").
		Preload("AcceptedProfile.Photos").Preload("AcceptedProfile.Profile").
		Find(&matches)
	
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	// Load AWS credentials from environment variables
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	// Process each match to add presigned URLs
	for i := range matches {
		// Process OfferedProfile
		processProfilePhotos(&matches[i].OfferedProfile, s3Client, bucketName, &wg)
		
		// Process AcceptedProfile
		processProfilePhotos(&matches[i].AcceptedProfile, s3Client, bucketName, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	c.IndentedJSON(http.StatusOK, matches)
}

// Helper function to process profile photos and generate presigned URLs
func processProfilePhotos(person *Person, s3Client *s3.S3, bucketName string, wg *sync.WaitGroup) {
	if person == nil || person.ID == 0 {
		return
	}

	// Find the profile photo in Photos
	var profilePhoto *ProfilePhoto
	for j := range person.Photos {
		if person.Photos[j].S3Key == fmt.Sprintf("%d/profile", person.ID) {
			profilePhoto = &person.Photos[j]
			break // Stop after finding the first match
		}
	}

	// If found, generate the signed URL and update Profile
	if profilePhoto != nil {
		req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(profilePhoto.S3Key),
		})
		url, err := req.Presign(180 * time.Minute)
		if err == nil {
			person.Profile = *profilePhoto // Copy all fields
			person.Profile.Url = url       // Update the URL
		}
	}

	// Now filter out profile photos from Photos
	var filteredPhotos []ProfilePhoto
	for _, photo := range person.Photos {
		if !strings.Contains(photo.S3Key, "/profile") {
			filteredPhotos = append(filteredPhotos, photo)
		}
	}
	person.Photos = filteredPhotos

	// Generate presigned URLs for regular photos
	for j := range person.Photos {
		wg.Add(1)
		go func(photo *ProfilePhoto) {
			defer wg.Done()
			// Generate presigned URL
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(photo.S3Key),
			})
			url, err := req.Presign(60 * time.Minute) // Presigned URL valid for 60 minutes
			if err != nil {
				fmt.Printf("Failed to generate presigned URL for S3Key %s: %v\n", photo.S3Key, err)
				return
			}
			photo.Url = url // Update the Url field with the presigned URL
		}(&person.Photos[j])
	}
}

func GetMatchByPersonID(c *gin.Context) {
	// Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var matches []Match
	// Preload the profiles with their photos
	results := db.Where("offered = ?", uint(id)).Or("accepted = ?", uint(id)).
		Preload("OfferedProfile.Photos").Preload("OfferedProfile.Profile").
		Preload("AcceptedProfile.Photos").Preload("AcceptedProfile.Profile").
		Find(&matches)
	
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	// Load AWS credentials from environment variables
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	// Process each match to add presigned URLs
	for i := range matches {
		// Process OfferedProfile
		processProfilePhotos(&matches[i].OfferedProfile, s3Client, bucketName, &wg)
		
		// Process AcceptedProfile
		processProfilePhotos(&matches[i].AcceptedProfile, s3Client, bucketName, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	c.IndentedJSON(http.StatusOK, matches)
}

func PostMatch(c *gin.Context) {
	var newMatch Match

	// Call BindJSON to bind the received JSON to Match struct
	if err := c.BindJSON(&newMatch); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newMatch.ID != 0 {
		results := db.Where("offered = ?", newMatch.Offered).Preload("OfferedProfile").Preload("AcceptedProfile").First(&newMatch)
		if results.Error != nil {
			fmt.Println(results.Error)
		}
		if results.RowsAffected == 0 {
			fmt.Printf("No match found with ID %d, checking if match exists between these two people", newMatch.ID)
			results := db.Where("offered = ? AND accepted = ?", newMatch.Offered, newMatch.Accepted).Or("offered = ? AND accepted = ?", newMatch.Accepted, newMatch.Offered).Preload("OfferedProfile").Preload("AcceptedProfile").First(&newMatch)
			if results.Error != nil {
				fmt.Println(results.Error)
			}
			if results.RowsAffected == 0 {
				fmt.Printf("No match found between these two people, creating a new match")
				newMatch.OfferedTime = time.Now()
				newMatch.AcceptedTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
				results := db.Create(newMatch) // pass a slice to insert multiple rows
				fmt.Println("Created rows: ", results.RowsAffected)
				if results.Error != nil {
					fmt.Println(results.Error)
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating match, existing match found and unable to create a new match"})
					return
				}
				newMatch.AcceptedTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
				c.IndentedJSON(http.StatusCreated, newMatch)
				return
			} else {
				fmt.Printf("Match found between these two people, although the match ID was wrong, updating the existing match")

			}
		}
		// Update the existing match
		newMatch.AcceptedTime = time.Now() // Update the AcceptedTime
		results = db.Model(&newMatch).Where("id = ?", newMatch.ID).Save(&newMatch)
		fmt.Println("Created rows: ", results.RowsAffected)
		if results.Error != nil {
			fmt.Println(results.Error)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating match, existing match found and unable to create a new match"})
			return
		}
		c.IndentedJSON(http.StatusOK, newMatch)
		return
	}

	// If Match.ID is 0, create a new match
	// newMatch.MatchID = nil;
	newMatch.OfferedTime = time.Now()
	newMatch.AcceptedTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	results := db.Create(&newMatch) // pass a slice to insert multiple row
	fmt.Println("Created rows: ", results.RowsAffected)
	if results.Error != nil {
		fmt.Println(results.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating match"})
		return
	}
	c.IndentedJSON(http.StatusCreated, newMatch)

}

// Troubleshooting and utility endpoints
func GetFaviconIco(c *gin.Context) {
	c.File("./urmid.svg") // Serve the favicon.ico file from the current directory
}

func GreetUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

func GreetUserByName(c *gin.Context) {
	name := c.Param("name")

	c.String(http.StatusOK, "Hello %s", name)
}

func ChatMessagesFromSQL(c *gin.Context) {
	var messages []ChatMessage
	if result := db.Find(&messages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, messages)
}

// Get chat messages for a specific match
func GetChatMessagesForMatch(c *gin.Context) {
	// Extract the match ID from the URL parameter
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Optional limit parameter for pagination
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50 // Default to 50 messages
	}

	// Optional offset parameter for pagination
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Check if MongoDB client is initialized
	if mongoClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB connection not available"})
		return
	}

	// Get chat history from MongoDB
	messages, err := loadChatHistoryFromMongo(uint(matchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load chat history: %v", err)})
		return
	}

	// If no messages exist, generate an introduction message
	if len(messages) == 0 {
		// Generate AI introduction
		introMessage, err := GenerateMatchIntroduction(uint(matchID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate introduction: %v", err)})
			return
		}
		
		// Add introduction message to MongoDB
		var sessionChat Conversation
		sessionChat.MatchID = uint(matchID)
		sessionChat.Messages = []ChatMessage{*introMessage}
		
		err = dumpChatHistoryToMongo(sessionChat)
		if err != nil {
			fmt.Printf("Failed to save introduction message: %v\n", err)
			// Continue anyway as we can still return the message
		}
		
		// Return the introduction message
		c.IndentedJSON(http.StatusOK, []ChatMessage{*introMessage})
		return
	}

	// Sort messages by time to ensure chronological order
	sortMessagesByTime(messages)

	// Calculate total message count for pagination info
	totalCount := len(messages)
	
	// Apply offset and limit for pagination
	start := offset
	end := offset + limit
	
	// Ensure start and end are within bounds
	if start >= totalCount {
		start = 0
		end = 0
	} else if end > totalCount {
		end = totalCount
	}
	
	// Slice messages based on pagination parameters
	pagedMessages := messages
	if start < totalCount {
		pagedMessages = messages[start:end]
	} else {
		pagedMessages = []ChatMessage{}
	}

	// Return the messages with pagination metadata
	response := struct {
		Messages   []ChatMessage `json:"messages"`
		Pagination struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
			Count  int `json:"count"`
		} `json:"pagination"`
	}{
		Messages: pagedMessages,
		Pagination: struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
			Count  int `json:"count"`
		}{
			Total:  totalCount,
			Offset: offset,
			Limit:  limit,
			Count:  len(pagedMessages),
		},
	}

	c.IndentedJSON(http.StatusOK, response)
}

// MarkMessagesAsRead marks all messages in a match as read for the current user
func MarkMessagesAsRead(c *gin.Context) {
	// Extract the match ID from the URL parameter
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Get current user ID from JWT token
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	// Get the match from the database
	var match Match
	result := db.First(&match, matchID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Verify that the user is part of this match
	if match.Offered != userID && match.Accepted != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this match"})
		return
	}

	// Reset unread count based on which user it is
	if userID == match.Offered {
		match.UnreadOffered = 0
	} else if userID == match.Accepted {
		match.UnreadAccepted = 0
	}

	// Save the match with updated counters
	result = db.Save(&match)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error resetting unread count: %v", result.Error)})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}

// GenerateDateSuggestionForMatch creates and sends a personalized date recommendation
func GenerateDateSuggestionForMatch(c *gin.Context) {
	// Extract the match ID from the URL parameter
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Get current user ID from JWT token
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	// Check if this user is part of the match
	var match Match
	result := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Verify that the user is part of this match
	if match.Offered != userID && match.Accepted != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this match"})
		return
	}

	// Get chat history from MongoDB
	messages, err := loadChatHistoryFromMongo(uint(matchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load chat history: %v", err)})
		return
	}

	// Sort messages by time to ensure chronological order
	sortMessagesByTime(messages)

	// Generate the date suggestion message
	dateMessage, err := GenerateDateSuggestion(uint(matchID), messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate date suggestion: %v", err)})
		return
	}

	// Add the message to the existing conversation in MongoDB
	err = addMessageToMongo(uint(matchID), *dateMessage)
	if err != nil {
		fmt.Printf("Warning: Failed to save date suggestion message to MongoDB: %v\n", err)
		// Continue anyway as we can still return the message
	}

	// Broadcast the new message to all connected clients in this chat
	broadcastMessageToMatch(uint(matchID), dateMessage)
	
	// Return the date suggestion message
	c.IndentedJSON(http.StatusOK, dateMessage)
}

// GenerateVibeChatForMatch creates and sends an AI-generated conversation starter
func GenerateVibeChatForMatch(c *gin.Context) {
	// Extract the match ID from the URL parameter
	matchIDStr := c.Param("id")
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Get current user ID from JWT token (optional validation)
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDValue.(uint)

	// Check if this user is part of the match
	var match Match
	result := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Verify that the user is part of this match
	if match.Offered != userID && match.Accepted != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not part of this match"})
		return
	}

	// Get chat history from MongoDB
	messages, err := loadChatHistoryFromMongo(uint(matchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load chat history: %v", err)})
		return
	}

	// Sort messages by time to ensure chronological order
	sortMessagesByTime(messages)

	// Generate the vibe chat message
	vibeMessage, err := GenerateVibeChat(uint(matchID), messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate vibe chat: %v", err)})
		return
	}

	// Add the message to the existing conversation in MongoDB
	err = addMessageToMongo(uint(matchID), *vibeMessage)
	if err != nil {
		fmt.Printf("Warning: Failed to save vibe chat message to MongoDB: %v\n", err)
		// Continue anyway as we can still return the message
	}

	// Broadcast the new message to all connected clients in this chat
	broadcastMessageToMatch(uint(matchID), vibeMessage)
	
	// Return the vibe chat message
	c.IndentedJSON(http.StatusOK, vibeMessage)
}

func ChatMessagesFromMongo(c *gin.Context) {
	// Extract the match ID from the URL parameter
	str := c.Param("id")
	matchID, err := strconv.Atoi(str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// Optional limit parameter for pagination
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50 // Default to 50 messages
	}

	// Check if MongoDB client is initialized
	if mongoClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB connection not available"})
		return
	}

	// Use a consistent database name across the application
	collection := mongoClient.Database("urmid").Collection("chathistory")
	
	// Create a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the conversation document for this match ID
	var conversation Conversation
	err = collection.FindOne(ctx, bson.M{"match_id": uint(matchID)}).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No history found is a normal situation for new conversations
			c.IndentedJSON(http.StatusOK, []ChatMessage{})
			return
		}
		
		// Other database errors should be reported
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error loading chat history: %v", err)})
		return
	}

	// Sort messages by time to ensure they appear in chronological order
	if len(conversation.Messages) > 0 {
		// Simple in-memory sort by time
		sortMessagesByTime(conversation.Messages)
		fmt.Printf("Loaded %d messages for match %d from MongoDB\n", len(conversation.Messages), matchID)
	}

	// If we have more messages than the limit, return only the most recent ones
	messages := conversation.Messages
	if len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}

	c.IndentedJSON(http.StatusOK, messages)
}