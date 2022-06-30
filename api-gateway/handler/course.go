package handler

import (
	"api-gateway/services/scraper"
	"api-gateway/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CourseHandler struct {
	courseService scraper.CourseService
}

func NewCourseHandler(DB *gorm.DB) *CourseHandler {
	return &CourseHandler{
		courseService: *scraper.NewCourseService(DB),
	}
}

func (s *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	year := c.Query("year")
	semester := c.Query("semester")
	name := c.Query("name")

	courses, err := s.courseService.GetAvailableCourses(year, semester, name)

	if err != nil {
		status_code, msg := utils.ExtractError(err)

		c.JSON(msg)
		return c.SendStatus(status_code)
	}

	c.JSON(courses)
	return c.SendStatus(200)
}
