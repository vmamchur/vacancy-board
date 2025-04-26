package repository

import (
	"context"

	"github.com/vmamchur/vacancy-board/db/generated"
	"github.com/vmamchur/vacancy-board/internal/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, dto model.CreateRefreshTokenDTO) (*model.RefreshToken, error)
	GetUser(ctx context.Context, refreshToken string) (*model.User, error)
	Revoke(ctx context.Context, refreshToken string) error
}

type refreshTokenRepository struct {
	q *generated.Queries
}

func NewRefreshTokenRepository(q *generated.Queries) RefreshTokenRepository {
	return &refreshTokenRepository{q}
}

func (r *refreshTokenRepository) Create(ctx context.Context, dto model.CreateRefreshTokenDTO) (*model.RefreshToken, error) {
	dbRefreshToken, err := r.q.CreateRefreshToken(ctx, generated.CreateRefreshTokenParams{
		Token:     dto.Token,
		UserID:    dto.UserID,
		ExpiresAt: dto.ExpiresAt,
		RevokedAt: dto.RevokedAt,
	})
	if err != nil {
		return nil, err
	}

	return toModelRefreshToken(dbRefreshToken), nil
}

func (r *refreshTokenRepository) GetUser(ctx context.Context, refreshToken string) (*model.User, error) {
	dbUser, err := r.q.GetUserFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return toModelUser(dbUser), nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, refreshToken string) error {
	err := r.q.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func toModelRefreshToken(rt generated.RefreshToken) *model.RefreshToken {
	return &model.RefreshToken{
		Token:     rt.Token,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
		UserID:    rt.UserID,
		ExpiresAt: rt.ExpiresAt,
		RevokedAt: rt.RevokedAt,
	}
}
