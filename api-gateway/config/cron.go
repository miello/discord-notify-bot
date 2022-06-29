package config

import (
	"api-gateway/services/scraper"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func StartUpdateJob(db *gorm.DB) (*gocron.Scheduler, error) {
	announcement := scraper.NewAnnouncementCron(db)
	material := scraper.NewMaterialCron(db)
	course := scraper.NewCourseCron(db)
	assignment := scraper.NewAssignmentCron(db)

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("12:00").Do(func() {
		log.Println("Start updating job")
		return

		course.UpdateCourses()
		announcement.UpdateAnnouncements()
		material.UpdateMaterial()
		assignment.UpdateAssignment()
	})

	s.StartAsync()

	return s, err
}
