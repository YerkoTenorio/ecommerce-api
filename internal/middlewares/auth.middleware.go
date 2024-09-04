package middlewares

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

//AuthMiddleware es un middleware que protege a las rutas que requieren autenticación

func AuthMiddleware(c *fiber.Ctx) error {
	// Extraer el token de la cabecera de autorización
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se proporcionó un token de autorización",
		})
	}

	// Parsear y verificar el token JWT usando RegisteredClaims
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(strings.TrimSpace(strings.Replace(tokenString, "Bearer", "", 1)), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token no válido o expirado",
		})
	}
	fmt.Println(claims)
	// Recuperar el usuario de la base de datos usando el ID del token
	var user models.User
	if err := utils.DB.Where("id = ?", claims.Subject).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Guardar el usuario en el contexto para que esté disponible en otras partes de la aplicación
	c.Locals("user", &user)

	return c.Next()
}

// RoleMiddleware es un middleware que protege a las rutas que requieren un rol específico

func RoleMiddleware(requireRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		if user.Role != requireRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No tienes permiso para acceder a esta ruta",
			})
		}
		return c.Next()
	}
}
