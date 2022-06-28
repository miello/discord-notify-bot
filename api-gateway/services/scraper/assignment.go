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

func UpdateAssignment() error {
	fmt.Println("Start update all assignment")
	BASE_URL := os.Getenv("BASE_URL")

	var all_course []models.Course
	tx := config.DB.Find(&all_course)

	if tx.Error != nil {
		log.Fatalf("Error occured. Error is: %s", tx.Error)
		return tx.Error
	}

	for _, row := range all_course {
		path := fmt.Sprintf("/?q=courseville/course/%v/assignment", row.ID)
		res, err := utils.GetHTML(path)

		if err != nil {
			log.Fatal(err)
			return err
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			return err
		}

		doc.Find("table[title='Assignment list'] > tbody > tr").Each(func(i int, s *goquery.Selection) {
			td_el := s.Find("td")
			title_col := td_el.Find("a")

			title := strings.TrimSpace(title_col.Text())

			href, _ := title_col.Attr("href")

			href_split := strings.Split(href, "/")
			assignment_id := href_split[len(href_split)-1]

			href = fmt.Sprintf("%v%v", BASE_URL, href)

			due_date := strings.Split(s.Find(".cv-due-col").Find(".sr-only").Text(), " ")
			due_date_text := strings.Join(due_date[2:], " ")

			config.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "href", "date", "course_id"}),
			}).Create(&models.Assignment{
				ID:       assignment_id,
				Title:    title,
				Href:     href,
				Date:     due_date_text,
				CourseID: row.ID,
			})
		})

		res.Body.Close()

		log.Printf("Update assignment %v successfully\n", row.Title)
		time.Sleep(10 * time.Second)
	}

	return nil
}
