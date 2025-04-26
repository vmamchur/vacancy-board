package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/vmamchur/vacancy-board/config"
	"github.com/vmamchur/vacancy-board/db/generated"
	"github.com/vmamchur/vacancy-board/internal/handler"
	"github.com/vmamchur/vacancy-board/internal/repository"
	"github.com/vmamchur/vacancy-board/internal/route"
	"github.com/vmamchur/vacancy-board/internal/service"
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

	userRepository := repository.NewUserRepository(q)
	refreshTokenRepository := repository.NewRefreshTokenRepository(q)
	authService := service.NewAuthService(userRepository, refreshTokenRepository, cfg.AppSecret)
	authHandler := handler.NewAuthHandler(authService)

	router := route.NewRouter(authHandler, cfg.AppSecret)

	log.Printf("Server listening on: %s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, router))
}
