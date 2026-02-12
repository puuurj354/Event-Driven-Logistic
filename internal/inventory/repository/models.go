package repository

import (
	"github.com/google/uuid"
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
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Stock int    `gorm:"default:0" json:"stock"`
}
type StockReservation struct {
	ID      uint          `gorm:"primaryKey" json:"id"`
	OrderID uuid.UUID     `gorm:"type:uuid" json:"order_id"`
	Status  ProductStatus `json:"status"`
}
