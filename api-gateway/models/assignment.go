package models

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string `json:"title"`
	Href     string `json:"href"`
	CourseID string `gorm:"index"`
	Course   Course
	DueDate  time.Time `json:"dueDate"`
}

type AssignmentView struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	DueDate string `json:"dueDate"`
}
