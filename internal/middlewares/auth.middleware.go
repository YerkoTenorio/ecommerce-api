package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

//JWTProtected es un middleware que protege a las rutas que requieren autenticación

func JWTProtected(c *fiber.Ctx) error {
	// extraer el token de la cabecera de autorización
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se proporcionó un token de autorización",
		})
	}

	// comprobar si el encabezado de autorización tiene el formato correcto

	tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se proporcionó un token de autorización",
		})
	}

	// verificar y parsear el token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Método de firma no válido")
		}
		// retornar la clave secreta

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token invalido o expirado",
		})

	}
	// continuar con la solicitud si el token es válido

	return c.Next()
}
