package service

import (
	"time"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"gorm.io/gorm"
)

type MatchService struct {
	db *gorm.DB
}

func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{db: db}
}

func (s *MatchService) GetMatchDetails(matchID uint) (domain.Match, error) {
	var match domain.Match
	if err := s.db.Preload("OfferedProfile").Preload("AcceptedProfile").
		First(&match, matchID).Error; err != nil {
		return domain.Match{}, err
	}
	return match, nil
}

func (s *MatchService) ListMatches() ([]domain.Match, error) {
	var matches []domain.Match
	if err := s.db.Preload("OfferedProfile").Preload("AcceptedProfile").Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) ListMatchesForPerson(personID uint) ([]domain.Match, error) {
	var matches []domain.Match
	if err := s.db.Where("offered = ?", personID).Or("accepted = ?", personID).
		Preload("OfferedProfile.Photos").Preload("OfferedProfile.Profile").
		Preload("AcceptedProfile.Photos").Preload("AcceptedProfile.Profile").
		Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchService) UpsertMatch(input domain.Match) (domain.Match, error) {
	match := input

	if match.ID != 0 {
		existing := domain.Match{}
		result := s.db.Where("offered = ?", match.Offered).
			Preload("OfferedProfile").
			Preload("AcceptedProfile").
			First(&existing)
		if result.Error == nil && result.RowsAffected != 0 {
			match = existing
			match.AcceptedTime = time.Now()
			if err := s.db.Model(&match).Where("id = ?", match.ID).Save(&match).Error; err != nil {
				return domain.Match{}, err
			}
			return match, nil
		}

		existing = domain.Match{}
		result = s.db.Where("offered = ? AND accepted = ?", input.Offered, input.Accepted).
			Or("offered = ? AND accepted = ?", input.Accepted, input.Offered).
			Preload("OfferedProfile").
			Preload("AcceptedProfile").
			First(&existing)
		if result.Error == nil && result.RowsAffected != 0 {
			match = existing
			match.AcceptedTime = time.Now()
			if err := s.db.Model(&match).Where("id = ?", match.ID).Save(&match).Error; err != nil {
				return domain.Match{}, err
			}
			return match, nil
		}

		match = input
		match.OfferedTime = time.Now()
		match.AcceptedTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		if err := s.db.Create(&match).Error; err != nil {
			return domain.Match{}, err
		}
		return match, nil
	}

	match.OfferedTime = time.Now()
	match.AcceptedTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	if err := s.db.Create(&match).Error; err != nil {
		return domain.Match{}, err
	}

	return match, nil
}
