package service

import (
	"context"

	"github.com/vmamchur/vacancy-board/internal/model"
	"github.com/vmamchur/vacancy-board/internal/repository"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *AuthService {
	return &AuthService{userRepository}
}

func (s *AuthService) Register(ctx context.Context, dto model.CreateUserDTO) (*model.User, error) {
	user, err := s.userRepository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}
