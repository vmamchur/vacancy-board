package scraper

import (
	"context"
	"log"
	"time"

	"github.com/vmamchur/vacancy-board/internal/model"
	"github.com/vmamchur/vacancy-board/internal/repository"
)

type Scraper interface {
	Scrape() ([]model.Vacancy, error)
}

type ScraperService struct {
	vacancyRepository repository.VacancyRepository
	scrapers          []Scraper
}

func NewScraper(vacancyRepository repository.VacancyRepository) *ScraperService {
	return &ScraperService{
		vacancyRepository: vacancyRepository,
		scrapers: []Scraper{
			DjinniScraper{},
		},
	}
}

func (s *ScraperService) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	s.scrapeAll(ctx)
}

func (s *ScraperService) scrapeAll(ctx context.Context) {
	for _, scr := range s.scrapers {
		vacancies, err := scr.Scrape()
		if err != nil {
			log.Printf("Error scraping: %v\n", err)
			continue
		}

		for _, v := range vacancies {
			dbVacancy, err := s.vacancyRepository.Create(ctx, model.CreateVacancyDTO{
				Title:       v.Title,
				CompanyName: v.CompanyName,
				Url:         v.Url,
			})
			if err != nil {
				log.Printf("Skipping duplicate or bad vacancy: %s\n", v.Url)
			}
			if dbVacancy != nil {
				log.Printf("Create vacancy: %s\n", v.Url)
			}
		}
	}
}
