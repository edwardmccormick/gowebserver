package service

import (
	"testing"

	"github.com/edwardmccormick/gowebserver/internal/domain"
)

func TestProfileServiceSaveProfileReconcilesPhotos(t *testing.T) {
	db := newTestDB(t)
	service := NewProfileService(db, nil, "")

	person := domain.Person{
		Name:         "Alice",
		Age:          28,
		LatLocation:  41.5,
		LongLocation: -87.6,
		Photos: []domain.ProfilePhoto{
			{S3Key: "1/profile", Caption: "avatar"},
			{S3Key: "1/image1", Caption: "first"},
			{S3Key: "1/image2", Caption: "second"},
		},
	}
	saved, err := service.SaveProfile(person)
	if err != nil {
		t.Fatalf("SaveProfile returned error: %v", err)
	}

	if saved.ID == 0 {
		t.Fatal("expected saved profile to have an ID")
	}

	updated, err := service.SaveProfile(domain.Person{
		ID:           saved.ID,
		Name:         "Alice Updated",
		Age:          29,
		LatLocation:  41.6,
		LongLocation: -87.7,
		Photos: []domain.ProfilePhoto{
			{S3Key: "1/profile", Caption: "avatar updated"},
			{S3Key: "1/image2", Caption: "second updated"},
			{S3Key: "1/image3", Caption: "third"},
		},
	})
	if err != nil {
		t.Fatalf("SaveProfile update returned error: %v", err)
	}

	personFromDB, err := service.GetPersonByID(updated.ID)
	if err != nil {
		t.Fatalf("GetPersonByID returned error: %v", err)
	}

	if personFromDB.Name != "Alice Updated" || personFromDB.Age != 29 {
		t.Fatalf("expected updated person fields, got %+v", personFromDB)
	}
	if len(personFromDB.Photos) != 3 {
		t.Fatalf("expected profile photo plus 2 gallery photos, got %+v", personFromDB.Photos)
	}

	photosByKey := map[string]domain.ProfilePhoto{}
	for _, photo := range personFromDB.Photos {
		photosByKey[photo.S3Key] = photo
	}

	if _, exists := photosByKey["1/image1"]; exists {
		t.Fatalf("expected removed gallery photo to be deleted, got %+v", personFromDB.Photos)
	}
	if photosByKey["1/profile"].Caption != "avatar updated" {
		t.Fatalf("expected profile photo caption update, got %+v", photosByKey["1/profile"])
	}
	if photosByKey["1/image2"].Caption != "second updated" {
		t.Fatalf("expected existing photo caption update, got %+v", photosByKey["1/image2"])
	}
	if photosByKey["1/image3"].Caption != "third" {
		t.Fatalf("expected new photo to be inserted, got %+v", photosByKey["1/image3"])
	}
}

func TestProfileServiceListPeopleNearLocationExcludesCurrentUser(t *testing.T) {
	db := newTestDB(t)
	service := NewProfileService(db, nil, "")

	people := []domain.Person{
		{Name: "Current", Age: 30, LatLocation: 41.0, LongLocation: -87.0},
		{Name: "Nearby", Age: 31, LatLocation: 41.4, LongLocation: -87.4},
		{Name: "Far Away", Age: 32, LatLocation: 50.0, LongLocation: -100.0},
	}
	for _, person := range people {
		if err := db.Create(&person).Error; err != nil {
			t.Fatalf("seed person: %v", err)
		}
	}

	results, err := service.ListPeopleNearLocation(41.0, -87.0, 1)
	if err != nil {
		t.Fatalf("ListPeopleNearLocation returned error: %v", err)
	}

	if len(results) != 1 || results[0].Name != "Nearby" {
		t.Fatalf("expected only nearby non-current user, got %+v", results)
	}
}

func TestProfileServiceGenerateUploadTargetsRequiresS3Client(t *testing.T) {
	service := NewProfileService(nil, nil, "")

	uploads, profiles, err := service.GenerateUploadTargets(5)
	if err == nil {
		t.Fatal("expected missing s3 client error")
	}
	if uploads != nil || profiles != nil {
		t.Fatalf("expected nil uploads on error, got uploads=%+v profiles=%+v", uploads, profiles)
	}
}
