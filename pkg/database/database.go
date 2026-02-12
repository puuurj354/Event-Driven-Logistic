package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB membuat koneksi ke PostgreSQL database
// Setiap service akan memanggil fungsi ini dengan DB_URL masing-masing
func ConnectDB(dbURL string) (*gorm.DB, error) {

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set ke logger.Silent untuk production
	}

	// Buka koneksi
	db, err := gorm.Open(postgres.Open(dbURL), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}


	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB.SetMaxIdleConns(10)   
	sqlDB.SetMaxOpenConns(100)  
	sqlDB.SetConnMaxLifetime(0) 

	log.Printf("✅ Database connected successfully")
	return db, nil
}


func PingDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}


func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	log.Println("🔌 Closing database connection...")
	return sqlDB.Close()
}
