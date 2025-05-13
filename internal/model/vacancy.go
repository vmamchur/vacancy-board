package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Title       string         `json:"title"`
	CompanyName sql.NullString `json:"company_name"`
	Url         string         `json:"url"`
}

type CreateVacancyDTO struct {
	Title       string         `json:"title"`
	CompanyName sql.NullString `json:"company_name"`
	Url         string         `json:"url"`
}
