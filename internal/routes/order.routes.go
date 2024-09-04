package routes

import (
	"ecommerce-api/internal/controllers"
	"ecommerce-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App) {

	//grupo de rutas para ordenes protegidas por jwt

	orderGroup := app.Group("/api/v1/orders", middlewares.AuthMiddleware)

	orderGroup.Post("/", controllers.CreateOrder)      //crear una orden
	orderGroup.Get("/", controllers.GetOrders)         //obtener todas las ordenes
	orderGroup.Get("/:id", controllers.GetOrder)       //obtener una orden por id
	orderGroup.Put("/:id", controllers.UpdateOrder)    //actualizar una orden por id
	orderGroup.Delete("/:id", controllers.DeleteOrder) //eliminar una orden por id

}
