package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gustcorrea/simple_auth_api/controllers"
)

var authRouter = fiber.New()

func GetAuthRoutes() *fiber.App {

	authRouter.Post("/token", controllers.Authenticate)

	return authRouter
}
