// /backend/src/models/book_inventory.go
package models

import "gorm.io/gorm"

type BookInventory struct {
	gorm.Model
	ISBN            string `gorm:"not null;uniqueIndex:idx_book_lib"`
	LibraryID       uint   `gorm:"not null;uniqueIndex:idx_book_lib"`
	Title           string `gorm:"not null"`
	Author          string `gorm:"not null"`
	Publisher       string `gorm:"not null"`
	Language        string `gorm:"not null"`
	Version         string `gorm:"not null"`
	TotalCopies     int    `gorm:"not null"`
	AvailableCopies int    `gorm:"not null"`
}
