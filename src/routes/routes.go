package routes

import (
	"github.com/gofiber/fiber/v2"
)

var globalRouter = fiber.New()

func sampleRoute(router *fiber.App) {
	router.Get("/", func(context *fiber.Ctx) error {
		return context.SendString("All fine")
	})
}

func SetupRoutes() *fiber.App {
	// place all routers from routes

	sampleRoute(globalRouter)

	globalRouter.Mount("/", GetUserRoutes())
	globalRouter.Mount("/", GetAuthRoutes())

	return globalRouter
}
