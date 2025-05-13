package repository

import (
	"context"

	"github.com/vmamchur/vacancy-board/db/generated"
	"github.com/vmamchur/vacancy-board/internal/model"
)

type VacancyRepository interface {
	Create(ctx context.Context, dto model.CreateVacancyDTO) (*model.Vacancy, error)
	GetAll(ctx context.Context) ([]model.Vacancy, error)
}

type vacancyRepository struct {
	q *generated.Queries
}

func NewVacancyRepository(q *generated.Queries) VacancyRepository {
	return &vacancyRepository{q}
}

func (r *vacancyRepository) Create(ctx context.Context, dto model.CreateVacancyDTO) (*model.Vacancy, error) {
	dbVacancy, err := r.q.CreateVacancy(ctx, generated.CreateVacancyParams{
		Title:       dto.Title,
		CompanyName: dto.CompanyName,
		Url:         dto.Url,
	})
	if err != nil {
		return nil, err
	}

	return toModelVacancy(dbVacancy), nil
}

func (r *vacancyRepository) GetAll(ctx context.Context) ([]model.Vacancy, error) {
	dbVacancies, err := r.q.GetAllVacancies(ctx)
	if err != nil {
		return nil, err
	}

	return toModelVacancies(dbVacancies), nil
}

func toModelVacancy(v generated.Vacancy) *model.Vacancy {
	return &model.Vacancy{
		ID:          v.ID,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		Title:       v.Title,
		CompanyName: v.CompanyName,
		Url:         v.Url,
	}
}

func toModelVacancies(vs []generated.Vacancy) []model.Vacancy {
	vacancies := make([]model.Vacancy, 0, len(vs))
	for _, v := range vs {
		vacancies = append(vacancies, *toModelVacancy(v))
	}
	return vacancies
}
