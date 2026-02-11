package broker

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectRabbitMQ(t *testing.T) {

	os.Setenv("MQ_URL", "amqp://guest:guest@localhost:5672/")


	conn := ConnectRabbitMQ()
	assert.NotNil(t, conn, "RabbitMQ connection should not be nil")


	defer conn.Close()
}

func TestCreateChannel(t *testing.T) {

	os.Setenv("MQ_URL", "amqp://guest:guest@localhost:5672/")


	conn := ConnectRabbitMQ()
	defer conn.Close()


	ch := CreateChannel(conn)
	assert.NotNil(t, ch, "RabbitMQ channel should not be nil")

	
	defer ch.Close()
}
