package broker

import (
	"encoding/json"
	"time"
)

const ExchangeName = "logistic.events"

type EventType string

const (
	OrderCreated EventType = "order.created"

	OrderCancelled EventType = "order.cancelled"

	PaymentSuccess EventType = "payment.success"

	PaymentFailed EventType = "payment.failed"

	StockReserved EventType = "stock.reserved"

	StockFailed EventType = "stock.failed"

	ShipmentCreated EventType = "shipment.created"

	ShipmentStatusUpdated EventType = "shipment.status_updated"
)

type Event struct {
	Type      EventType       `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	ItemName   string  `json:"item_name"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type PaymentSuccessPayload struct {
	OrderID   string  `json:"order_id"`
	PaymentID uint    `json:"payment_id"`
	Amount    float64 `json:"amount"`
}

type PaymentFailedPayload struct {
	OrderID   string `json:"order_id"`
	PaymentID uint   `json:"payment_id"`
	Reason    string `json:"reason"`
}

type StockReservedPayload struct {
	OrderID       string `json:"order_id"`
	ProductID     uint   `json:"product_id"`
	Quantity      int    `json:"quantity"`
	ReservationID uint   `json:"reservation_id"`
}

type StockFailedPayload struct {
	OrderID  string `json:"order_id"`
	ItemName string `json:"item_name"`
	Reason   string `json:"reason"`
}

func NewEvent(eventType EventType, payload interface{}) (*Event, error) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   json.RawMessage(payloadBytes),
	}, nil
}
