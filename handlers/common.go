package handlers

import (
	"context"

	"github.com/BaneleJerry/auth-service/middleware"
	"github.com/golang-jwt/jwt"
)

func GetUserFromContext(ctx context.Context) (*jwt.StandardClaims, bool) {
	claims, ok := ctx.Value(middleware.ContextKeyUserClaims).(*jwt.StandardClaims)
	return claims, ok
}
