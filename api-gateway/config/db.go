package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"api-gateway/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() error {
	var err error

	DB_HOST := os.Getenv("DATABASE_HOST")
	DB_PORT := os.Getenv("DATABASE_PORT")
	DB_USER := os.Getenv("DATABASE_USER")
	DB_PASS := os.Getenv("DATABASE_PASSWORD")
	DB_NAME := os.Getenv("DATABASE_NAME")

	DB_CONFIG := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Bangkok", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)

	for {
		DB, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  DB_CONFIG,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})

		if err == nil {
			break
		}

		log.Println("Failed to connect to database retry in 5 seconds")
		time.Sleep(time.Second * 5)
	}

	log.Println("Init database successfully")

	return nil
}

func DBMigrate() error {
	return DB.AutoMigrate(&models.Material{}, &models.Course{}, &models.Assignment{}, &models.Announcement{})
}
