package gowebserver

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   newUser.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"email": newUser.Email,
	})
	tokenString, err := token.SignedString(signingKey())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")
	var uploadURLs []ProfilePhoto
	var profileUploadURLs []ProfilePhoto

	profileKey := fmt.Sprintf("%d/profile", newUser.ID)
	profileReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileURL, err := profileReq.Presign(180 * time.Minute)
	if err == nil {
		profileUploadURLs = append(profileUploadURLs, ProfilePhoto{
			Url:      profileURL,
			S3Key:    profileKey,
			PersonID: newUser.ID,
		})
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d/image%d", newUser.ID, i+1)
		req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		url, err := req.Presign(180 * time.Minute)
		if err != nil {
			fmt.Printf("Failed to generate presigned URL for key %s: %v\n", key, err)
			continue
		}
		uploadURLs = append(uploadURLs, ProfilePhoto{
			Url:      url,
			S3Key:    key,
			PersonID: newUser.ID,
		})
	}

	resp := struct {
		Token             string         `json:"token"`
		ID                uint           `json:"id"`
		UploadUrls        []ProfilePhoto `json:"upload_urls,omitempty"`
		ProfileUploadUrls []ProfilePhoto `json:"profile_upload_urls,omitempty"`
	}{
		Token:             tokenString,
		ID:                newUser.ID,
		UploadUrls:        uploadURLs,
		ProfileUploadUrls: profileUploadURLs,
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

	var user *User
	db.Where("email = ?", req.Email).First(&user)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"email": user.Email,
	})
	tokenString, err := token.SignedString(signingKey())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	var person = Person{ID: user.ID}
	db.Preload("Photos").Preload("Profile").First(&person)

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AWS session"})
		return
	}

	s3Client := s3.New(awsSession)
	bucketName := os.Getenv("AWS_S3_BUCKET")
	var uploadURLs []ProfilePhoto
	var profileUploadURLs []ProfilePhoto

	profileKey := fmt.Sprintf("%d/profile", user.ID)
	profileReq, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileURL, err := profileReq.Presign(180 * time.Minute)
	if err == nil {
		profileUploadURLs = append(profileUploadURLs, ProfilePhoto{
			Url:      profileURL,
			S3Key:    profileKey,
			PersonID: user.ID,
		})
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d/image%d", user.ID, i+1)
		req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		url, err := req.Presign(180 * time.Minute)
		if err != nil {
			fmt.Printf("Failed to generate presigned URL for key %s: %v\n", key, err)
			continue
		}
		uploadURLs = append(uploadURLs, ProfilePhoto{
			Url:      url,
			S3Key:    key,
			PersonID: user.ID,
		})
	}

	var wg sync.WaitGroup
	for j := range person.Photos {
		wg.Add(1)
		go func(photo *ProfilePhoto) {
			defer wg.Done()
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(photo.S3Key),
			})
			url, err := req.Presign(180 * time.Minute)
			if err != nil {
				fmt.Printf("Failed to generate presigned URL for S3Key %s: %v\n", photo.S3Key, err)
				return
			}
			photo.Url = url
		}(&person.Photos[j])
	}

	profileGetReq, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(profileKey),
	})
	profileViewURL, err := profileGetReq.Presign(180 * time.Minute)
	if err == nil {
		if person.Profile.ID == 0 {
			person.Profile = ProfilePhoto{
				PersonID: user.ID,
				S3Key:    profileKey,
				Url:      profileViewURL,
			}
		} else {
			person.Profile.Url = profileViewURL
		}
	}

	wg.Wait()

	resp := struct {
		Token             string         `json:"token"`
		Person            Person         `json:"person"`
		UploadUrls        []ProfilePhoto `json:"upload_urls,omitempty"`
		ProfileUploadUrls []ProfilePhoto `json:"profile_upload_urls,omitempty"`
	}{
		Token:             tokenString,
		Person:            person,
		UploadUrls:        uploadURLs,
		ProfileUploadUrls: profileUploadURLs,
	}

	c.JSON(http.StatusOK, resp)
}

func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return signingKey(), nil
	})
	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	revokeToken(token)
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}
