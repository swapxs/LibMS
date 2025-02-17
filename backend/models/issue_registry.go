package models

import "time"

type IssueRegistry struct {
	IssueID            uint       `gorm:"primaryKey;autoIncrement"`
	ISBN               string     `gorm:"not null"`
	ReaderID           uint       `gorm:"not null"`
	IssueApproverID    uint       `gorm:"not null"`
	IssueStatus        string     `gorm:"not null"` // e.g., "Issued", "Returned"
	IssueDate          time.Time  `gorm:"autoCreateTime"`
	ExpectedReturnDate time.Time  `gorm:"not null"`
	ReturnDate         *time.Time
	ReturnApproverID   *uint
	LibraryID          uint `gorm:"not null"`
}
