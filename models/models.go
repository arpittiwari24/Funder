package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name" gorm:"text; not null; default: null"`
	Description string `json:"description" gorm:"text; not null; default: null"`
	Image       string `json:"image" gorm:"text; not null; default: null"`
	Users       int    `json:"users" gorm:"int; not null; default: null"`
	Email       string `json:"email" gorm:"text;not null; default: null"`
	Url         string `json:"url" gorm:"text; not null; default: null"`
}
