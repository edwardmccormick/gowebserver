package service

import (
	"fmt"
	"math"
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

type SearchCriterion struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Value    *int   `json:"value,omitempty"`
	Min      *int   `json:"min,omitempty"`
	Max      *int   `json:"max,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

type SearchOptions struct {
	Distance     float64                    `json:"distance"`
	Gender       string                     `json:"gender,omitempty"`
	Preference   string                     `json:"preference,omitempty"`
	Relationship string                     `json:"relationship,omitempty"`
	Criteria     map[string]SearchCriterion `json:"criteria"`
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

func (s *ProfileService) SearchPeople(currentUserID uint, options SearchOptions) ([]domain.Person, error) {
	currentUser, err := s.GetPersonByID(currentUserID)
	if err != nil {
		return nil, err
	}

	query := s.db.Model(&domain.Person{}).
		Where("id <> ?", currentUserID).
		Preload("Photos").
		Preload("Profile")

	if desiredGender := desiredGenderIdentity(options.Preference); desiredGender != "" {
		query = query.Where("gender_identity = ?", desiredGender)
	}

	if seekerGender := normalizeSearchGender(options.Gender); seekerGender != "" {
		query = query.Where("(interested_in = ? OR interested_in = ?)", seekerInterestValue(seekerGender), "anybody")
	}

	if desiredRelationship := normalizeRelationshipGoal(options.Relationship); desiredRelationship != "" {
		query = query.Where("relationship_goal = ?", desiredRelationship)
	}

	if currentUser.LatLocation != 0 && currentUser.LongLocation != 0 && options.Distance > 0 {
		latRange, longRange := searchBounds(currentUser.LatLocation, options.Distance)
		query = query.Where(
			"lat_location BETWEEN ? AND ? AND long_location BETWEEN ? AND ?",
			currentUser.LatLocation-latRange,
			currentUser.LatLocation+latRange,
			currentUser.LongLocation-longRange,
			currentUser.LongLocation+longRange,
		)
	}

	for _, criterion := range options.Criteria {
		if criterion.Enabled || criterion.Category != "" {
			column, ok := searchColumnForCategory(criterion.Category)
			if !ok {
				continue
			}

			switch criterion.Type {
			case "min":
				if criterion.Value != nil {
					query = query.Where(column+" >= ?", *criterion.Value)
				}
			case "max":
				if criterion.Value != nil {
					query = query.Where(column+" <= ?", *criterion.Value)
				}
			case "range":
				if criterion.Min != nil && criterion.Max != nil {
					minValue := *criterion.Min
					maxValue := *criterion.Max
					if minValue > maxValue {
						minValue, maxValue = maxValue, minValue
					}
					query = query.Where(column+" BETWEEN ? AND ?", minValue, maxValue)
				}
			default:
				if criterion.Value != nil {
					query = query.Where(column+" = ?", *criterion.Value)
				}
			}
		}
	}

	var people []domain.Person
	if err := query.Find(&people).Error; err != nil {
		return nil, err
	}

	if currentUser.LatLocation == 0 || currentUser.LongLocation == 0 || options.Distance <= 0 {
		return people, nil
	}

	filtered := make([]domain.Person, 0, len(people))
	for _, person := range people {
		if person.LatLocation == 0 && person.LongLocation == 0 {
			continue
		}
		if haversineMiles(currentUser.LatLocation, currentUser.LongLocation, person.LatLocation, person.LongLocation) <= options.Distance {
			filtered = append(filtered, person)
		}
	}

	return filtered, nil
}

func searchBounds(lat, distanceMiles float64) (latRange float64, longRange float64) {
	const milesPerLatDegree = 69.0
	latRange = distanceMiles / milesPerLatDegree

	cosLat := math.Cos(lat * math.Pi / 180)
	if math.Abs(cosLat) < 0.0001 {
		return latRange, 180
	}

	longRange = distanceMiles / (milesPerLatDegree * cosLat)
	return latRange, longRange
}

func haversineMiles(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMiles = 3958.8
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMiles * c
}

func searchColumnForCategory(category string) (string, bool) {
	switch category {
	case "bouginess", "bougieness":
		return "bouginess", true
	case "cats":
		return "cats", true
	case "dogs":
		return "dogs", true
	case "drinking":
		return "drinking", true
	case "energy_level", "energy_levels":
		return "energy_level", true
	case "food":
		return "food", true
	case "importance_of_politics":
		return "importance_of_politics", true
	case "kids":
		return "kids", true
	case "outdoorsyness", "outdoorsy_ness":
		return "outdoorsyness", true
	case "religion":
		return "religion", true
	case "smoking":
		return "smoking", true
	case "travel":
		return "travel", true
	default:
		return "", false
	}
}

func normalizeSearchGender(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "a woman", "woman", "female":
		return "woman"
	case "a man", "man", "male":
		return "man"
	case "nonbinary", "non-binary", "nonbinary folks":
		return "nonbinary"
	default:
		return ""
	}
}

func desiredGenderIdentity(preference string) string {
	switch strings.TrimSpace(strings.ToLower(preference)) {
	case "women", "woman":
		return "woman"
	case "men", "man":
		return "man"
	case "nonbinary folks", "nonbinary", "non-binary":
		return "nonbinary"
	default:
		return ""
	}
}

func seekerInterestValue(gender string) string {
	switch gender {
	case "woman":
		return "women"
	case "man":
		return "men"
	case "nonbinary":
		return "nonbinary"
	default:
		return ""
	}
}

func normalizeRelationshipGoal(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", "?", "who cares":
		return ""
	default:
		return strings.TrimSpace(strings.ToLower(value))
	}
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
