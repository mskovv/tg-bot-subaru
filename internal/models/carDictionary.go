package models

import "gorm.io/gorm"

type CarDictionary struct {
	gorm.Model
	Mark     string `gorm:"text;not null;default:null"`
	CarModel string `gorm:"text;not null;default:null"`
}
