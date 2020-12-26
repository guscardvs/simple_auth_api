package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gustcorrea/simple_auth_api/database"
	"github.com/gustcorrea/simple_auth_api/routes"
)

func main() {
	app := fiber.New()

	database.InitDatabase()

	app.Mount("/", routes.SetupRoutes())

	app.Listen(":3000")
}
