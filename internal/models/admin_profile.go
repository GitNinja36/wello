package models

import "time"

type AdminProfile struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string    `gorm:"uniqueIndex" json:"userId"`
	User       User      `gorm:"constraint:OnDelete:CASCADE"`
	Position   string    `json:"position"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
