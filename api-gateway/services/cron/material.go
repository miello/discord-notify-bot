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

type materialData struct {
	title string
	href  string
}

type MaterialCron struct {
	DB     *gorm.DB
	course *CourseCron
}

func NewMaterialCron(db *gorm.DB) *MaterialCron {
	return &MaterialCron{
		DB:     db,
		course: NewCourseCron(db),
	}
}

func extractMaterialDetail(s *goquery.Selection) materialData {
	BASE_URL := os.Getenv("BASE_URL")
	title_el := s.Find("td[data-col='title'] > a")

	title := title_el.Text()
	href, _ := title_el.Attr("href")

	action_el := s.Find("td[data-col='action'] > a")
	if len(action_el.Nodes) != 0 {
		href, _ = action_el.Attr("href")
		if !strings.HasPrefix(href, "https://") {
			href = fmt.Sprintf("%v%v", BASE_URL, href)
		}
	} else {
		href = fmt.Sprintf("%v%v", BASE_URL, href)
	}

	return materialData{
		title,
		href,
	}
}

// This supposes to be used only in internal cron job, it must not leak to handler
func (c *MaterialCron) UpdateMaterial() error {
	fmt.Println("Start update all materials")

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

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
			return err
		}

		folder := doc.Find("section[aria-label='Course Materials'] *[data-folder]")
		others := doc.Find("section[aria-label='Course Materials'] > * > table tbody tr")

		folder.Each(func(i int, s *goquery.Selection) {
			folder_title := s.Find("button div[data-part='title']").Text()
			s.Find("table > tbody > tr").Each(func(j int, mat *goquery.Selection) {
				material := extractMaterialDetail(mat)
				material_struct := models.Material{
					Title:      material.title,
					Href:       material.href,
					FolderName: folder_title,
					CourseID:   row.ID,
				}

				c.DB.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "href"}},
					DoUpdates: clause.AssignmentColumns([]string{"title", "href", "folder_name", "course_id"}),
				}).Create(&material_struct)
			})
		})

		others.Each(func(i int, s *goquery.Selection) {
			material := extractMaterialDetail(s)
			material_struct := models.Material{
				Title:      material.title,
				Href:       material.href,
				FolderName: "Others",
				CourseID:   row.ID,
			}

			c.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "href"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "href", "folder_name", "course_id"}),
			}).Create(&material_struct)
		})

		res.Body.Close()
		log.Printf("Update materials %v successfully\n", row.Title)
		time.Sleep(5 * time.Second)
	}

	return nil
}
