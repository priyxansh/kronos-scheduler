package routes

import (
	"kronos-scheduler/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/healthz", handlers.Healthz)
	app.Post("/prioritize", handlers.Prioritize)
}
