package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	Password      string `gorm:"not null"` // stored as bcrypt hash
	ContactNumber string
	Role          string `gorm:"not null"` // "Owner", "LibraryAdmin", or "Reader"
	LibraryID     uint   // foreign key to Library
}
