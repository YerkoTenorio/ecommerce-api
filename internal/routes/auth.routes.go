package routes

import (
	"ecommerce-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/api/v1/register", controllers.Register)
	app.Post("/api/v1/login", controllers.Login)
}
