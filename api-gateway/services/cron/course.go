package cron

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

type CourseCron struct {
	db *gorm.DB
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

func (c *CourseCron) GetTargetScraperCourse(all_course *[]models.Course) error {
	year := os.Getenv("YEAR")
	semester := os.Getenv("SEMESTER")

	var query models.Course
	if year != "" {
		query.Year, _ = strconv.Atoi(year)
	}

	if semester != "" {
		query.Semester, _ = strconv.Atoi(semester)
	}

	tx := c.db.Where(&query).Find(all_course)
	return tx.Error
}
