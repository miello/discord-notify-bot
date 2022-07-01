package cron

import (
	"api-gateway/models"
	"api-gateway/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AnnouncementCron struct {
	db     *gorm.DB
	course *CourseCron
}

func NewAnnouncementCron(db *gorm.DB) *AnnouncementCron {
	return &AnnouncementCron{
		db:     db,
		course: NewCourseCron(db),
	}
}

// This supposes to be used only in internal cron job, it must not leak to handler
func (c *AnnouncementCron) UpdateAnnouncements() error {
	BASE_URL := os.Getenv("BASE_URL")

	var all_course []models.Course
	err := c.course.GetTargetScraperCourse(&all_course)

	if err != nil {
		return err
	}

	for _, row := range all_course {
		path := fmt.Sprintf("/?q=courseville/course/%v", row.ID)
		res, err := utils.GetHTML(path)

		if err != nil {
			log.Fatalf("Error occured. Error is: %s", err.Error())
			return err
		}

		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			return fmt.Errorf("error with status code: %v", res.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
			return err
		}

		loc, _ := time.LoadLocation("Asia/Bangkok")

		doc.Find("table[title='Course announcements'] > tbody > tr").Each(func(i int, s *goquery.Selection) {
			td_root := s.Find("td")

			date_root := td_root.Children().First()
			desc_root := td_root.Next().Children().First()

			split_date := strings.Split(date_root.Text(), " ")

			date := fmt.Sprintf("20%v-%v-%v", split_date[2], split_date[1], split_date[0])

			time_date, _ := time.ParseInLocation("2006-Jan-02", date, loc)

			title := desc_root.Text()
			href, _ := desc_root.Attr("href")

			href_split := strings.Split(href, "/")
			id := href_split[len(href_split)-1]

			href = fmt.Sprintf("%v%v", BASE_URL, href)

			announcement := models.Announcement{
				ID:          id,
				Title:       title,
				Href:        href,
				PublishDate: time_date,
				CourseID:    row.ID,
			}

			c.db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "href", "publish_date", "course_id"}),
			}).Create(&announcement)
		})

		res.Body.Close()
		log.Printf("Update announcement %v successfully\n", row.Title)
		time.Sleep(5 * time.Second)
	}
	return nil
}
