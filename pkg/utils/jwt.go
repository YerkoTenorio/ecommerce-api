package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(UserID uint, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": UserID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
