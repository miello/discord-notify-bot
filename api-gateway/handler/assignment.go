package handler

import (
	"api-gateway/services/scraper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AssignmentHandler struct {
	assignmentService scraper.AssignmentService
}

func NewAssignmentHandler(DB *gorm.DB) *AssignmentHandler {
	return &AssignmentHandler{
		assignmentService: *scraper.NewAssignmentService(DB),
	}
}

func (h *AssignmentHandler) GetAssignments(c *fiber.Ctx) error {

	return nil
}
