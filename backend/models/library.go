package models

import "gorm.io/gorm"

type Library struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Users []User
	Books []BookInventory
}
