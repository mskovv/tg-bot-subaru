package models

import "time"

type Appointment struct {
	ID          uint      `gorm:"primaryKey"`
	Date        time.Time `gorm:"not null"`
	Time        time.Time `gorm:"not null"`
	MasterID    uint      `gorm:"not null"` // Ссылка на мастера
	CarModel    string    `gorm:"not null"`
	Description string
}
