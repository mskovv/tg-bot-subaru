package models

import (
	"fmt"
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

func (a Appointment) String() string {
	return fmt.Sprintf(
		"Дата: %s\nВремя: %s\nМарка: %s\nМодель: %s\nЗадача: %s",
		a.Date.Format("2006-01-02"), // Formatting the date as YYYY-MM-DD
		a.Time.Format("15:04"),      // Formatting the time as HH:mm
		a.CarMark,
		a.CarModel,
		a.Description,
	)
}
