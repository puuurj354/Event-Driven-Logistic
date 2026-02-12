package repository

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OrderID       uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	PaymentStatus string    `json:"payment_status"`
	PaidAt        time.Time `json:"paid_at"`
}
