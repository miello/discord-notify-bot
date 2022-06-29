package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"api-gateway/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.DBConnect()
	err = config.DBMigrate()

	if err != nil {
		log.Fatal("Failed to migrate schema")
	}

	scheduler, job_err := config.StartUpdateJob()

	if job_err != nil {
		log.Fatal("Failed to init cron job")
	}

	app, err := config.SetupFiber()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
		scheduler.Clear()
	}()

	if err != nil {
		log.Fatal("Failed to init fiber")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3030"
	}

	if err := app.Listen(fmt.Sprintf(":%v", PORT)); err != nil {
		fmt.Println(err.Error())
	}
}
