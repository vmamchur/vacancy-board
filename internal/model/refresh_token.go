package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

type CreateRefreshTokenDTO struct {
	Token     string       `json:"token"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

type RefreshTokensDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type RevokeRefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}
