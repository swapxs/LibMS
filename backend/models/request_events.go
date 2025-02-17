package models

import "time"

type RequestEvent struct {
	ReqID        uint       `gorm:"primaryKey;autoIncrement"`
	BookID       string     `gorm:"not null"` // ISBN of the book
	ReaderID     uint       `gorm:"not null"`
	RequestDate  time.Time  `gorm:"autoCreateTime"`
	ApprovalDate *time.Time
	ApproverID   *uint
	RequestType  string `gorm:"not null"` // "Issue", "Approved", "Rejected"
}
