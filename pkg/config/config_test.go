package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	
	os.Setenv("DB_URL", "postgres://postgres:postgres@localhost:5432/logistic_db?sslmode=disable")
	os.Setenv("MQ_URL", "amqp://guest:guest@localhost:5672/")
	os.Setenv("PORT_ORDER", "8080")

	cfg := LoadConfig()

	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/logistic_db?sslmode=disable", cfg.Database.URL)
	assert.Equal(t, "amqp://guest:guest@localhost:5672/", cfg.RabbitMQ.URL)
	assert.Equal(t, "8080", cfg.Server.Port)
}
