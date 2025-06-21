package models

import (
	"time"
)

type MedicalCheck struct {
	ID             string      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AppointmentID  string      `gorm:"uniqueIndex" json:"appointmentId"`
	Appointment    Appointment `gorm:"foreignKey:AppointmentID"`
	Type           TestType    `gorm:"type:text;default:'BLOOD'" json:"type"`
	Location       string      `json:"location"`
	TeamAssigned   *string     `json:"teamAssigned,omitempty"`
	Status         TestStatus  `gorm:"type:text;default:'PENDING'" json:"status"`
	ReportUploaded bool        `gorm:"default:false" json:"reportUploaded"`
	ReportUrl      *string     `json:"reportUrl,omitempty"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}
