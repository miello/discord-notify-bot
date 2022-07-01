package scraper

import (
	"api-gateway/models"
	"api-gateway/utils"
	"time"

	"gorm.io/gorm"
)

type AnnouncementService struct {
	db            *gorm.DB
	courseService *CourseService
}

func NewAnnouncementService(db *gorm.DB) *AnnouncementService {
	return &AnnouncementService{
		db:            db,
		courseService: NewCourseService(db),
	}
}

func convertToAnnouncementView(announcement *models.Announcement) models.AnnouncementView {
	return models.AnnouncementView{
		Title: announcement.Title,
		Href:  announcement.Href,
		Date:  announcement.PublishDate.Format(time.RFC3339),
	}
}

func (c *AnnouncementService) GetAnnouncements(id string) ([]models.AnnouncementView, error) {
	found, err := c.courseService.IsCourseIdExists(id)

	if err != nil {
		return nil, utils.CreateError(500, err.Error())
	}

	if !found {
		return nil, utils.CreateError(404, "Not found, maybe api owner does not attend this course")
	}

	query := models.Announcement{
		ID: id,
	}

	var raw_announcement []models.Announcement

	tx := c.db.Where(&query).Find(&raw_announcement)

	if tx.Error != nil {
		return nil, utils.CreateError(500, tx.Error.Error())
	}

	var announcement_view []models.AnnouncementView

	for _, announcement := range raw_announcement {
		announcement_view = append(announcement_view, convertToAnnouncementView(&announcement))
	}

	return announcement_view, nil
}
