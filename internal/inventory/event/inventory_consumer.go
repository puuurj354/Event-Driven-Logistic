package event

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type InventoryConsumer struct {
	consumer  *broker.Consumer
	publisher *broker.Publisher
	svc       service.InventoryService
}

func NewInventoryConsumer(
	consumer *broker.Consumer,
	publisher *broker.Publisher,
	svc service.InventoryService,
) *InventoryConsumer {
	return &InventoryConsumer{
		consumer:  consumer,
		publisher: publisher,
		svc:       svc,
	}
}

func (ic *InventoryConsumer) StartListening() error {

	if err := ic.consumer.Subscribe(
		"inventory.order.created",
		"order.created",
		ic.handleOrderCreated,
	); err != nil {
		return err
	}

	log.Println("✅ Inventory Consumer: listening for order.created")
	return nil
}

func (ic *InventoryConsumer) handleOrderCreated(event broker.Event) error {

	var payload broker.OrderCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("❌ Gagal parse OrderCreatedPayload: %v", err)
		return err
	}

	log.Printf("📨 Order.created diterima di Inventory: OrderID=%s, Item=%s, Qty=%d",
		payload.OrderID, payload.ItemName, payload.Quantity)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		log.Printf("❌ Format OrderID tidak valid: %s", payload.OrderID)
		return err
	}

	reservation, err := ic.svc.ReserveStock(orderID, payload.ItemName, payload.Quantity)
	if err != nil {

		log.Printf("⚠️ Stok gagal direservasi: %v", err)

		failPayload := broker.StockFailedPayload{
			OrderID:  payload.OrderID,
			ItemName: payload.ItemName,
			Reason:   err.Error(),
		}

		if pubErr := ic.publisher.PublishEvent(broker.StockFailed, failPayload); pubErr != nil {
			log.Printf("❌ Gagal publish stock.failed: %v", pubErr)
		}

		return nil
	}

	successPayload := broker.StockReservedPayload{
		OrderID:       payload.OrderID,
		ProductID:     reservation.ProductID,
		Quantity:      reservation.Quantity,
		ReservationID: reservation.ID,
	}

	if err := ic.publisher.PublishEvent(broker.StockReserved, successPayload); err != nil {
		log.Printf("❌ Gagal publish stock.reserved: %v", err)
		return err
	}

	log.Printf("✅ Stok berhasil direservasi dan event dipublish: OrderID=%s, Qty=%d",
		payload.OrderID, payload.Quantity)

	return nil
}
