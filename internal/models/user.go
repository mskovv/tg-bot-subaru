package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"text;not null;default:null`
	LastName string `json:"LastName" gorm:"text;not null;default:null`
	Nickname string `json:"Nickname" gorm:"text;null;default:null`
	Role     string `gorm:"not null"` // "director" или "master"
}
