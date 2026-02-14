package service

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/order/repository"
)

type OrderService interface {
	CreateOrder(req CreateOrderRequest) (*repository.Order, error)

	GetOrderByID(id uuid.UUID) (*repository.Order, error)

	GetOrdersByCustomerID(customerID string) ([]repository.Order, error)

	UpdateOrderStatus(id uuid.UUID, status repository.OrderStatus) error
}

type CreateOrderRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	ItemName   string  `json:"item_name" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,gt=0"`
	TotalPrice float64 `json:"total_price" binding:"required,gt=0"`
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(req CreateOrderRequest) (*repository.Order, error) {

	if req.Quantity <= 0 {
		return nil, errors.New("quantity harus lebih dari 0")
	}

	if req.TotalPrice <= 0 {
		return nil, errors.New("total price harus lebih dari 0")
	}

	if req.CustomerID == "" {
		return nil, errors.New("customer_id tidak boleh kosong")
	}

	order := &repository.Order{
		CustomerID: req.CustomerID,
		ItemName:   req.ItemName,
		Quantity:   req.Quantity,
		TotalPrice: req.TotalPrice,
		Status:     repository.PENDING,
	}

	if err := s.repo.Create(order); err != nil {
		log.Printf("❌ Gagal membuat order: %v", err)
		return nil, err
	}

	log.Printf("✅ Order berhasil dibuat: ID=%s, Customer=%s, Item=%s",
		order.ID, order.CustomerID, order.ItemName)

	return order, nil
}

func (s *orderService) GetOrderByID(id uuid.UUID) (*repository.Order, error) {
	order, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("❌ Order tidak ditemukan: ID=%s, error=%v", id, err)
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetOrdersByCustomerID(customerID string) ([]repository.Order, error) {
	if customerID == "" {
		return nil, errors.New("customer_id tidak boleh kosong")
	}

	orders, err := s.repo.GetByCustomerID(customerID)
	if err != nil {
		log.Printf("❌ Gagal mengambil order untuk customer: %s, error=%v", customerID, err)
		return nil, err
	}
	return orders, nil
}

func (s *orderService) UpdateOrderStatus(id uuid.UUID, status repository.OrderStatus) error {

	if !status.IsValid() {
		return errors.New("status order tidak valid: " + string(status))
	}

	currentOrder, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("order tidak ditemukan: " + id.String())
	}

	if currentOrder.Status == repository.CANCELLED {
		return errors.New("tidak dapat mengubah status order yang sudah CANCELLED")
	}

	if currentOrder.Status == repository.SHIPPED && status == repository.PENDING {
		return errors.New("tidak dapat mengembalikan status SHIPPED ke PENDING")
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		log.Printf("❌ Gagal update status order: ID=%s, error=%v", id, err)
		return err
	}

	log.Printf("✅ Status order diperbarui: ID=%s, %s → %s",
		id, currentOrder.Status, status)

	return nil 
}
