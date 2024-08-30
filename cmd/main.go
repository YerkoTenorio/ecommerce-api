package main

import (
	"ecommerce-api/config"
	"ecommerce-api/internal/routes"
	"ecommerce-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadEnv()
	config.GetJWTSecret()
	utils.ConnectDB()
	utils.AutoMigrateDB()

	app := fiber.New()
	routes.RegisterRoutes(app)
	app.Listen(":3000")

}
