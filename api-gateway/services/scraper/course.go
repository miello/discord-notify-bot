package scraper

import (
	"api-gateway/models"
	"api-gateway/types"
	"api-gateway/utils"
	"strconv"

	"gorm.io/gorm"
)

type CourseService struct {
	db *gorm.DB
}

func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{
		db,
	}
}

func convertCourseToView(course *models.Course) types.CourseView {
	return types.CourseView{
		ID:       course.ID,
		Key:      course.Key,
		Title:    course.Title,
		Href:     course.Href,
		Semester: course.Semester,
		Year:     course.Year,
	}
}

func (c *CourseService) GetAvailableCourses(year string, semester string, name string) ([]types.CourseView, error) {
	var all_course []models.Course
	var err error

	query := &models.Course{}

	if semester != "" {
		var semester_num int
		semester_num, err = strconv.Atoi(semester)
		query.Semester = semester_num
	}

	if err != nil {
		return nil, utils.CreateError(400, "Failed to parse semester")
	}

	if year != "" {
		var year_num int
		year_num, err = strconv.Atoi(year)
		query.Year = year_num
	}

	if name != "" {
		query.Title = name
	}

	if err != nil {
		return nil, utils.CreateError(400, "Failed to parse year")
	}

	tx := c.db.Where(&query).Find(&all_course)
	if tx.Error != nil {
		return nil, utils.CreateError(500, tx.Error.Error())
	}

	var converted_course []types.CourseView

	for _, c2 := range all_course {
		converted_course = append(converted_course, convertCourseToView(&c2))
	}

	return converted_course, nil
}

func (c *CourseService) GetCourseIdByName(name string) (string, error) {
	var course models.Course

	tx := c.db.First(&models.Course{
		Title: name,
	}).Find(&course)

	if tx.Error != nil {
		return "", utils.CreateError(404, "Not found course")
	}

	return course.ID, nil
}

func (c *CourseService) IsCourseIdExists(id string) (bool, error) {
	var cnt int64

	tx := c.db.Model(&models.Course{}).Where(&models.Course{
		ID: id,
	}).Count(&cnt)

	if tx.Error != nil {
		return false, utils.CreateError(500, tx.Error.Error())
	}

	return cnt != 0, nil
}
