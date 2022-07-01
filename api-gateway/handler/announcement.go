package handler

import (
	"api-gateway/services/scraper"

	"gorm.io/gorm"
)

type AnnouncementHandler struct {
	service *scraper.AnnouncementService
}

func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
	return &AnnouncementHandler{
		service: scraper.NewAnnouncementService(db),
	}
}
