package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	PENDING   OrderStatus = "PENDING"
	PAID      OrderStatus = "PAID"
	SHIPPED   OrderStatus = "SHIPPED"
	CANCELLED OrderStatus = "CANCELLED"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case PENDING, PAID, SHIPPED, CANCELLED:
		return true
	}
	return false
}

type Order struct {
	ID         uuid.UUID   `gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	CustomerID string      `gorm:"not null"`
	ItemName   string      `gorm:"not null"`
	Quantity   int         `gorm:"not null"`
	TotalPrice float64     `gorm:"not null"`
	Status     OrderStatus `gorm:"type:varchar(20);default:PENDING;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()

	if !o.Status.IsValid() {
		return fmt.Errorf("invalid order status: %s", o.Status)
	}

	return nil
}
