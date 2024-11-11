package hotel

import (
	"github.com/google/uuid"
)

type ContactInfo struct {
	InfoType    string `json:"info_type"` // Telefon, E-mail, Konum gibi
	InfoContent string `json:"info_content"`
}

type Hotel struct {
	ID           uuid.UUID     `json:"id"`
	OwnerName    string        `json:"owner_name"`
	OwnerSurname string        `json:"owner_surname"`
	CompanyTitle string        `json:"company_title"`
	Contacts     []ContactInfo `json:"contacts"`
}

// NewHotel creates a new Hotel instance
func NewHotel(ownerName, ownerSurname, companyTitle string, contacts []ContactInfo) *Hotel {
	return &Hotel{
		ID:           uuid.New(),
		OwnerName:    ownerName,
		OwnerSurname: ownerSurname,
		CompanyTitle: companyTitle,
		Contacts:     contacts,
	}
}
