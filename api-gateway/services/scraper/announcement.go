package scraper

import (
	"api-gateway/models"
	"api-gateway/types"
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

func convertToShortAnnouncement(announcement *models.Announcement) types.ShortAnnouncement {
	return types.ShortAnnouncement{
		Title: announcement.Title,
		Href:  announcement.Href,
		Date:  announcement.PublishDate.Format(time.RFC3339),
	}
}

func (c *AnnouncementService) GetAnnouncements(id string, page int, limit int) (types.AnnouncementView, error) {
	found, err := c.courseService.IsCourseIdExists(id)
	var res types.AnnouncementView

	if err != nil {
		return res, utils.CreateError(500, err.Error())
	}

	if !found {
		return res, utils.CreateError(404, "Not found, maybe api owner does not attend this course")
	}

	query := models.Announcement{
		CourseID: id,
	}

	var raw_announcement []models.Announcement
	var number_of_announcement int64

	tx := c.db.Model(&query).Where(&query)

	tx = tx.Count(&number_of_announcement).Offset(utils.GetOffset(page, limit)).Limit(limit).Find(&raw_announcement)

	if tx.Error != nil {
		return res, utils.CreateError(500, tx.Error.Error())
	}

	var short_announcements []types.ShortAnnouncement

	for _, announcement := range raw_announcement {
		short_announcements = append(short_announcements, convertToShortAnnouncement(&announcement))
	}

	res = types.AnnouncementView{
		Announcements: short_announcements,
		Metadata: types.PaginateMetadata{
			CurrentPage: page,
			TotalPages:  utils.GetTotalPages(number_of_announcement, limit),
			TotalItems:  int(number_of_announcement),
		},
	}

	return res, nil
}

func (c *AnnouncementService) GetOverviewAnnouncements(id []string, page int, limit int) (types.OverviewAnnouncementView, error) {
	offset := utils.GetOffset(page, limit)
	date := time.Now().Add(-14 * 24 * time.Hour)
	now_date := time.Now()

	var res types.OverviewAnnouncementView
	var announcement []models.Announcement
	var count int64

	tx := c.db.Model(&models.Announcement{}).Where(gorm.Expr("publish_date > ? AND publish_date < ?", date, now_date)).Where(gorm.Expr("course_id IN (?)", id))
	if tx.Error != nil {
		return res, utils.CreateError(500, tx.Error.Error())
	}

	tx = tx.Count(&count).Order("publish_date desc").Offset(offset).Limit(limit).Preload("Course").Find(&announcement)
	if tx.Error != nil {
		return res, utils.CreateError(500, tx.Error.Error())
	}

	var announcement_view []types.OverviewAnnouncement = make([]types.OverviewAnnouncement, 0)

	for _, a := range announcement {
		announcement_view = append(announcement_view, convertToOverviewAnnouncement(a))
	}

	res = types.OverviewAnnouncementView{
		Announcements: announcement_view,
		Metadata: types.PaginateMetadata{
			CurrentPage: page,
			TotalPages:  utils.GetTotalPages(count, limit),
			TotalItems:  int(count),
		},
	}

	return res, nil
}

func convertToOverviewAnnouncement(a models.Announcement) types.OverviewAnnouncement {
	return types.OverviewAnnouncement{
		CourseTitle: a.Course.Title,
		Title:       a.Title,
		Href:        a.Href,
		Date:        a.PublishDate.Format(time.RFC3339),
	}
}
