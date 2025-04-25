package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vmamchur/vacancy-board/internal/model"
	"github.com/vmamchur/vacancy-board/internal/service"
	"github.com/vmamchur/vacancy-board/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto model.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid input", err)
		return
	}

	user, err := h.authService.Register(r.Context(), dto)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, user)
}
