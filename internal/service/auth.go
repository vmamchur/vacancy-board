package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vmamchur/vacancy-board/config"
	"github.com/vmamchur/vacancy-board/internal/model"
	"github.com/vmamchur/vacancy-board/internal/repository"
	"github.com/vmamchur/vacancy-board/pkg/auth"
)

type AuthService struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	appSecret              string
}

func NewAuthService(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	appSecret string,
) *AuthService {
	return &AuthService{userRepository, refreshTokenRepository, appSecret}
}

func (s *AuthService) Register(ctx context.Context, dto model.CreateUserDTO) (*model.AuthResponseDTO, error) {
	hashedPassword, err := auth.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	dto.Password = hashedPassword

	_, err = s.userRepository.GetByEmail(ctx, dto.Email)
	if err == nil {
		return nil, errors.New("User with this email already exists")
	}

	user, err := s.userRepository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	tokens, err := s.CreateTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) Login(ctx context.Context, dto model.LoginDTO) (*model.AuthResponseDTO, error) {
	user, err := s.userRepository.GetByEmail(ctx, dto.Email)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	err = auth.CheckPasswordHash(dto.Password, user.Password)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	tokens, err := s.CreateTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*model.AuthResponseDTO, error) {
	user, err := s.refreshTokenRepository.GetUser(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("Token is invalid")
	}

	tokens, err := s.CreateTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	err := s.refreshTokenRepository.Revoke(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) CreateTokens(ctx context.Context, userID uuid.UUID) (*model.AuthResponseDTO, error) {
	rawRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.refreshTokenRepository.Create(ctx, model.CreateRefreshTokenDTO{
		Token:     rawRefreshToken,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(config.RefreshTokenTTL),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := auth.MakeJWT(userID, s.appSecret, config.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponseDTO{
		RefreshToken: refreshToken.Token,
		AccessToken:  accessToken,
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		return nil, errors.New("Couldn't get user")
	}

	return user, nil
}
