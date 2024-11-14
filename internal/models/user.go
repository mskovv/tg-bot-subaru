package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string        `gorm:"text;not null;default:null"`
	LastName     string        `gorm:"text;not null;default:null"`
	Nickname     string        `gorm:"text;null;default:null"`
	Role         string        `gorm:"not null"`
	Appointments []Appointment `gorm:"foreignKey:MasterID"`
}
