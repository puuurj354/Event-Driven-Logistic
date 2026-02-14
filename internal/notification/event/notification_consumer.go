package event

import (
	"encoding/json"
	"log"

	"github.com/purnama/Event-Driven-Logistic/internal/notification/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type NotificationConsumer struct {
	consumer *broker.Consumer
	svc      service.NotificationService
}

func NewNotificationConsumer(consumer *broker.Consumer, svc service.NotificationService) *NotificationConsumer {
	return &NotificationConsumer{
		consumer: consumer,
		svc:      svc,
	}
}

type eventBinding struct {
	Queue      string
	RoutingKey string
}

func (nc *NotificationConsumer) StartListening() error {
	bindings := []eventBinding{
		{Queue: "notif.order.created", RoutingKey: "order.created"},
		{Queue: "notif.payment.success", RoutingKey: "payment.success"},
		{Queue: "notif.payment.failed", RoutingKey: "payment.failed"},
		{Queue: "notif.stock.reserved", RoutingKey: "stock.reserved"},
		{Queue: "notif.stock.failed", RoutingKey: "stock.failed"},
	}

	for _, b := range bindings {

		err := nc.consumer.Subscribe(
			b.Queue,
			b.RoutingKey,
			nc.handleEvent,
		)
		if err != nil {
			return err
		}
	}

	log.Println("✅ Notification Consumer: listening for ALL events")
	return nil
}
func (nc *NotificationConsumer) handleEvent(event broker.Event) error {

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(event.Payload, &payloadMap); err != nil {
		log.Printf("⚠️ Gagal parse payload, using empty: %v", err)
		payloadMap = map[string]interface{}{}
	}

	orderID := ""
	if oid, ok := payloadMap["order_id"]; ok {
		orderID, _ = oid.(string)
	}
	return nc.svc.ProcessEvent(
		string(event.Type),
		orderID,
		string(event.Payload),
	)
}
