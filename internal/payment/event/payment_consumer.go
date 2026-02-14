package event

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type PaymentConsumer struct {
	consumer *broker.Consumer
	svc      service.PaymentService
}

func NewPaymentConsumer(consumer *broker.Consumer, svc service.PaymentService) *PaymentConsumer {
	return &PaymentConsumer{
		consumer: consumer,
		svc:      svc,
	}
}

func (pc *PaymentConsumer) StartListening() error {
	if err := pc.consumer.Subscribe(
		"payment.order.created",
		"order.created",
		pc.handleOrderCreated,
	); err != nil {
		return err
	}

	log.Println("✅ Payment Consumer: listening for order.created")
	return nil
}

func (pc *PaymentConsumer) handleOrderCreated(event broker.Event) error {
	var payload broker.OrderCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("❌ Gagal parse OrderCreatedPayload: %v", err)
		return err
	}

	log.Printf("📨 Order.created diterima di Payment: OrderID=%s, Amount=%.2f",
		payload.OrderID, payload.TotalPrice)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		log.Printf("❌ Format OrderID tidak valid: %s", payload.OrderID)
		return err
	}

	payment, err := pc.svc.CreatePayment(orderID, payload.TotalPrice)
	if err != nil {
		log.Printf("❌ Gagal membuat payment: %v", err)
		return err
	}

	log.Printf("✅ Payment PENDING dibuat: PaymentID=%d, OrderID=%s, Amount=%.2f",
		payment.ID, payload.OrderID, payload.TotalPrice)

	return nil
}
