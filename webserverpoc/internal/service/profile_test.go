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

func TestProfileServiceSearchPeopleAppliesDistanceAndCriteria(t *testing.T) {
	db := newTestDB(t)
	service := NewProfileService(db, nil, "")

	people := []domain.Person{
		{
			Name:         "Current",
			Age:          30,
			LatLocation:  29.4241,
			LongLocation: -98.4936,
			Details: domain.Details{
				Dogs:        5,
				EnergyLevel: 5,
			},
		},
		{
			Name:         "Strong Match",
			Age:          31,
			LatLocation:  29.4600,
			LongLocation: -98.5000,
			Details: domain.Details{
				Dogs:        8,
				EnergyLevel: 7,
			},
		},
		{
			Name:         "Wrong Trait",
			Age:          32,
			LatLocation:  29.4500,
			LongLocation: -98.4800,
			Details: domain.Details{
				Dogs:        2,
				EnergyLevel: 7,
			},
		},
		{
			Name:         "Too Far",
			Age:          33,
			LatLocation:  30.1000,
			LongLocation: -98.9000,
			Details: domain.Details{
				Dogs:        8,
				EnergyLevel: 7,
			},
		},
	}
	for _, person := range people {
		if err := db.Create(&person).Error; err != nil {
			t.Fatalf("seed person: %v", err)
		}
	}

	minDogs := 7
	minEnergy := 6
	results, err := service.SearchPeople(1, SearchOptions{
		Distance: 25,
		Criteria: map[string]SearchCriterion{
			"dogs": {
				Category: "dogs",
				Type:     "min",
				Value:    &minDogs,
			},
			"energy": {
				Category: "energy_levels",
				Type:     "min",
				Value:    &minEnergy,
			},
		},
	})
	if err != nil {
		t.Fatalf("SearchPeople returned error: %v", err)
	}

	if len(results) != 1 || results[0].Name != "Strong Match" {
		t.Fatalf("expected only Strong Match, got %+v", results)
	}
}

func TestProfileServiceSearchPeopleSupportsRangeAliases(t *testing.T) {
	db := newTestDB(t)
	service := NewProfileService(db, nil, "")

	people := []domain.Person{
		{
			Name:         "Current",
			Age:          30,
			LatLocation:  41.0,
			LongLocation: -87.0,
		},
		{
			Name:         "Outdoorsy",
			Age:          31,
			LatLocation:  41.1,
			LongLocation: -87.1,
			Details: domain.Details{
				Outdoorsyness: 7,
				Bouginess:     4,
			},
		},
		{
			Name:         "Not Outdoorsy Enough",
			Age:          31,
			LatLocation:  41.1,
			LongLocation: -87.1,
			Details: domain.Details{
				Outdoorsyness: 2,
				Bouginess:     4,
			},
		},
	}
	for _, person := range people {
		if err := db.Create(&person).Error; err != nil {
			t.Fatalf("seed person: %v", err)
		}
	}

	minOutdoor := 5
	maxOutdoor := 8
	exactBougie := 4
	results, err := service.SearchPeople(1, SearchOptions{
		Distance: 50,
		Criteria: map[string]SearchCriterion{
			"outdoor": {
				Category: "outdoorsy_ness",
				Type:     "range",
				Min:      &minOutdoor,
				Max:      &maxOutdoor,
			},
			"bougie": {
				Category: "bougieness",
				Type:     "exact",
				Value:    &exactBougie,
			},
		},
	})
	if err != nil {
		t.Fatalf("SearchPeople returned error: %v", err)
	}

	if len(results) != 1 || results[0].Name != "Outdoorsy" {
		t.Fatalf("expected only Outdoorsy, got %+v", results)
	}
}

func TestProfileServiceSearchPeopleAppliesIdentityPreferenceAndRelationship(t *testing.T) {
	db := newTestDB(t)
	service := NewProfileService(db, nil, "")

	people := []domain.Person{
		{
			Name:             "Current",
			Age:              30,
			GenderIdentity:   "man",
			InterestedIn:     "women",
			RelationshipGoal: "dating",
			LatLocation:      29.4241,
			LongLocation:     -98.4936,
		},
		{
			Name:             "Compatible Woman",
			Age:              29,
			GenderIdentity:   "woman",
			InterestedIn:     "men",
			RelationshipGoal: "dating",
			LatLocation:      29.4300,
			LongLocation:     -98.4900,
		},
		{
			Name:             "Wrong Audience",
			Age:              29,
			GenderIdentity:   "woman",
			InterestedIn:     "women",
			RelationshipGoal: "dating",
			LatLocation:      29.4300,
			LongLocation:     -98.4900,
		},
		{
			Name:             "Wrong Goal",
			Age:              29,
			GenderIdentity:   "woman",
			InterestedIn:     "men",
			RelationshipGoal: "friendship",
			LatLocation:      29.4300,
			LongLocation:     -98.4900,
		},
	}
	for _, person := range people {
		if err := db.Create(&person).Error; err != nil {
			t.Fatalf("seed person: %v", err)
		}
	}

	results, err := service.SearchPeople(1, SearchOptions{
		Distance:     25,
		Gender:       "a man",
		Preference:   "women",
		Relationship: "dating",
	})
	if err != nil {
		t.Fatalf("SearchPeople returned error: %v", err)
	}

	if len(results) != 1 || results[0].Name != "Compatible Woman" {
		t.Fatalf("expected only Compatible Woman, got %+v", results)
	}
}
