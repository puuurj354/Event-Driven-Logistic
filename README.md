# Event-Driven Logistic System

## Project Structure

```
├── cmd/                         # Entry point aplikasi
│   ├── order-service/           # main.go untuk Order Service
│   ├── payment-service/         # main.go untuk Payment Service
│   ├── inventory-service/       # main.go untuk Inventory Service
│   └── notification-service/    # main.go untuk Notification Service
│
├── internal/                    # Kode private yang tidak bisa di-import oleh aplikasi lain
│   ├── order/                   # Logika spesifik Order Service
│   │   ├── delivery/            # Handler (Gin HTTP handlers)
│   │   ├── repository/          # Database logic (GORM)
│   │   ├── service/             # Business logic
│   │   └── event/               # RabbitMQ Publisher/Subscriber logic
│   ├── payment/                 # Logika spesifik Payment Service
│   ├── inventory/               # Logika spesifik Inventory Service
│   └── notification/            # WebSocket & Notification logic
│
├── pkg/                         # Shared library (bisa digunakan antar service)
│   ├── broker/                  # Helper untuk koneksi RabbitMQ
│   ├── database/                # Helper untuk koneksi PostgreSQL
│   ├── config/                  # Viper configuration loader
│   └── response/                # Standarisasi format JSON response
│
├── api/                         # Dokumentasi API (Swagger/Postman collection)
├── deployments/                 # Docker & Docker Compose files
│   └── docker-compose.yaml
├── templates/                   # HTML Templates (HTMX)
├── .env                         # Variabel lingkungan (DB_URL, MQ_URL, dll)
├── go.mod                       # Go modules dependency
└── README.md                    # Dokumentasi proyek
```
