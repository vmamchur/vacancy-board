package route

import (
	"net/http"

	"github.com/vmamchur/vacancy-board/internal/handler"
	"github.com/vmamchur/vacancy-board/internal/middleware"
)

func RegisterAuthRoutes(mux *http.ServeMux, authHandler *handler.AuthHandler, appSecret string) {
	jwtMiddleware := middleware.NewJWTMiddleware(appSecret)

	mux.Handle("POST /auth/register", http.HandlerFunc(authHandler.Register))
	mux.Handle("POST /auth/login", http.HandlerFunc(authHandler.Login))
	mux.Handle("POST /auth/refresh", http.HandlerFunc(authHandler.RefreshTokens))

	mux.Handle("POST /auth/logout", jwtMiddleware(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /auth/me", jwtMiddleware(http.HandlerFunc(authHandler.GetMe)))
}
