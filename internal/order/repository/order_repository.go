package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *Order) error

	GetByID(id uuid.UUID) (*Order, error)

	GetByCustomerID(customerID string) ([]Order, error)

	UpdateStatus(id uuid.UUID, status OrderStatus) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uuid.UUID) (*Order, error) {
	var order Order
	err := r.db.First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByCustomerID(customerID string) ([]Order, error) {
	var orders []Order
	err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status OrderStatus) error {
	return r.db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}
