package config

import (
	"api-gateway/handler"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"gorm.io/gorm"
)

func generateKey(c *fiber.Ctx) string {
	return c.Request().URI().String()
}

func SetupFiber(db *gorm.DB) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	app.Use(cache.New(cache.Config{
		Expiration:   60 * time.Minute,
		KeyGenerator: generateKey,
		CacheControl: true,
	}))

	courseHandler := handler.NewCourseHandler(db)
	assignmentHandler := handler.NewAssignmentHandler(db)
	materialHandler := handler.NewMaterialHandler(db)
	announcementHandler := handler.NewAnnouncementHandler(db)

	app.Get("/api/courses", courseHandler.GetAllCourses)
	app.Get("/api/:id/assignments", assignmentHandler.GetAssignments)
	app.Get("/api/:id/materials", materialHandler.GetMaterials)
	app.Get("/api/:id/announcements", announcementHandler.GetAnnouncement)
	app.Get("/api/assignments/overview", assignmentHandler.GetOverviewAssignment)
	app.Get("/api/announcements/overview", announcementHandler.GetOverviewAnnouncement)

	return app, nil
}
