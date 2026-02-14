package event

import (
	"log"

	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type OrderPublisher struct {
	publisher *broker.Publisher
}

func NewOrderPublisher(publisher *broker.Publisher) *OrderPublisher {
	return &OrderPublisher{publisher: publisher}
}

func (p *OrderPublisher) PublishOrderCreated(payload broker.OrderCreatedPayload) error {
	log.Printf("📤 Publishing order.created: OrderID=%s, Item=%s, Qty=%d",
		payload.OrderID, payload.ItemName, payload.Quantity)

	err := p.publisher.PublishEvent(broker.OrderCreated, payload)
	if err != nil {
		log.Printf("❌ Gagal publish order.created: %v", err)
		return err
	}

	log.Printf("✅ Event order.created berhasil dipublish: OrderID=%s", payload.OrderID)
	return nil
}

func (p *OrderPublisher) PublishOrderCancelled(orderID string) error {
	log.Printf("📤 Publishing order.cancelled: OrderID=%s", orderID)

	payload := map[string]string{"order_id": orderID}

	err := p.publisher.PublishEvent(broker.OrderCancelled, payload)
	if err != nil {
		log.Printf("❌ Gagal publish order.cancelled: %v", err)
		return err
	}

	return nil
}
