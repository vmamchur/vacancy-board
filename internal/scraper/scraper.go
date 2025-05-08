package scraper

import (
	"log"
	"time"

	"github.com/vmamchur/vacancy-board/internal/repository"
)

type Scraper interface {
	Scrape() error
}

type ScraperService struct {
	vacancyRepository repository.VacancyRepository
	scrapers          []Scraper
}

func NewScraper(vacancyRepository repository.VacancyRepository) *ScraperService {
	return &ScraperService{
		vacancyRepository: vacancyRepository,
		scrapers: []Scraper{
			DjinniScraper{vacancyRepository: vacancyRepository},
		},
	}
}

func (s *ScraperService) Run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	s.scrapeAll()
}

func (s *ScraperService) scrapeAll() {
	for _, scr := range s.scrapers {
		err := scr.Scrape()
		if err != nil {
			log.Printf("Error scraping: %v\n", err)
			continue
		}
	}
}
