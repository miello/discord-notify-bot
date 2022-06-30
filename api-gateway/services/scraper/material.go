package scraper

import (
	"gorm.io/gorm"
)

type MaterialService struct {
	DB *gorm.DB
}

func NewMaterialService(db *gorm.DB) *MaterialService {
	return &MaterialService{
		DB: db,
	}
}
