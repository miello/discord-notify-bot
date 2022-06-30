package handler

import (
	"api-gateway/models"
	"api-gateway/services/scraper"
	"api-gateway/utils"

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
	id := c.Params("id")

	if id == "" {
		c.JSON(&models.ResponseError{
			Msg: "Required course id",
		})
		return c.SendStatus(400)
	}

	res, err := h.assignmentService.GetAssignments(id)
	if err != nil {
		status_code, body := utils.ExtractError(err)

		c.JSON(body)
		return c.SendStatus(status_code)
	}

	c.JSON(res)

	return c.SendStatus(200)
}
