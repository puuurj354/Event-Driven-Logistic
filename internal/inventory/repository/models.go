package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductStatusReleased  ProductStatus = "RELEASED"
	ProductStatusReserved  ProductStatus = "RESERVED"
	ProductStatusConfirmed ProductStatus = "CONFIRMED"
)

func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusReleased, ProductStatusReserved, ProductStatusConfirmed:
		return true
	}
	return false
}

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;uniqueIndex" json:"name"`
	Stock     int       `gorm:"default:0;check:stock >= 0" json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type StockReservation struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	OrderID   uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	ProductID uint          `gorm:"not null;index" json:"product_id"`
	Quantity  int           `gorm:"not null;check:quantity > 0" json:"quantity"`
	Status    ProductStatus `gorm:"type:varchar(20);default:RESERVED;not null" json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Stock < 0 {
		return fmt.Errorf("product stock cannot be negative: %d", p.Stock)
	}

	return nil
}

func (sr *StockReservation) BeforeCreate(tx *gorm.DB) (err error) {
	if !sr.Status.IsValid() {
		return fmt.Errorf("invalid stock reservation status: %s", sr.Status)
	}
	if sr.Quantity <= 0 {
		return fmt.Errorf("reservation quantity must be positive: %d", sr.Quantity)
	}

	return nil
}

func (sr *StockReservation) BeforeUpdate(tx *gorm.DB) (err error) {
	if !sr.Status.IsValid() {
		return fmt.Errorf("invalid stock reservation status: %s", sr.Status)
	}

	return nil
}
