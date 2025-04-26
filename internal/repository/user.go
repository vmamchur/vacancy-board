package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vmamchur/vacancy-board/db/generated"
	"github.com/vmamchur/vacancy-board/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, dto model.CreateUserDTO) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Get(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type userRepository struct {
	q *generated.Queries
}

func NewUserRepository(q *generated.Queries) UserRepository {
	return &userRepository{q}
}

func (r *userRepository) Create(ctx context.Context, dto model.CreateUserDTO) (*model.User, error) {
	dbUser, err := r.q.CreateUser(ctx, generated.CreateUserParams{
		Email:    dto.Email,
		Password: dto.Password,
	})
	if err != nil {
		return nil, err
	}

	return toModelUser(dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toModelUser(dbUser), nil
}

func (r *userRepository) Get(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	dbUser, err := r.q.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toModelUser(dbUser), nil
}

func toModelUser(u generated.User) *model.User {
	return &model.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
		Password:  u.Password,
	}
}
