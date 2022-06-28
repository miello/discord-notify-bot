package models

import (
	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string `gorm:"index"`
	Href     string
	CourseID string
	Course   Course
	Date     string
}
