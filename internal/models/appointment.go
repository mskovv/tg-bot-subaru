package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	gorm.Model
	Date        time.Time `gorm:"not null;index:idx_appointment_date"`
	Time        time.Time `gorm:"type:time;not null"`
	CarMark     string    `gorm:"not null"`
	CarModel    string    `gorm:"not null"`
	Description string
}

func (a Appointment) String() string {
	return fmt.Sprintf(
		"Дата: %s\nВремя: %s\nМарка: %s\nМодель: %s\nЗадача: %s",
		a.Date.Format("2006-01-02"),
		a.Time.Format("15:04"),
		a.CarMark,
		a.CarModel,
		a.Description,
	)
}
