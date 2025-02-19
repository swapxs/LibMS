// /backend/models/library.go
package models

import "gorm.io/gorm"

type Library struct {
	gorm.Model
	Name  string          `gorm:"unique;not null"`
	Users []User          `gorm:"not null"`
	Books []BookInventory `gorm:"not null"`
}
