package config

import (
	"api-gateway/handler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupFiber(db *gorm.DB) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	courseHandler := handler.NewCourseHandler(db)
	assignmentHandler := handler.NewAssignmentHandler(db)
	materialHandler := handler.NewMaterialHandler(db)
	announcementHandler := handler.NewAnnouncementHandler(db)

	app.Get("/api/courses", courseHandler.GetAllCourses)
	app.Get("/api/:id/assignments", assignmentHandler.GetAssignments)
	app.Get("/api/:id/materials", materialHandler.GetMaterials)
	app.Get("/api/:id/announcements", announcementHandler.GetAnnouncement)

	return app, nil
}
