package config

import (
	"api-gateway/services/scraper"
	"time"

	"github.com/go-co-op/gocron"
)

func StartUpdateJob() (*gocron.Scheduler, error) {
	announcement := scraper.NewAnnouncementCron(DB)
	material := scraper.NewMaterialCron(DB)
	course := scraper.NewCourseCron(DB)
	assignment := scraper.NewAssignmentCron(DB)

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("12:00").Do(func() {
		course.UpdateCourses()
		announcement.UpdateAnnouncements()
		material.UpdateMaterial()
		assignment.UpdateAssignment()
	})

	s.StartAsync()

	return s, err
}
