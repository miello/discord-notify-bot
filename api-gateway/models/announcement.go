package models

import "gorm.io/gorm"

type Announcement struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string `json:"title" gorm:"index"`
	Href     string `json:"href"`
	CourseID string
	Course   Course
	Date     string `json:"publishDate"`
}

type AnnouncementView struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Date  string `json:"publishDate"`
}
