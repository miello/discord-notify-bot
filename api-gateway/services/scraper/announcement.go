package scraper

import (
	"api-gateway/config"
	"api-gateway/models"
	"api-gateway/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm/clause"
)

func UpdateAnnouncements() error {
	BASE_URL := os.Getenv("BASE_URL")

	var all_course []models.Course
	tx := config.DB.Find(&all_course)

	if tx.Error != nil {
		return tx.Error
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

		doc.Find("table[title='Course announcements'] > tbody > tr").Each(func(i int, s *goquery.Selection) {
			td_root := s.Find("td")

			date_root := td_root.Children().First()
			desc_root := td_root.Next().Children().First()

			date := date_root.Text()
			title := desc_root.Text()
			href, _ := desc_root.Attr("href")

			href_split := strings.Split(href, "/")
			id := href_split[len(href_split)-1]

			href = fmt.Sprintf("%v%v", BASE_URL, href)

			announcement := models.Announcement{
				ID:       id,
				Title:    title,
				Href:     href,
				Date:     date,
				CourseID: row.ID,
			}

			config.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "href", "date", "course_id"}),
			}).Create(&announcement)
		})

		time.Sleep(5 * time.Second)
		res.Body.Close()
	}
	return nil
}
