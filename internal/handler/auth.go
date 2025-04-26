package handler

import (
	"net/http"

	"github.com/vmamchur/vacancy-board/internal/middleware"
	"github.com/vmamchur/vacancy-board/internal/model"
	"github.com/vmamchur/vacancy-board/internal/service"
	"github.com/vmamchur/vacancy-board/pkg/httputil"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto model.CreateUserDTO

	if !httputil.DecodeJSONBody(w, r, &dto) {
		return
	}

	tokens, err := h.authService.Register(r.Context(), dto)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto model.LoginDTO

	if !httputil.DecodeJSONBody(w, r, &dto) {
		return
	}

	tokens, err := h.authService.Login(r.Context(), dto)
	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var dto model.RefreshTokensDTO

	if !httputil.DecodeJSONBody(w, r, &dto) {
		return
	}

	tokens, err := h.authService.RefreshTokens(r.Context(), dto.RefreshToken)
	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	var dto model.RevokeRefreshTokenDTO

	if !httputil.DecodeJSONBody(w, r, &dto) {
		return
	}

	err := h.authService.RevokeRefreshToken(r.Context(), dto.RefreshToken)
	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
	}

	httputil.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	user, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
	}

	httputil.RespondWithJSON(w, http.StatusOK, user)
}
