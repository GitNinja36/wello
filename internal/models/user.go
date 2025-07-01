package models

import (
	"time"
)

type User struct {
	ID                string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string         `json:"name"`
	Email             string         `gorm:"uniqueIndex" json:"email"`
	Phone             string         `gorm:"uniqueIndex" json:"phone"`
	Role              Role           `gorm:"type:text;default:'PATIENT'" json:"role"`
	Age               int            `json:"age"`
	Gender            string         `json:"gender"`
	Bio               string         `json:"bio"`
	Address           string         `json:"address"`
	Password          string         `json:"-"`
	Verified          bool           `gorm:"default:false" json:"verified"`
	IsApproved        bool           `gorm:"default:false" json:"isApproved"`
	RequestedAsDoctor bool           `gorm:"default:false" json:"requestedAsDoctor"`
	PhotoURL          *string        `json:"photoUrl,omitempty"`
	Appointments      []Appointment  `gorm:"foreignKey:PatientID" json:"appointments"`
	Orders            []Order        `json:"orders"`
	AdminProfile      *AdminProfile  `gorm:"foreignKey:UserID" json:"adminProfile"`
	DoctorProfile     *DoctorProfile `gorm:"foreignKey:UserID" json:"doctorProfile,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}
