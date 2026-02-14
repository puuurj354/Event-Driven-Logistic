package event

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type DeliveryConsumer struct {
	consumer *broker.Consumer
	svc      service.ShipmentService
}

func NewDeliveryConsumer(consumer *broker.Consumer, svc service.ShipmentService) *DeliveryConsumer {
	return &DeliveryConsumer{
		consumer: consumer,
		svc:      svc,
	}
}
func (dc *DeliveryConsumer) StartListening() error {
	if err := dc.consumer.Subscribe(
		"delivery.payment.success",
		"payment.success",
		dc.handlePaymentSuccess,
	); err != nil {
		return err
	}

	log.Println("✅ Delivery Consumer: listening for payment.success")
	return nil
}

func (dc *DeliveryConsumer) handlePaymentSuccess(event broker.Event) error {
	var payload broker.PaymentSuccessPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("❌ Gagal parse PaymentSuccessPayload: %v", err)
		return err
	}

	log.Printf("📨 Payment.success diterima di Delivery: OrderID=%s, Amount=%.2f",
		payload.OrderID, payload.Amount)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		log.Printf("❌ Format OrderID tidak valid: %s", payload.OrderID)
		return err
	}

	shipment, err := dc.svc.CreateShipment(orderID, "Auto-Assigned")
	if err != nil {
		log.Printf("❌ Gagal membuat shipment: %v", err)
		return err
	}

	log.Printf("✅ Shipment dibuat: ShipmentID=%d, OrderID=%s, Status=%s",
		shipment.ID, payload.OrderID, shipment.Status)

	return nil
}
