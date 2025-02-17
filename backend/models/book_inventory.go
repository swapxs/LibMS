package models

import "gorm.io/gorm"

type BookInventory struct {
	gorm.Model
	ISBN            string `gorm:"not null;uniqueIndex:idx_book_lib"`
	LibraryID       uint   `gorm:"not null;uniqueIndex:idx_book_lib"`
	Title           string `gorm:"not null"`
	Author         string `gorm:"not null"`
    Publisher       string `grom:"not null"`
    Language        string `grom:"not null"`
    Version         string  `grom:not null`
	TotalCopies     int `gorm:"not null"`
	AvailableCopies int `gorm:"not null"`
}
