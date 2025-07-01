package models

import (
	"time"
)

type Appointment struct {
	ID              string            `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PatientID       string            `json:"patientId"`
	Patient         User              `gorm:"foreignKey:PatientID"`
	DoctorProfileID string            `json:"doctorProfileId"`
	DoctorProfile   DoctorProfile     `gorm:"foreignKey:DoctorProfileID" json:"DoctorProfile"`
	Mode            AppointmentMode   `gorm:"type:text;default:'ONLINE'" json:"mode"`
	Status          AppointmentStatus `gorm:"type:text;default:'PENDING'" json:"status"`
	ScheduledAt     time.Time         `json:"scheduledAt"`
	MeetingLink     *string           `json:"meetingLink,omitempty"`
	Location        *string           `json:"location,omitempty"`
	FeePaid         bool              `gorm:"default:false" json:"feePaid"`
	Summary         *string           `json:"summary,omitempty"`
	Rating          *int              `json:"rating"`
	Review          *string           `json:"review"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}
