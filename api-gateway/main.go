package main

import (
	"fmt"
	"log"
	"os"

	"api-gateway/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.DBConnect()
	config.DBMigrate()
	app, err := config.SetupFiber()

	if err != nil {
		log.Fatal("Failed to init fiber")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3030"
	}

	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	if err := app.Listen(fmt.Sprintf(":%v", PORT)); err != nil {
		fmt.Println(err.Error())
	}
}
