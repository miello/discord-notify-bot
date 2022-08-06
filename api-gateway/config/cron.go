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

	location, err := time.LoadLocation("Asia/Bangkok")

	if err != nil {
		log.Println("Unfortunately can't load a location")
		log.Println(err)
	} else {
		s.ChangeLocation(location)
	}

	_, err = s.Every(1).Day().At("0:00").Do(func() {
		log.Println("Start updating job at", time.Now())

		course.UpdateCourses()
		assignment.UpdateAssignment()
		announcement.UpdateAnnouncements()
		material.UpdateMaterial()
	})

	s.StartAsync()

	return s, err
}
