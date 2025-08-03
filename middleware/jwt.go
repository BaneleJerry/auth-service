package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/BaneleJerry/auth-service/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	ContextKeyUserClaims contextKey = "userClaims"
)

// JWTAuthMiddleware validates JWT and passes claims in context
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, claims, err := utils.ValidateJWT(tokenStr)
		if err != nil || !token.Valid {

			// Check if the error is a validation error and if it is expired
			if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return

			}

			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ContextKeyUserClaims, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
