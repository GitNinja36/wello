package models

import (
	"time"
)

type Order struct {
	ID            string        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        string        `json:"userId"`
	User          User          `gorm:"foreignKey:UserID"`
	Items         string        `gorm:"type:jsonb" json:"items"`
	Amount        float64       `json:"amount"`
	Status        OrderStatus   `gorm:"type:text;default:'PENDING'" json:"status"`
	PaymentMethod PaymentMethod `gorm:"type:text;default:'ONLINE'" json:"paymentMethod"`
	PaymentID     *string       `json:"paymentId,omitempty"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}
