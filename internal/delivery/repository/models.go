package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShipmentStatus string

const (
	ShipmentStatusPickingUp ShipmentStatus = "PICKING_UP"
	ShipmentStatusOnTheWay  ShipmentStatus = "ON_THE_WAY"
	ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
)

func (s ShipmentStatus) IsValid() bool {
	switch s {
	case ShipmentStatusPickingUp, ShipmentStatusOnTheWay, ShipmentStatusDelivered:
		return true
	}
	return false
}

type Shipment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderID     uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	CourierName string         `gorm:"not null" json:"courier_name"`
	CurrentLat  float64        `gorm:"type:decimal(10,8)" json:"current_lat"`
	CurrentLong float64        `gorm:"type:decimal(11,8)" json:"current_long"`
	Status      ShipmentStatus `gorm:"type:varchar(20);default:PICKING_UP;not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"` 
}

func (s *Shipment) BeforeCreate(tx *gorm.DB) (err error) {
	if !s.Status.IsValid() {
		return fmt.Errorf("invalid shipment status: %s", s.Status)
	}
	if s.CurrentLat < -90 || s.CurrentLat > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", s.CurrentLat)
	}
	if s.CurrentLong < -180 || s.CurrentLong > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", s.CurrentLong)
	}

	return nil
}

func (s *Shipment) BeforeUpdate(tx *gorm.DB) (err error) {
	if !s.Status.IsValid() {
		return fmt.Errorf("invalid shipment status: %s", s.Status)
	}

	return nil
}
