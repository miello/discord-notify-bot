package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	Href       string `gorm:"primaryKey"`
	Title      string
	FolderName string
	CourseID   string `gorm:"index"`
	Course     Course
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type File struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

type MaterialView struct {
	FolderName string `json:"folderName"`
	File       []File `json:"file"`
}
