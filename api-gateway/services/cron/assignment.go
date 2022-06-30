package cron

import (
	"api-gateway/models"
	"api-gateway/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssignmentCron struct {
	db *gorm.DB
}

type loadAssignmentBody struct {
	Data struct {
		Html string `json:"html"`
	} `json:"data"`
	Next   int    `json:"next"`
	All    bool   `json:"all"`
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func NewAssignmentCron(db *gorm.DB) *AssignmentCron {
	return &AssignmentCron{
		db,
	}
}

func extractAssignment(s *goquery.Selection, courseId string) models.Assignment {
	loc, _ := time.LoadLocation("Asia/Bangkok")

	BASE_URL := os.Getenv("BASE_URL")
	td_el := s.Find("td")
	title_col := td_el.Find("a")

	title := strings.TrimSpace(title_col.Text())

	href, _ := title_col.Attr("href")

	href_split := strings.Split(href, "/")
	assignment_id := href_split[len(href_split)-1]

	href = fmt.Sprintf("%v%v", BASE_URL, href)

	due_date := strings.Split(s.Find(".cv-due-col").Find(".sr-only").Text(), " ")
	due_date_text := strings.Join(due_date[2:], " ")

	split_date := strings.Split(due_date_text, " ")

	raw_date := fmt.Sprintf("%v %v, %v %v:00", split_date[1], split_date[0], split_date[2], split_date[4])

	date_time, err := time.ParseInLocation("January 2, 2006 15:04:05", raw_date, loc)
	if err != nil {
		fmt.Println(err.Error())
	}

	return models.Assignment{
		ID:       assignment_id,
		Title:    title,
		Href:     href,
		DueDate:  date_time,
		CourseID: courseId,
	}
}

func loadMoreAssignment(db *gorm.DB, courseId string) error {

	body := &loadAssignmentBody{
		Status: 1,
		Next:   5,
	}

	for body.Status != 0 && !body.All {
		println(body.Status, body.All, body.Next)

		form := map[string]string{
			"cv_cid": courseId,
			"next":   strconv.Itoa(body.Next),
		}

		err := utils.GetJSONByFormDataReq("POST", "?q=courseville/ajax/loadmoreassignmentrows", &form, body)

		if err != nil {
			return fmt.Errorf("%v", err.Error())
		}

		html := fmt.Sprintf("<html><body><table>%v</table></body></html>", body.Data.Html)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

		if err != nil {
			return fmt.Errorf("%v", err.Error())
		}

		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			assignment := extractAssignment(s, courseId)
			db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "href", "date", "course_id"}),
			}).Create(&assignment)
		})

		time.Sleep(5 * time.Second)
	}

	return nil
}

// This supposes to be used only in internal cron job, it must not leak to handler
func (c *AssignmentCron) UpdateAssignment() error {
	fmt.Println("Start update all assignment")

	var all_course []models.Course
	tx := c.db.Find(&all_course)

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

		ch := make(chan bool, 1)

		go func() {
			time.Sleep(5 * time.Second)
			loadMoreAssignment(c.db, row.ID)
			ch <- true
		}()

		doc.Find("table[title='Assignment list'] > tbody > tr").Each(func(i int, s *goquery.Selection) {
			assignment := extractAssignment(s, row.ID)

			c.db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "href", "due_date", "course_id"}),
			}).Create(&assignment)
		})

		<-ch

		res.Body.Close()

		log.Printf("Update assignment %v successfully\n", row.Title)
		time.Sleep(5 * time.Second)
	}

	return nil
}
