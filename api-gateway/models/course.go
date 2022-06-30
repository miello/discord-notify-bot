package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Key      string `gorm:"index"`
	Title    string `gorm:"index"`
	Href     string
	Semester int
	Year     int
}

type CourseView struct {
	ID       string `json:"key"`
	Key      string `json:"courseId"`
	Title    string `json:"courseTitle"`
	Href     string `json:"courseHref"`
	Semester int    `json:"semester"`
	Year     int    `json:"year"`
}
