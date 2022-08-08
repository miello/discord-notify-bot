package handler

import (
	"api-gateway/services/scraper"
	"api-gateway/types"
	"api-gateway/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MaterialHandler struct {
	service *scraper.MaterialService
}

func NewMaterialHandler(db *gorm.DB) *MaterialHandler {
	return &MaterialHandler{
		service: scraper.NewMaterialService(db),
	}
}

func (h *MaterialHandler) GetMaterials(c *fiber.Ctx) error {
	id := c.Params("id")
	folderName := c.Query("folder")

	if id == "" {
		c.JSON(&types.ResponseError{
			Msg: "Required course id",
		})
		return c.SendStatus(400)
	}

	res, err := h.service.GetMaterials(id, folderName)
	if err != nil {
		status_code, body := utils.ExtractError(err)

		c.JSON(body)
		return c.SendStatus(status_code)
	}

	c.JSON(res)

	return c.SendStatus(200)
}
