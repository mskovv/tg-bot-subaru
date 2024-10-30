package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"text;not null;default:null`
	LastName string `json:"LastName" gorm:"text;not null;default:null`
	Nickname string `json:"Nickname" gorm:"text;not null;default:null`
}
