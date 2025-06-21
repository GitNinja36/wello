package models

import (
	"time"
)

type DoctorProfile struct {
	ID                string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID            string    `gorm:"uniqueIndex" json:"userId"`
	User              User      `gorm:"constraint:OnDelete:CASCADE"`
	Specialization    string    `json:"specialization"`
	LicenseNumber     string    `json:"licenseNumber"`
	ConsultationFees  float64   `json:"consultationFees"`
	AvailabilitySlots string    `gorm:"type:jsonb" json:"availabilitySlots"` // JSON array
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
