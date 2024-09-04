package routes

import (
	"ecommerce-api/internal/controllers"
	"ecommerce-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	//Rutas protegidas para todos los usuarios autenticados
	UserRoutes := app.Group("/api/v1/user", middlewares.AuthMiddleware)

	UserRoutes.Get("/profile", controllers.GetProfile)
	UserRoutes.Put("/profile", controllers.UpdateProfile)

	//Rutas protegidas para los usuarios con rol de administrador
	AdminRoutes := app.Group("/api/v1/admin/users", middlewares.AuthMiddleware, middlewares.RoleMiddleware("admin"))

	AdminRoutes.Get("/", controllers.GetAllUsers)
	AdminRoutes.Get("/:id", controllers.GetUser)
	AdminRoutes.Put("/:id", controllers.UpdateUser)
	AdminRoutes.Delete("/:id", controllers.DeleteUser)

}
