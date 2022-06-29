package scraper

import (
	"api-gateway/services/scraper"
	"strconv"
	"strings"

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

	courses, err := s.courseService.GetAvailableCourses(year, semester)

	if err != nil {
		arr := strings.Split(err.Error(), ": ")
		status_code, _ := strconv.Atoi(arr[0])

		msg := strings.Join(arr[1:], " ")

		return c.Status(status_code).SendString(msg)
	}

	c.JSON(courses)
	return c.SendStatus(200)
}
