package gowebserver

import (
	"fmt"

	"gorm.io/gorm"
)

func MigrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&UserSession{},
		&Person{},
		&Details{},
		&ProfilePhoto{},
		&Match{},
		&ChatMessage{},
		&MatchReadState{},
	)
}

func SeedDemoData(db *gorm.DB) error {
	if err := seedUsers(db); err != nil {
		return err
	}
	if err := seedPeople(db); err != nil {
		return err
	}
	if err := seedMatches(db); err != nil {
		return err
	}
	return nil
}

func seedPeople(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Person{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count people: %w", err)
	}
	if count > 0 {
		return backfillSeedPeople(db)
	}
	seedPeople := seededPeople()
	if err := db.Create(&seedPeople).Error; err != nil {
		return fmt.Errorf("seed people: %w", err)
	}
	return nil
}

func seedUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	seedUsers := make([]User, len(users))
	for i, user := range users {
		seedUsers[i] = user
		seedUsers[i].ID = uint(i + 1)
	}

	if err := db.Create(&seedUsers).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	return nil
}

func seedMatches(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Match{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count matches: %w", err)
	}
	if count > 0 {
		return nil
	}
	if err := db.Create(&Matches).Error; err != nil {
		return fmt.Errorf("seed matches: %w", err)
	}
	return nil
}

func seededPeople() []Person {
	seedPeople := make([]Person, len(people))
	copy(seedPeople, people)

	for i := range seedPeople {
		if identity, ok := seedIdentities[seedPeople[i].ID]; ok {
			seedPeople[i].GenderIdentity = identity.GenderIdentity
			seedPeople[i].InterestedIn = identity.InterestedIn
			seedPeople[i].RelationshipGoal = identity.RelationshipGoal
		}
	}

	return seedPeople
}

func backfillSeedPeople(db *gorm.DB) error {
	for personID, identity := range seedIdentities {
		if err := db.Model(&Person{}).
			Where("id = ? AND (gender_identity = '' OR interested_in = '' OR relationship_goal = '')", personID).
			Updates(map[string]interface{}{
				"gender_identity":   identity.GenderIdentity,
				"interested_in":     identity.InterestedIn,
				"relationship_goal": identity.RelationshipGoal,
			}).Error; err != nil {
			return fmt.Errorf("backfill seeded person %d: %w", personID, err)
		}
	}

	return nil
}
