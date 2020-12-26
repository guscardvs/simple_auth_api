package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gustcorrea/simple_auth_api/controllers"
)

var userRouter = fiber.New()

// GetUserRoutes provides userRouter with all routes created
func GetUserRoutes() *fiber.App {

	userRouter.Post("/create", controllers.CreateUser)

	privateRoute := userRouter.Group("/user")
	privateRoute.Use(controllers.SecureAuth)
	privateRoute.Get("/", controllers.GetUser)
	privateRoute.Put("/", controllers.EditUser)
	privateRoute.Put("/password", controllers.ChangePassword)

	return userRouter
}
