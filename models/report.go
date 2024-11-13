package models

import (
	"time"

	"github.com/google/uuid"
)

type ReportStatus string

const (
	Pending   ReportStatus = "In Progress"
	Completed ReportStatus = "Completed"
)

type Report struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Location    string       `json:"location"`
	HotelCount  int          `json:"hotel_count"`
	PhoneCount  int          `json:"phone_count"`
	RequestedAt time.Time    `json:"requested_at"`
	Status      ReportStatus `json:"status"`
}

// NewReport creates a new Report instance
func NewReport(location string, hotelCount, phoneCount int) *Report {
	return &Report{
		ID:          uuid.New(),
		Location:    location,
		HotelCount:  hotelCount,
		PhoneCount:  phoneCount,
		RequestedAt: time.Now(),
		Status:      Pending,
	}
}
