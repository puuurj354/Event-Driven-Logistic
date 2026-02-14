package event

import (
	"log"

	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

type PaymentPublisher struct {
	publisher *broker.Publisher
}

func NewPaymentPublisher(publisher *broker.Publisher) *PaymentPublisher {
	return &PaymentPublisher{publisher: publisher}
}

func (p *PaymentPublisher) PublishPaymentSuccess(payload broker.PaymentSuccessPayload) error {
	log.Printf("📤 Publishing payment.success: OrderID=%s, PaymentID=%d",
		payload.OrderID, payload.PaymentID)

	err := p.publisher.PublishEvent(broker.PaymentSuccess, payload)
	if err != nil {
		log.Printf("❌ Gagal publish payment.success: %v", err)
		return err
	}

	log.Printf("✅ Event payment.success berhasil dipublish")
	return nil
}

func (p *PaymentPublisher) PublishPaymentFailed(payload broker.PaymentFailedPayload) error {
	log.Printf("📤 Publishing payment.failed: OrderID=%s, Reason=%s",
		payload.OrderID, payload.Reason)

	err := p.publisher.PublishEvent(broker.PaymentFailed, payload)
	if err != nil {
		log.Printf("❌ Gagal publish payment.failed: %v", err)
		return err
	}

	return nil
}
