package repository

import (
	"gorm.io/gorm"
)


type PaymentRepository interface {

	CreatePayment(payment *Payment) error

	GetPaymentByOrderID(orderID string) (*Payment, error)
	GetPaymentByID(id uint) (*Payment, error)

	UpdatePayment(payment *Payment) error
}


type paymentRepository struct {
	db *gorm.DB 
}


func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db} 
}

func (r *paymentRepository) CreatePayment(payment *Payment) error {
	return r.db.Create(payment).Error 
}

func (r *paymentRepository) GetPaymentByOrderID(orderID string) (*Payment, error) {
	var payment Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error 
	if err != nil {
		return nil, err 
	}
	return &payment, nil
}

func (r *paymentRepository) GetPaymentByID(id uint) (*Payment, error) {
	var payment Payment
	err := r.db.First(&payment, id).Error 
	if err != nil {
		return nil, err 
	}
	return &payment, nil 
}

func (r *paymentRepository) UpdatePayment(payment *Payment) error {
	return r.db.Save(payment).Error 
}
