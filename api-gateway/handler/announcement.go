package handler

import (
	"api-gateway/services/scraper"
	"api-gateway/types"
	"api-gateway/utils"
	"strconv"

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
	_id := c.Params("id", "")
	_page := c.Query("page", "1")
	_limit := c.Query("limit", "10")

	id := _id
	page, _ := strconv.Atoi(_page)
	limit, _ := strconv.Atoi(_limit)

	if id == "" {
		c.JSON(&types.ResponseError{
			Msg: "Required course id",
		})
		return c.SendStatus(400)
	}

	res, err := h.service.GetAnnouncements(id, page, limit)
	if err != nil {
		status_code, body := utils.ExtractError(err)

		c.JSON(body)
		return c.SendStatus(status_code)
	}

	c.JSON(res)

	return nil
}

func (h *AnnouncementHandler) GetOverviewAnnouncement(c *fiber.Ctx) error {
	var query types.IGetOverviewQuery
	c.QueryParser(&query)

	if query.Page == 0 {
		query.Page = 1
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	res, err := h.service.GetOverviewAnnouncements(query.Id, query.Page, query.Limit)
	if err != nil {
		status_code, body := utils.ExtractError(err)

		c.JSON(body)
		return c.SendStatus(status_code)
	}

	c.JSON(res)

	return nil
}
