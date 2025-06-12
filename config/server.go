package config

import (
	"kronos-scheduler/routes"

	"github.com/gofiber/fiber/v2"
)

func CreateServer() *fiber.App {
	app := fiber.New()

	routes.SetupRoutes(app)

	return app
}
