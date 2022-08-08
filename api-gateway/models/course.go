package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	Key          string `gorm:"index"`
	Title        string `gorm:"index"`
	Href         string
	Semester     int
	Year         int
	Announcement []Announcement
	Material     []Material
	Assignment   []Assignment
}
