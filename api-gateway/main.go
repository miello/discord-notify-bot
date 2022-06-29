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

	db := config.DBConnect()
	err = config.DBMigrate(db)

	if err != nil {
		log.Fatal("Failed to migrate schema")
	}

	scheduler, job_err := config.StartUpdateJob(db)

	if job_err != nil {
		log.Fatal("Failed to init cron job")
	}

	app, err := config.SetupFiber(db)

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
