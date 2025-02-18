package models

import (
	"gorm.io/gorm"
	"time"
)

type IssueRegistry struct {
	gorm.Model
	ISBN               string    `gorm:"not null" json:"isbn"`
	ReaderID           uint      `gorm:"not null" json:"reader_id"`
	IssueApproverID    uint      `gorm:"not null" json:"issue_approver_id"`
	IssueStatus        string    `gorm:"not null" json:"issue_status"`
	IssueDate          time.Time `gorm:"autoCreateTime" json:"issue_date"`
	ExpectedReturnDate time.Time `gorm:"not null" json:"expected_return_date"`
	ReturnDate         *time.Time `json:"return_date"`
	ReturnApproverID   *uint     `json:"return_approver_id"`
	LibraryID          uint      `gorm:"not null" json:"library_id"`
}
