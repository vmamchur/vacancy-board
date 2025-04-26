package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/vmamchur/vacancy-board/pkg/auth"
	"github.com/vmamchur/vacancy-board/pkg/httputil"
)

const userIDContextKey = "userID"

func NewJWTMiddleware(appSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetBearerToken(r.Header)
			if err != nil {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
				return
			}

			userID, err := auth.ValidateJWT(token, appSecret)
			if err != nil {
				httputil.RespondWithError(w, http.StatusUnauthorized, "Invalid token", err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("User ID not found in context")
	}
	return userID, nil
}
