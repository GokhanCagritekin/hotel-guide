// internal/hotels/hotel.go

package hotels

import (
	"time"

	"github.com/google/uuid"
)

// Hotel yapısı
type Hotel struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ContactInfo []Contact `json:"contact_info"`
	CreatedAt   time.Time `json:"created_at"`
}

// Contact yapısı - iletişim bilgilerini temsil eder
type Contact struct {
	Type  string `json:"type"`  // Örnek: "Telefon", "Email", "Konum"
	Value string `json:"value"` // Örnek: Telefon numarası veya email adresi
}

// Yeni bir otel oluşturmak için fonksiyon
func NewHotel(name string, contactInfo []Contact) *Hotel {
	return &Hotel{
		ID:          uuid.New(),
		Name:        name,
		ContactInfo: contactInfo,
		CreatedAt:   time.Now(),
	}
}
