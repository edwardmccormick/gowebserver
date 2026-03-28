package service

import (
	"testing"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"gorm.io/gorm"
)

func TestMatchServiceUpsertMatchCreatesNewMatch(t *testing.T) {
	db := newTestDB(t)
	service := NewMatchService(db)

	match, err := service.UpsertMatch(domain.Match{
		Offered:  10,
		Accepted: 20,
	})
	if err != nil {
		t.Fatalf("UpsertMatch returned error: %v", err)
	}

	if match.ID == 0 {
		t.Fatal("expected created match to have an ID")
	}
	if match.OfferedTime.IsZero() {
		t.Fatal("expected offered time to be set")
	}
	if !match.AcceptedTime.Equal(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected sentinel accepted time, got %v", match.AcceptedTime)
	}

	var count int64
	if err := db.Model(&domain.Match{}).Count(&count).Error; err != nil {
		t.Fatalf("count matches: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 persisted match, got %d", count)
	}
}

func TestMatchServiceUpsertMatchAcceptsExistingReciprocalMatch(t *testing.T) {
	db := newTestDB(t)
	existing := domain.Match{
		Offered:      30,
		Accepted:     40,
		OfferedTime:  time.Now().Add(-time.Hour),
		AcceptedTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("seed match: %v", err)
	}

	service := NewMatchService(db)
	match, err := service.UpsertMatch(domain.Match{
		Model:    gorm.Model{ID: 999},
		Offered:  40,
		Accepted: 30,
	})
	if err != nil {
		t.Fatalf("UpsertMatch returned error: %v", err)
	}

	if match.ID != existing.ID {
		t.Fatalf("expected existing match ID %d, got %d", existing.ID, match.ID)
	}
	if !match.AcceptedTime.After(existing.AcceptedTime) {
		t.Fatalf("expected accepted time to be updated, got %v", match.AcceptedTime)
	}

	var count int64
	if err := db.Model(&domain.Match{}).Count(&count).Error; err != nil {
		t.Fatalf("count matches: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected no duplicate match rows, got %d", count)
	}
}

func TestMatchServiceListMatchesForPersonPreloadsProfiles(t *testing.T) {
	db := newTestDB(t)
	offered := domain.Person{Name: "Alice", Age: 29, LatLocation: 41.0, LongLocation: -87.0}
	accepted := domain.Person{Name: "Bob", Age: 31, LatLocation: 42.0, LongLocation: -88.0}
	if err := db.Create(&offered).Error; err != nil {
		t.Fatalf("seed offered person: %v", err)
	}
	if err := db.Create(&accepted).Error; err != nil {
		t.Fatalf("seed accepted person: %v", err)
	}
	if err := db.Create(&domain.ProfilePhoto{PersonID: offered.ID, S3Key: "1/image1", Caption: "offered"}).Error; err != nil {
		t.Fatalf("seed offered photo: %v", err)
	}
	if err := db.Create(&domain.ProfilePhoto{PersonID: accepted.ID, S3Key: "2/image1", Caption: "accepted"}).Error; err != nil {
		t.Fatalf("seed accepted photo: %v", err)
	}
	if err := db.Create(&domain.Match{Offered: offered.ID, Accepted: accepted.ID}).Error; err != nil {
		t.Fatalf("seed match: %v", err)
	}

	service := NewMatchService(db)
	matches, err := service.ListMatchesForPerson(offered.ID)
	if err != nil {
		t.Fatalf("ListMatchesForPerson returned error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].OfferedProfile.Name != "Alice" || matches[0].AcceptedProfile.Name != "Bob" {
		t.Fatalf("expected preloaded profiles, got %+v", matches[0])
	}
	if len(matches[0].OfferedProfile.Photos) != 1 || len(matches[0].AcceptedProfile.Photos) != 1 {
		t.Fatalf("expected preloaded photos, got %+v", matches[0])
	}
}
