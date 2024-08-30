package routes

import (
	"ecommerce-api/internal/controllers"
	"ecommerce-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) {

	//grupo de rutas para productos protegidas por jwt

	productGroup := app.Group("/api/v1/products", middlewares.JWTProtected)

	productGroup.Post("/", controllers.CreateProduct)      //crear un producto
	productGroup.Get("/", controllers.GetProducts)         //obtener todos los productos
	productGroup.Get("/:id", controllers.GetProduct)       //obtener un producto por id
	productGroup.Put("/:id", controllers.UpdateProduct)    //actualizar un producto por id
	productGroup.Delete("/:id", controllers.DeleteProduct) //eliminar un producto por id

}
