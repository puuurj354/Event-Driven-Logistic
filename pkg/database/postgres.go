package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10) // keep 10 connections idle in the pool

	sqlDB.SetMaxOpenConns(100) // allow up to 100 open connections to the database

	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection initialized successfully")

	return db
}
