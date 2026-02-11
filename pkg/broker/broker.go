package broker

import (
	"log"

	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() *amqp.Connection {
	cfg := config.LoadConfig()
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	log.Println("RabbitMQ connected successfully")
	return conn
}

func CreateChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	return ch
}
