package service

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/repository"
)

type PaymentService interface {
	CreatePayment(orderID uuid.UUID, amount float64) (*repository.Payment, error)

	GetPaymentByOrderID(orderID string) (*repository.Payment, error)

	ConfirmPayment(paymentID uint) (*repository.Payment, error)

	FailPayment(paymentID uint) error
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) CreatePayment(orderID uuid.UUID, amount float64) (*repository.Payment, error) {
	if amount <= 0 {
		return nil, errors.New("jumlah pembayaran harus lebih dari 0")
	}

	payment := &repository.Payment{
		OrderID:       orderID,
		Amount:        amount,
		PaymentStatus: repository.PaymentStatusPending,
	}

	if err := s.repo.CreatePayment(payment); err != nil {
		log.Printf("❌ Gagal membuat payment: OrderID=%s, error=%v", orderID, err)
		return nil, err
	}

	log.Printf("✅ Payment dibuat: ID=%d, OrderID=%s, Amount=%.2f, Status=PENDING",
		payment.ID, orderID, amount)

	return payment, nil
}

func (s *paymentService) GetPaymentByOrderID(orderID string) (*repository.Payment, error) {
	if orderID == "" {
		return nil, errors.New("order_id tidak boleh kosong")
	}

	payment, err := s.repo.GetPaymentByOrderID(orderID)
	if err != nil {
		log.Printf("❌ Payment tidak ditemukan untuk order: %s, error=%v", orderID, err)
		return nil, err
	}
	return payment, nil
}

func (s *paymentService) ConfirmPayment(paymentID uint) (*repository.Payment, error) {
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return nil, errors.New("payment tidak ditemukan")
	}

	if payment.PaymentStatus != repository.PaymentStatusPending {
		return nil, errors.New("hanya payment dengan status PENDING yang bisa dikonfirmasi")
	}

	now := time.Now()
	payment.PaymentStatus = repository.PaymentStatusCompleted
	payment.PaidAt = &now
	if err := s.repo.UpdatePayment(payment); err != nil {
		log.Printf("❌ Gagal konfirmasi payment: ID=%d, error=%v", paymentID, err)
		return nil, err
	}

	log.Printf("✅ Payment dikonfirmasi: ID=%d, OrderID=%s, Status=COMPLETED",
		payment.ID, payment.OrderID)

	return payment, nil
}

func (s *paymentService) FailPayment(paymentID uint) error {
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return errors.New("payment tidak ditemukan")
	}

	if payment.PaymentStatus != repository.PaymentStatusPending {
		return errors.New("hanya payment dengan status PENDING yang bisa di-fail")
	}
	payment.PaymentStatus = repository.PaymentStatusFailed
	if err := s.repo.UpdatePayment(payment); err != nil {
		log.Printf("❌ Gagal mengubah payment ke FAILED: ID=%d, error=%v", paymentID, err)
		return err
	}

	log.Printf("⚠️ Payment gagal: ID=%d, OrderID=%s, Status=FAILED",
		payment.ID, payment.OrderID)

	return nil
}
