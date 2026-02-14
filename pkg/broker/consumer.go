package broker

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventHandler func(event Event) error

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(channel *amqp.Channel) (*Consumer, error) {

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

	return &Consumer{channel: channel}, nil
}

func (c *Consumer) Subscribe(queueName, routingKey string, handler EventHandler) error {

	q, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("gagal declare queue '%s': %w", queueName, err)
	}

	err = c.channel.QueueBind(
		q.Name,
		routingKey,
		ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("gagal bind queue '%s' ke routing key '%s': %w", queueName, routingKey, err)
	}

	log.Printf("✅ Queue '%s' bound to exchange '%s' with key '%s'",
		queueName, ExchangeName, routingKey)

	msgs, err := c.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("gagal start consuming dari queue '%s': %w", queueName, err)
	}

	go func() {
		log.Printf("👂 Consumer listening on queue '%s' for '%s'...", queueName, routingKey)

		for msg := range msgs {
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Gagal parse event dari queue '%s': %v", queueName, err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("📨 Event diterima: type=%s, queue=%s", event.Type, queueName)

			if err := handler(event); err != nil {
				log.Printf("❌ Gagal proses event '%s': %v", event.Type, err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
			log.Printf("✅ Event '%s' berhasil diproses dari queue '%s'", event.Type, queueName)
		}

		log.Printf("⚠️ Consumer stopped for queue '%s'", queueName)
	}()

	return nil
}
