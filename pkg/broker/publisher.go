package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) (*Publisher, error) {
	err := channel.ExchangeDeclare(
		ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal declare exchange '%s': %w", ExchangeName, err)
	}

	log.Printf("✅ Exchange '%s' (topic) berhasil di-declare", ExchangeName)

	return &Publisher{channel: channel}, nil
}

func (p *Publisher) Publish(event *Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("gagal serialize event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		ExchangeName,
		string(event.Type),
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    event.Timestamp,
			Type:         string(event.Type),
		},
	)
	if err != nil {
		return fmt.Errorf("gagal publish event '%s': %w", event.Type, err)
	}

	log.Printf("📤 Event published: type=%s, exchange=%s", event.Type, ExchangeName)

	return nil
}

func (p *Publisher) PublishEvent(eventType EventType, payload interface{}) error {
	event, err := NewEvent(eventType, payload)
	if err != nil {
		return fmt.Errorf("gagal membuat event: %w", err)
	}

	return p.Publish(event)
}
