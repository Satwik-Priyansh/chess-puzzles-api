package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKeyType string

const UserIDKey contextKeyType = "userID"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Bearer realm="API Access", error="invalid_token"`)
				http.Error(w, "Unauthorized: Missing authorization header.", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization header format.", http.StatusUnauthorized)
				return
			}
			userID, err := ValidateToken(parts[1], secret)
			if err != nil {
				http.Error(w, "token validation failed.", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}

}
