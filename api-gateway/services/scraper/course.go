package scraper

import (
	"api-gateway/models"
	"api-gateway/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourseService struct {
	db *gorm.DB
}

type CourseCron struct {
	db *gorm.DB
}

func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{
		db,
	}
}

func NewCourseCron(db *gorm.DB) *CourseCron {
	return &CourseCron{
		db,
	}
}

// This supposes to be used only in internal cron job, it must not leak to handler
func (c *CourseCron) UpdateCourses() error {
	BASE_URL := os.Getenv("BASE_URL")
	res, err := utils.GetHTML("/?q=courseville")

	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	doc.Find("*[course_no]").Each(func(i int, s *goquery.Selection) {
		course_id, _ := s.Attr("cv_cid")
		course_no, _ := s.Attr("course_no")
		title, _ := s.Attr("title")
		href, _ := s.Attr("href")

		semester, _ := s.Attr("semester")
		semester_int, _ := strconv.Atoi(semester)

		year, _ := s.Attr("year")
		year_int, _ := strconv.Atoi(year)

		course := models.Course{
			Title:    title,
			Key:      course_no,
			Href:     fmt.Sprintf("%v%v", BASE_URL, href),
			ID:       course_id,
			Semester: semester_int,
			Year:     year_int,
		}

		c.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "key", "href", "semester", "year"}),
		}).Create(&course)
	})

	time.Sleep(5 * time.Second)
	log.Println("Update available courses successfully")

	return nil
}

func (c *CourseService) GetAvailableCourses(year string, semester string) ([]models.Course, error) {
	var all_course []models.Course
	var err error

	query := &models.Course{}

	if semester != "" {
		var semester_num int
		semester_num, err = strconv.Atoi(semester)
		query.Semester = semester_num
	}

	if err != nil {
		return nil, fmt.Errorf("400: Failed to parse semester")
	}

	if year != "" {
		var year_num int
		year_num, err = strconv.Atoi(year)
		query.Year = year_num
	}

	if err != nil {
		return nil, fmt.Errorf("400: Failed to parse year")
	}

	tx := c.db.Where(&query).Find(&all_course)
	if tx.Error != nil {
		return nil, fmt.Errorf("500: %v", tx.Error.Error())
	}

	return all_course, nil
}
