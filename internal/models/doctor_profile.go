package models

import (
	"time"
)

type DoctorProfile struct {
	ID                string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID            string    `gorm:"uniqueIndex" json:"userId"`
	User              *User     `gorm:"constraint:OnDelete:CASCADE;" json:"user"`
	Specialization    string    `json:"specialization"`
	LicenseNumber     string    `json:"licenseNumber"`
	ConsultationFees  float64   `json:"consultationFees"`
	AvailabilitySlots string    `gorm:"type:jsonb" json:"availabilitySlots"`
	PhotoURL          *string   `json:"photoUrl,omitempty"`
	IsPending         bool      `gorm:"default:true" json:"isPending"`
	ApprovedBy        *string   `json:"approvedBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Bio               string    `json:"bio"`
	Experience        string    `json:"experience"`
	ClinicName        string    `json:"clinicName"`
	Certifications    string    `json:"certifications"`
	TotalPatients     int       `json:"totalPatients"`
	Rating            float64   `json:"rating"`
	Reviews           []Review  `gorm:"foreignKey:DoctorID" json:"reviews"`
}
