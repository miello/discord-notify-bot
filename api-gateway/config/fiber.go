package config

import (
	"api-gateway/handler"
	"fmt"

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

	app.Get("/api/test/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	app.Get("/api/courses", courseHandler.GetAllCourses)
	app.Get("/api/:id/assignments", assignmentHandler.GetAssignments)
	app.Get("/api/:id/materials", materialHandler.GetMaterials)

	return app, nil
}
