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
	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	http.HandleFunc("/auth/register", authHandler.Register)

	log.Printf("Server listening on: %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
