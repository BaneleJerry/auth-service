package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var ErrMissingSecretKey = jwt.NewValidationError("missing SECRET_KEY", jwt.ValidationErrorUnverifiable)

func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "", ErrMissingSecretKey
	}

	expStr := os.Getenv("JWT_EXPIRATION")
	expSeconds, err := strconv.Atoi(expStr)
	if err != nil || expSeconds <= 0 {
		expSeconds = 3600 // default to 1 hour
	}

	claims := jwt.StandardClaims{
		Subject: userID.String(),
		Issuer:   "auth-service",
		ExpiresAt: time.Now().Add(time.Duration(expSeconds) * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, *jwt.StandardClaims, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return nil, nil, ErrMissingSecretKey
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
	}

	return token, claims, nil
}
