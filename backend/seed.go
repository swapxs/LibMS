// /backend/seed.go
package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Adjust the import path to match your module name and location of your models.
	"github.com/swapxs/LibMS/backend/models"
)

func main() {
	// Read the DSN from an environment variable.
	dsn := "host=localhost user=postgres password=postgres dbname=lms port=5432 sslmode=disable TimeZone=UTC"

	// Connect to the database.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the models.
	err = db.AutoMigrate(
		&models.Library{},
		&models.User{},
		&models.BookInventory{},
		&models.RequestEvent{},
		&models.IssueRegistry{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	// Seed Libraries: 3 libraries.
	libraryNames := []string{"Library A", "Library B", "Library C"}
	var libraries []models.Library
	for _, name := range libraryNames {
		lib := models.Library{
			Name: name,
		}
		if err := db.Create(&lib).Error; err != nil {
			log.Fatalf("Error creating library %s: %v", name, err)
		}
		fmt.Printf("Created library: %s with ID %d\n", name, lib.ID)
		libraries = append(libraries, lib)
	}

	// Hash common password "123123123".
	password := "123123123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}
	hashedStr := string(hashedPassword)

	// Seed Users:
	// For each library: 1 Owner, 2 Admins, 4 Readers.
	for _, lib := range libraries {
		// Create Owner.
		ownerEmail := fmt.Sprintf("owner%d@example.com", lib.ID)
		owner := models.User{
			Name:          fmt.Sprintf("Owner %d", lib.ID),
			Email:         ownerEmail,
			Password:      hashedStr,
			ContactNumber: "123123123",
			Role:          "Owner",
			LibraryID:     lib.ID,
		}
		if err := db.Create(&owner).Error; err != nil {
			log.Fatalf("Error creating owner for library %d: %v", lib.ID, err)
		}
		fmt.Printf("Created Owner: %s\n", ownerEmail)

		// Create 2 Admins.
		for i := 1; i <= 2; i++ {
			adminEmail := fmt.Sprintf("admin%d_%d@example.com", lib.ID, i)
			admin := models.User{
				Name:          fmt.Sprintf("Admin %d-%d", lib.ID, i),
				Email:         adminEmail,
				Password:      hashedStr,
				ContactNumber: "123123123",
				Role:          "LibraryAdmin",
				LibraryID:     lib.ID,
			}
			if err := db.Create(&admin).Error; err != nil {
				log.Fatalf("Error creating admin for library %d: %v", lib.ID, err)
			}
			fmt.Printf("Created Admin: %s\n", adminEmail)
		}

		// Create 4 Readers.
		for i := 1; i <= 4; i++ {
			readerEmail := fmt.Sprintf("reader%d_%d@example.com", lib.ID, i)
			reader := models.User{
				Name:          fmt.Sprintf("Reader %d-%d", lib.ID, i),
				Email:         readerEmail,
				Password:      hashedStr,
				ContactNumber: "123123123",
				Role:          "Reader",
				LibraryID:     lib.ID,
			}
			if err := db.Create(&reader).Error; err != nil {
				log.Fatalf("Error creating reader for library %d: %v", lib.ID, err)
			}
			fmt.Printf("Created Reader: %s\n", readerEmail)
		}
	}

	// Seed Books:
	// For each library, create 12 books (total of 36 books across 3 libraries).
	for _, lib := range libraries {
		for i := 1; i <= 12; i++ {
			isbn := fmt.Sprintf("ISBN-%d-%d", lib.ID, i)
			book := models.BookInventory{
				ISBN:            isbn,
				LibraryID:       lib.ID,
				Title:           fmt.Sprintf("Book Title %d-%d", lib.ID, i),
				Author:          fmt.Sprintf("Author %d-%d", lib.ID, i),
				Publisher:       fmt.Sprintf("Publisher %d", lib.ID),
				Language:        "English",
				Version:         "1st",
				TotalCopies:     10,
				AvailableCopies: 10,
			}
			if err := db.Create(&book).Error; err != nil {
				log.Fatalf("Error creating book for library %d: %v", lib.ID, err)
			}
			fmt.Printf("Created Book: %s in Library %d\n", isbn, lib.ID)
		}
	}

	fmt.Println("Seeding completed successfully.")
}
