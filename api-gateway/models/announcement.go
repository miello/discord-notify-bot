package models

import "gorm.io/gorm"

type Announcement struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string `gorm:"index"`
	URL      string
	CourseID string
	Course   Course
	Date     string
}
