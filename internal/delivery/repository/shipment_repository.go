package repository

import (
	"gorm.io/gorm"
)

type ShipmentRepository interface {
	CreateShipment(shipment *Shipment) error
	GetShipmentByOrderID(orderID string) (*Shipment, error)
	GetShipmentByID(id uint) (*Shipment, error)
	UpdateShipment(shipment *Shipment) error

	UpdateLocation(shipmentID uint, lat, long float64) error

	UpdateStatus(shipmentID uint, status ShipmentStatus) error
}

type shipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) CreateShipment(shipment *Shipment) error {
	return r.db.Create(shipment).Error
}

func (r *shipmentRepository) GetShipmentByOrderID(orderID string) (*Shipment, error) {
	var shipment Shipment
	err := r.db.Where("order_id = ?", orderID).First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) GetShipmentByID(id uint) (*Shipment, error) {
	var shipment Shipment
	err := r.db.First(&shipment, id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *shipmentRepository) UpdateShipment(shipment *Shipment) error {
	return r.db.Save(shipment).Error
}

func (r *shipmentRepository) UpdateLocation(shipmentID uint, lat, long float64) error {
	return r.db.Model(&Shipment{}).Where("id = ?", shipmentID).Updates(map[string]interface{}{
		"current_lat":  lat,
		"current_long": long,
	}).Error
}

func (r *shipmentRepository) UpdateStatus(shipmentID uint, status ShipmentStatus) error {
	return r.db.Model(&Shipment{}).Where("id = ?", shipmentID).Update("status", status).Error
}
