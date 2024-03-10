package models

import "gorm.io/gorm"
import "time"

type Word struct {
	gorm.Model
	Text         string
	Translation1 string
	Translation2 string
	ImagePath    string
	UserID       uint
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}