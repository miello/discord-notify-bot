package config

import (
	"api-gateway/services/cron"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func StartUpdateJob(db *gorm.DB) (*gocron.Scheduler, error) {
	announcement := cron.NewAnnouncementCron(db)
	material := cron.NewMaterialCron(db)
	course := cron.NewCourseCron(db)
	assignment := cron.NewAssignmentCron(db)

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("12:00").Do(func() {
		log.Println("Start updating job")
		return

		course.UpdateCourses()
		announcement.UpdateAnnouncements()
		material.UpdateMaterial()
	})
	assignment.UpdateAssignment()

	s.StartAsync()

	return s, err
}
