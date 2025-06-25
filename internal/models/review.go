package models

import "time"

type Review struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DoctorID  string    `json:"doctorId"`
	PatientID string    `json:"patientId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
