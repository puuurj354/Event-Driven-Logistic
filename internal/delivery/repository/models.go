package repository

import (
	"github.com/google/uuid"
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
	OrderID     uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	CourierName string         `json:"courier_name"`
	CurrentLat  float64        `json:"current_lat"`
	CurrentLong float64        `json:"current_long"`
	Status      ShipmentStatus `json:"status"`
}

