package models

import "gorm.io/gorm"

type Material struct {
	gorm.Model
	Title    string
	Href     string
	CourseID string
	Course   Course
}
