package models

import (
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	ID          uint      `gorm:"primaryKey"`
	Date        time.Time `gorm:"not null"`
	CarModel    string    `gorm:"not null"`
	CarMark     string    `gorm:"not null"`
	Description string
	gorm.Model
}
