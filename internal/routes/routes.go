package routes

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	// LLAMADO A TODAS LAS RUTAS
	AuthRoutes(app)
	ProductRoutes(app)
	OrderRoutes(app)
	UserRoutes(app)

}
