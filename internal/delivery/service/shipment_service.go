package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/repository"
)

type ShipmentService interface {
	CreateShipment(orderID uuid.UUID, courierName string) (*repository.Shipment, error)

	GetShipmentByOrderID(orderID string) (*repository.Shipment, error)

	UpdateShipmentStatus(shipmentID uint, status repository.ShipmentStatus) error

	UpdateLocation(shipmentID uint, lat, long float64) error
}

type shipmentService struct {
	repo repository.ShipmentRepository
}

func NewShipmentService(repo repository.ShipmentRepository) ShipmentService {
	return &shipmentService{repo: repo}
}

func (s *shipmentService) CreateShipment(orderID uuid.UUID, courierName string) (*repository.Shipment, error) {

	if courierName == "" {
		return nil, errors.New("nama kurir tidak boleh kosong")
	}

	// Bangun model Shipment
	shipment := &repository.Shipment{
		OrderID:     orderID,
		CourierName: courierName,
		Status:      repository.ShipmentStatusPickingUp,
		CurrentLat:  0,
		CurrentLong: 0, // Koordinat awal (belum ada tracking)
	}

	// Simpan ke database
	if err := s.repo.CreateShipment(shipment); err != nil {
		log.Printf("❌ Gagal membuat shipment: OrderID=%s, error=%v", orderID, err)
		return nil, err
	}

	log.Printf("✅ Shipment dibuat: ID=%d, OrderID=%s, Kurir=%s, Status=PICKING_UP",
		shipment.ID, orderID, courierName)

	return shipment, nil
}

func (s *shipmentService) GetShipmentByOrderID(orderID string) (*repository.Shipment, error) {

	if orderID == "" {
		return nil, errors.New("order_id tidak boleh kosong")
	}

	shipment, err := s.repo.GetShipmentByOrderID(orderID)
	if err != nil {
		log.Printf("❌ Shipment tidak ditemukan: OrderID=%s, error=%v", orderID, err)
		return nil, err
	}
	return shipment, nil
}

func (s *shipmentService) UpdateShipmentStatus(shipmentID uint, status repository.ShipmentStatus) error {

	if !status.IsValid() {
		return errors.New("status shipment tidak valid: " + string(status))
	}

	shipment, err := s.repo.GetShipmentByID(shipmentID)
	if err != nil {
		return fmt.Errorf("shipment tidak ditemukan: ID=%d, error: %w", shipmentID, err)
	}

	if !isValidStatusTransition(shipment.Status, status) {
		return fmt.Errorf("transisi status tidak valid: %s → %s", shipment.Status, status)
	}

	if err := s.repo.UpdateStatus(shipmentID, status); err != nil {
		log.Printf("❌ Gagal update status shipment: ID=%d, error=%v", shipmentID, err)
		return err
	}

	log.Printf("✅ Status shipment diperbarui: ID=%d, %s → %s",
		shipmentID, shipment.Status, status)

	return nil
}

func (s *shipmentService) UpdateLocation(shipmentID uint, lat, long float64) error {

	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude tidak valid: %f (harus antara -90 dan 90)", lat)
	}

	if long < -180 || long > 180 {
		return fmt.Errorf("longitude tidak valid: %f (harus antara -180 dan 180)", long)
	}

	if err := s.repo.UpdateLocation(shipmentID, lat, long); err != nil {
		log.Printf("❌ Gagal update lokasi shipment: ID=%d, error=%v", shipmentID, err)
		return err
	}

	log.Printf("📍 Lokasi kurir diperbarui: ShipmentID=%d, Lat=%f, Long=%f",
		shipmentID, lat, long)

	return nil
}

func isValidStatusTransition(from, to repository.ShipmentStatus) bool {
	validTransitions := map[repository.ShipmentStatus][]repository.ShipmentStatus{
		repository.ShipmentStatusPickingUp: {repository.ShipmentStatusOnTheWay},
		repository.ShipmentStatusOnTheWay:  {repository.ShipmentStatusDelivered},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
