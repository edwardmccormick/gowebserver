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

	// Return the created Person record
	c.IndentedJSON(http.StatusCreated, newPerson)
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

func GetPeople(c *gin.Context) {
	var people []Person
	if result := db.Preload("Photos").Preload("Profile").Find(&people); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, people)
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
	results := db.Where("offered = ?", uint(id)).Or("accepted = ?", uint(id)).Preload("OfferedProfile").Preload("AcceptedProfile").Find(&matches)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, matches)
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
	results := db.Where("offered = ?", uint(id)).Or("accepted = ?", uint(id)).Preload("OfferedProfile").Preload("AcceptedProfile").Find(&matches)
	if results.Error != nil {
		fmt.Println(results.Error)
	}
	if results.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "matches for person not found"})
		return
	}

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

	// fmt.Println(people) // Print the people slice to the console for debugging pu
	c.IndentedJSON(http.StatusOK, messages)
}

func ChatMessagesFromMongo(c *gin.Context) {
	// // Load the configuration
	// var config *Config
	// var err error

	// if isRunningInDockerContainer() {
	// 	config, err = LoadConfig("./config.json") // Adjust the path as needed
	// 	if err != nil {
	// 		fmt.Println("Error loading config:", err)
	// 		return
	// 	}
	// } else {
	// 	config, err = LoadConfig("./configlocal.json") // Adjust the path as needed
	// 	if err != nil {
	// 		fmt.Println("Error loading config:", err)
	// 		return
	// 	}
	// }

	// mongoClient, err := ConnectToMongoDBWithConfig(config)
	// if err != nil {
	// 	fmt.Println("Error connecting to MongoDB:", err)
	// 	return
	// }
	// fmt.Println("Connected to MongoDB.")

	// // Extract the ID from the URL parameter
	str := c.Param("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	collection := mongoClient.Database("gowebserver").Collection("chathistory")

	var result struct {
		History []ChatMessage `bson:"history"`
	}
	err = collection.FindOne(context.TODO(), bson.M{"ID": id}).Decode(&result)
	if err != nil {
		fmt.Printf("Error loading chat history from MongoDB: %v\n", err)
		return
	}

	// Send the last 30 messages to the user

	c.IndentedJSON(http.StatusOK, collection)
}

