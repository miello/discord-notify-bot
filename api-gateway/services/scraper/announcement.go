package scraper

import (
	"gorm.io/gorm"
)

type AnnouncementService struct {
	db *gorm.DB
}

func NewAnnouncementService(db *gorm.DB) *AnnouncementService {
	return &AnnouncementService{
		db,
	}
}
