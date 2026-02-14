package event

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/order/repository"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type OrderConsumer struct {
	consumer *broker.Consumer
	repo     repository.OrderRepository
}

func NewOrderConsumer(consumer *broker.Consumer, repo repository.OrderRepository) *OrderConsumer {
	return &OrderConsumer{
		consumer: consumer,
		repo:     repo,
	}
}

func (oc *OrderConsumer) StartListening() error {

	if err := oc.consumer.Subscribe(
		"order.payment.success",
		"payment.success",
		oc.handlePaymentSuccess,
	); err != nil {
		return err
	}

	if err := oc.consumer.Subscribe(
		"order.payment.failed",
		"payment.failed",
		oc.handlePaymentFailed,
	); err != nil {
		return err
	}

	if err := oc.consumer.Subscribe(
		"order.stock.failed",
		"stock.failed",
		oc.handleStockFailed,
	); err != nil {
		return err
	}

	log.Println("✅ Order Consumer: listening for payment.success, payment.failed, stock.failed")
	return nil
}

func (oc *OrderConsumer) handlePaymentSuccess(event broker.Event) error {

	var payload broker.PaymentSuccessPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("❌ Gagal parse PaymentSuccessPayload: %v", err)
		return err
	}

	log.Printf("📨 Payment sukses diterima: OrderID=%s, Amount=%.2f",
		payload.OrderID, payload.Amount)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		log.Printf("❌ Format OrderID tidak valid: %s", payload.OrderID)
		return err
	}

	if err := oc.repo.UpdateStatus(orderID, repository.PAID); err != nil {
		log.Printf("❌ Gagal update order status ke PAID: %v", err)
		return err
	}

	log.Printf("✅ Order %s status diperbarui: PENDING → PAID", payload.OrderID)
	return nil
}

func (oc *OrderConsumer) handlePaymentFailed(event broker.Event) error {

	var payload broker.PaymentFailedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("❌ Gagal parse PaymentFailedPayload: %v", err)
		return err
	}

	log.Printf("📨 Payment gagal diterima: OrderID=%s, Reason=%s",
		payload.OrderID, payload.Reason)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return err
	}

	if err := oc.repo.UpdateStatus(orderID, repository.CANCELLED); err != nil {
		log.Printf("❌ Gagal update order status ke CANCELLED: %v", err)
		return err
	}

	log.Printf("✅ Order %s status diperbarui: PENDING → CANCELLED (payment failed)", payload.OrderID)
	return nil
}
func (oc *OrderConsumer) handleStockFailed(event broker.Event) error {

	var payload broker.StockFailedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("❌ Gagal parse StockFailedPayload: %v", err)
		return err
	}

	log.Printf("📨 Stock gagal diterima: OrderID=%s, Item=%s, Reason=%s",
		payload.OrderID, payload.ItemName, payload.Reason)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return err
	}

	if err := oc.repo.UpdateStatus(orderID, repository.CANCELLED); err != nil {
		log.Printf("❌ Gagal update order status ke CANCELLED: %v", err)
		return err
	}

	log.Printf("✅ Order %s status diperbarui: PENDING → CANCELLED (stock failed)", payload.OrderID) // Log
	return nil                                                                                      // Sukses
}
