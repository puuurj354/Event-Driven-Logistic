package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	URL string
}

type RabbitMQConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() *Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			URL: viper.GetString("DB_URL"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: viper.GetString("MQ_URL"),
		},
		Server: ServerConfig{
			Port: viper.GetString("PORT_ORDER"), // This can be overridden per service
		},
	}

	return config
}
