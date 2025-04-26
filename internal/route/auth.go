package route

import (
	"net/http"

	"github.com/vmamchur/vacancy-board/internal/handler"
)

func RegisterAuthRoutes(mux *http.ServeMux, authHandler *handler.AuthHandler) {
	mux.HandleFunc("/auth/register", authHandler.Register)
}
