package config

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func SetupFiber() (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	return app, nil
}
