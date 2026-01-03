package middleware

import (
	"context"
	"net/http"
	"strconv"

	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWT(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "Unauthorized: No token found", http.StatusUnauthorized)
				return
			}

			claims, err := domain.ValidateToken(cookie.Value, cfg.JWTSecret)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims["sub"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int64, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, false
		}
		return id, true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}
