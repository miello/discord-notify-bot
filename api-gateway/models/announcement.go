package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Title       string `json:"title"`
	Href        string `json:"href"`
	CourseID    string `gorm:"index"`
	Course      Course
	PublishDate time.Time `json:"publishDate"`
}

type AnnouncementView struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Date  string `json:"publishDate"`
}
