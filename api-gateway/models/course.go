package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Key      string `gorm:"index" json:"courseId"`
	Title    string `gorm:"index" json:"courseTitle"`
	Href     string `json:"courseHref"`
	Semester int    `json:"semester"`
	Year     int    `json:"year"`
}

type CourseView struct {
	Key      string `json:"courseId"`
	Title    string `json:"courseTitle"`
	Href     string `json:"courseHref"`
	Semester int    `json:"semester"`
	Year     int    `json:"year"`
}
