package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID uint, secret string) (string, error) {
	// Crear las claims registradas con el ID del usuario en el campo Subject
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),                          // ID del usuario como Subject
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Expiración en 24 horas
		IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Fecha de emisión
		Issuer:    "ecommerce-api",                                    // Nombre de tu aplicación
		// Puedes agregar más claims estándar si lo necesitas
	}

	// Crear el token con las claims utilizando el método de firma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	return token.SignedString([]byte(secret))
}
