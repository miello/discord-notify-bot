package handler

import (
	"api-gateway/models"
	"api-gateway/services/scraper"
	"api-gateway/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AnnouncementHandler struct {
	service *scraper.AnnouncementService
}

func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
	return &AnnouncementHandler{
		service: scraper.NewAnnouncementService(db),
	}
}

func (h *AnnouncementHandler) GetAnnouncement(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		c.JSON(&models.ResponseError{
			Msg: "Required course id",
		})
		return c.SendStatus(400)
	}

	res, err := h.service.GetAnnouncements(id)
	if err != nil {
		status_code, body := utils.ExtractError(err)

		c.JSON(body)
		return c.SendStatus(status_code)
	}

	c.JSON(res)

	return nil
}
