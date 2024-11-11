// internal/reports/report.go

package reports

import (
	"time"

	"github.com/google/uuid"
)

// Report yapısı - bir raporu temsil eder
type Report struct {
	ID          uuid.UUID `json:"id"`
	Location    string    `json:"location"`
	HotelCount  int       `json:"hotel_count"`
	PhoneCount  int       `json:"phone_count"`
	RequestedAt time.Time `json:"requested_at"`
	Status      string    `json:"status"` // "Hazırlanıyor", "Tamamlandı"
}

// Yeni bir rapor talebi oluşturmak için fonksiyon
func NewReport(location string) *Report {
	return &Report{
		ID:          uuid.New(),
		Location:    location,
		RequestedAt: time.Now(),
		Status:      "Hazırlanıyor",
	}
}

// Rapor durumunu güncellemek için fonksiyon
func (r *Report) CompleteReport(hotelCount, phoneCount int) {
	r.HotelCount = hotelCount
	r.PhoneCount = phoneCount
	r.Status = "Tamamlandı"
}
