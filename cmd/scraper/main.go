package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/vmamchur/vacancy-board/config"
	"github.com/vmamchur/vacancy-board/db/generated"
	"github.com/vmamchur/vacancy-board/internal/repository"
	"github.com/vmamchur/vacancy-board/internal/scraper"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.Name, cfg.DB.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s\n", err)
	}
	defer db.Close()

	q := generated.New(db)

	vacancyRepository := repository.NewVacancyRepository(q)
	scraper := scraper.NewScraper(vacancyRepository)

	log.Println("Scraper service running...")
	scraper.Run()
}
