package models

import (
	"gorm.io/gorm"
)

type Assignment struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string `json:"title"`
	Href     string `json:"href"`
	CourseID string
	Course   Course
	Date     string `json:"dueDate"`
}

type AssignmentView struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Date  string `json:"dueDate"`
}
