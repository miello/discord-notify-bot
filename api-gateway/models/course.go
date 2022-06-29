package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Key      string `json:"courseId" gorm:"index"`
	Title    string `gorm:"index"`
	Href     string `json:"courseHref"`
	Semester int    `json:"semester"`
	Year     int    `json:"year"`
}
