package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	// Try to find .env in current directory first
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")     // Current directory (cmd/service-name)
	viper.AddConfigPath("./cmd") // In case running from root

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  Warning: Error reading config file: %v", err)
		log.Printf("📍 Current directory: %s", getCurrentDir())
		log.Printf("💡 Using environment variables instead...")
	} else {
		log.Printf("✅ Config loaded from: %s", viper.ConfigFileUsed())
	}

	config := &Config{
		Database: DatabaseConfig{
			URL: viper.GetString("DB_URL"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: viper.GetString("MQ_URL"),
		},
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
		},
	}

	// Validation
	if config.Database.URL == "" {
		log.Fatal("❌ DB_URL is required but not set")
	}
	if config.RabbitMQ.URL == "" {
		log.Fatal("❌ MQ_URL is required but not set")
	}
	if config.Server.Port == "" {
		log.Fatal("❌ PORT is required but not set")
	}

	return config
}

// LoadConfigFromPath membaca .env dari path spesifik
// Useful untuk testing atau custom setup
func LoadConfigFromPath(configPath string) *Config {
	dir := filepath.Dir(configPath)
	fileName := filepath.Base(configPath)

	viper.SetConfigName(fileName)
	viper.SetConfigType("env")
	viper.AddConfigPath(dir)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file from %s: %v", configPath, err)
	}

	log.Printf("✅ Config loaded from: %s", configPath)

	return &Config{
		Database: DatabaseConfig{
			URL: viper.GetString("DB_URL"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: viper.GetString("MQ_URL"),
		},
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
		},
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}


func (c *Config) PrintConfig() {
	fmt.Println("📋 Current Configuration:")
	fmt.Printf("   Database: %s\n", maskURL(c.Database.URL))
	fmt.Printf("   RabbitMQ: %s\n", maskURL(c.RabbitMQ.URL))
	fmt.Printf("   Server Port: %s\n", c.Server.Port)
}

// maskURL menyembunyikan password dalam URL untuk logging
func maskURL(url string) string {
	if len(url) > 20 {
		return url[:20] + "...****"
	}
	return "****"
}
