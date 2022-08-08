package models

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string
	Href     string
	CourseID string `gorm:"index"`
	Course   Course
	DueDate  time.Time `gorm:"not null"`
}
