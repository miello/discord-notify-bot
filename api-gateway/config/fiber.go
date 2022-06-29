package config

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetupFiber() (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	return app, nil
}
