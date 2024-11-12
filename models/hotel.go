package models

import (
	"github.com/google/uuid"
)

type ContactInfo struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	HotelID     uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;" json:"hotel_id"`
	InfoType    string    `json:"info_type"` // Telefon, E-mail, Konum gibi
	InfoContent string    `json:"info_content"`
}

type Hotel struct {
	ID           uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OwnerName    string        `json:"owner_name"`
	OwnerSurname string        `json:"owner_surname"`
	CompanyTitle string        `json:"company_title"`
	ContactInfos []ContactInfo `gorm:"foreignKey:HotelID;references:ID;constraint:OnDelete:CASCADE;"`
}

type HotelOfficial struct {
	OwnerName    string `json:"owner_name"`
	OwnerSurname string `json:"owner_surname"`
	CompanyTitle string `json:"company_title"`
}

// NewHotel creates a new Hotel instance
func NewHotel(ownerName, ownerSurname, companyTitle string, contacts []ContactInfo) *Hotel {
	return &Hotel{
		ID:           uuid.New(),
		OwnerName:    ownerName,
		OwnerSurname: ownerSurname,
		CompanyTitle: companyTitle,
		ContactInfos: contacts,
	}
}
