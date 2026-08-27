package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"todos/internal/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

type TokenParser interface {
	ParseToken(token string) (int, error)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

func AuthMiddleware(parser TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				utils.WriteJSONResponseError(w, http.StatusUnauthorized, "header authorization is required")
				return
			}

			ok := strings.HasPrefix(authorization, "Bearer ")
			if !ok {
				utils.WriteJSONResponseError(w, http.StatusUnauthorized, "doesn't have bearer prefix")
				return
			}

			tokenString := strings.TrimPrefix(authorization, "Bearer ")

			if len(tokenString) == 0 {
				utils.WriteJSONResponseError(w, http.StatusUnauthorized, "token is required")
				return
			}

			userID, err := parser.ParseToken(tokenString)
			if err != nil {
				slog.Info("invalid jwt token", "error", err)
				utils.WriteJSONResponseError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}

}
