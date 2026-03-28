package service

import (
	"context"

	"github.com/edwardmccormick/gowebserver/internal/store"
)

type AuthService struct {
	users store.UserStore
}

func NewAuthService(users store.UserStore) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}
