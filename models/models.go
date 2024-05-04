package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name" gorm:"text; not null; default: null"`
	Description string `json:"description" gorm:"text; not null; default: null"`
	Image       string `json:"image" gorm:"text; not null; default: null"`
	Users       int    `json:"users" gorm:"text; not null; default: null"`
}