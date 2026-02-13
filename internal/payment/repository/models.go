package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending, PaymentStatusCompleted, PaymentStatusFailed, PaymentStatusRefunded:
		return true
	}
	return false
}

type Payment struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	OrderID       uuid.UUID     `gorm:"type:uuid;not null;index" json:"order_id"`
	Amount        float64       `gorm:"not null;check:amount >= 0" json:"amount"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);default:PENDING;not null" json:"status"`
	PaidAt        *time.Time    `json:"paid_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if !p.PaymentStatus.IsValid() {
		return fmt.Errorf("invalid payment status: %s", p.PaymentStatus)
	}
	if p.Amount < 0 {
		return fmt.Errorf("payment amount cannot be negative: %f", p.Amount)
	}

	return nil
}

func (p *Payment) BeforeUpdate(tx *gorm.DB) (err error) {
	if !p.PaymentStatus.IsValid() {
		return fmt.Errorf("invalid payment status: %s", p.PaymentStatus)
	}

	return nil
}
