package models

import (
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	gorm.Model
	Date        time.Time `gorm:"not null"`
	Time        time.Time `gorm:"not null"`
	CarMark     string    `gorm:"not null"`
	CarModel    string    `gorm:"not null"`
	Description string
}
