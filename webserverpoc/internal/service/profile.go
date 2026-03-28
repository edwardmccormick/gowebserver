package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/edwardmccormick/gowebserver/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProfileService struct {
	db         *gorm.DB
	s3Client   *s3.S3
	bucketName string
}

func NewProfileService(db *gorm.DB, s3Client *s3.S3, bucketName string) *ProfileService {
	return &ProfileService{
		db:         db,
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (s *ProfileService) SaveProfile(input domain.Person) (domain.Person, error) {
	person := input

	if err := s.db.Omit("Photos", "Profile").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&person).Error; err != nil {
		return domain.Person{}, err
	}

	var existingPhotos []domain.ProfilePhoto
	if err := s.db.Where("person_id = ?", person.ID).Find(&existingPhotos).Error; err != nil {
		return domain.Person{}, err
	}

	existingS3Keys := make(map[string]uint)
	for _, photo := range existingPhotos {
		existingS3Keys[photo.S3Key] = photo.ID
	}

	requestedPhotoKeys := make(map[string]bool)
	for i := range person.Photos {
		photo := &person.Photos[i]
		photo.PersonID = person.ID
		requestedPhotoKeys[photo.S3Key] = true

		if existingID, exists := existingS3Keys[photo.S3Key]; exists {
			photo.ID = existingID
			if err := s.db.Model(&domain.ProfilePhoto{}).Where("id = ?", existingID).Updates(map[string]interface{}{
				"caption": photo.Caption,
			}).Error; err != nil {
				return domain.Person{}, err
			}
		} else {
			photo.ID = 0
			if err := s.db.Create(photo).Error; err != nil {
				return domain.Person{}, err
			}
		}
	}

	for _, existingPhoto := range existingPhotos {
		if !requestedPhotoKeys[existingPhoto.S3Key] && !strings.HasSuffix(existingPhoto.S3Key, "/profile") {
			if err := s.db.Delete(&domain.ProfilePhoto{}, existingPhoto.ID).Error; err != nil {
				continue
			}
		}
	}

	return person, nil
}

func (s *ProfileService) GenerateUploadTargets(personID uint) ([]domain.ProfilePhoto, []domain.ProfilePhoto, error) {
	if s.s3Client == nil {
		return nil, nil, fmt.Errorf("s3 client is not initialized")
	}

	var uploadURLs []domain.ProfilePhoto
	var profileUploadURLs []domain.ProfilePhoto

	profileKey := fmt.Sprintf("%d/profile", personID)
	profileReq, _ := s.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(profileKey),
	})
	profileURL, err := profileReq.Presign(180 * time.Minute)
	if err == nil {
		profileUploadURLs = append(profileUploadURLs, domain.ProfilePhoto{
			Url:      profileURL,
			S3Key:    profileKey,
			PersonID: personID,
		})
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d/image%d", personID, i+1)
		req, _ := s.s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		})
		url, err := req.Presign(180 * time.Minute)
		if err != nil {
			continue
		}
		uploadURLs = append(uploadURLs, domain.ProfilePhoto{
			Url:      url,
			S3Key:    key,
			PersonID: personID,
		})
	}

	return uploadURLs, profileUploadURLs, nil
}

func (s *ProfileService) ListUsers() ([]domain.User, error) {
	var users []domain.User
	if err := s.db.Preload("Person").Preload("Person.Photos").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *ProfileService) ListPeople() ([]domain.Person, error) {
	var people []domain.Person
	if err := s.db.Preload("Photos").Preload("Profile").Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (s *ProfileService) GetPersonByID(id uint) (domain.Person, error) {
	person := domain.Person{ID: id}
	if err := s.db.Preload("Photos").Preload("Profile").First(&person).Error; err != nil {
		return domain.Person{}, err
	}
	return person, nil
}

func (s *ProfileService) ListPeopleNearLocation(lat, long float64, excludeUserID uint) ([]domain.Person, error) {
	var people []domain.Person
	if err := s.db.Where("lat_location BETWEEN ? AND ? AND long_location BETWEEN ? AND ? AND id != ?",
		lat-1, lat+1, long-1, long+1, excludeUserID).
		Preload("Photos").Preload("Profile").Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (s *ProfileService) ListPeopleWithinBounds(minLat, maxLat, minLong, maxLong float64) ([]domain.Person, error) {
	var people []domain.Person
	if err := s.db.Where("lat_location BETWEEN ? AND ? AND long_location BETWEEN ? AND ?", minLat, maxLat, minLong, maxLong).
		Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (s *ProfileService) ListPhotosByPersonID(personID uint) ([]domain.ProfilePhoto, error) {
	var photos []domain.ProfilePhoto
	if err := s.db.Where("person_id = ?", personID).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (s *ProfileService) AttachPhotoViewURLs(photos []domain.ProfilePhoto, duration time.Duration) {
	if s.s3Client == nil {
		return
	}

	for i := range photos {
		req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(photos[i].S3Key),
		})
		url, err := req.Presign(duration)
		if err == nil {
			photos[i].Url = url
		}
	}
}

func (s *ProfileService) AttachPhotoManagementURLs(photos []domain.ProfilePhoto, duration time.Duration) {
	if s.s3Client == nil {
		return
	}

	for i := range photos {
		getReq, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(photos[i].S3Key),
		})
		if getURL, err := getReq.Presign(duration); err == nil {
			photos[i].Url = getURL
		}

		putReq, _ := s.s3Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(photos[i].S3Key),
		})
		if putURL, err := putReq.Presign(duration); err == nil {
			photos[i].Upload = putURL
		}

		deleteReq, _ := s.s3Client.DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(photos[i].S3Key),
		})
		if deleteURL, err := deleteReq.Presign(duration); err == nil {
			photos[i].Delete = deleteURL
		}
	}
}

func (s *ProfileService) HydrateMatchMedia(matches []domain.Match) {
	if s.s3Client == nil {
		return
	}

	for i := range matches {
		s.attachPersonMedia(&matches[i].OfferedProfile)
		s.attachPersonMedia(&matches[i].AcceptedProfile)
	}
}

func (s *ProfileService) attachPersonMedia(person *domain.Person) {
	if person == nil || person.ID == 0 {
		return
	}

	var profilePhoto *domain.ProfilePhoto
	for i := range person.Photos {
		if person.Photos[i].S3Key == fmt.Sprintf("%d/profile", person.ID) {
			profilePhoto = &person.Photos[i]
			break
		}
	}

	if profilePhoto != nil {
		req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(profilePhoto.S3Key),
		})
		if url, err := req.Presign(180 * time.Minute); err == nil {
			person.Profile = *profilePhoto
			person.Profile.Url = url
		}
	}

	var filteredPhotos []domain.ProfilePhoto
	for _, photo := range person.Photos {
		if !strings.Contains(photo.S3Key, "/profile") {
			filteredPhotos = append(filteredPhotos, photo)
		}
	}
	person.Photos = filteredPhotos
	s.AttachPhotoViewURLs(person.Photos, 60*time.Minute)
}
