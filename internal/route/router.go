package route

import (
	"net/http"

	"github.com/vmamchur/vacancy-board/internal/handler"
)

func NewRouter(authHandler *handler.AuthHandler, appSecret string) http.Handler {
	mux := http.NewServeMux()

	RegisterAuthRoutes(mux, authHandler, appSecret)

	return mux
}
