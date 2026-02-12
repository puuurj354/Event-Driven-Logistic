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


---

## 🚀 Project Overview: Event-Driven Logistics Platform

Proyek ini adalah simulasi sistem logistik sederhana (pengiriman barang/makanan) skala mikro yang dibangun dengan arsitektur **Microservices** berbasis **Event-Driven**. Fokus utama proyek ini adalah menangani sinkronisasi data antar layanan secara asinkron (tidak saling menunggu) dan memberikan pembaruan status secara *real-time* kepada pengguna.

### 🛠️ Tech Stack

* **Backend:** Go (Golang) dengan Framework **Gin**.
* **Frontend:** HTMX Template.
* **Message Broker:** **RabbitMQ** (sebagai jantung komunikasi antar layanan).
* **Database:** **PostgreSQL** (menggunakan **GORM** untuk manajemen data).
* **Real-time:** **WebSockets** (untuk notifikasi instan ke frontend).
* **Infrastructure:** **Docker & Docker Compose** (untuk orkestrasi layanan).
* **Maps Integration:** Integrasi koordinat untuk pelacakan kurir.



---

### 🏛️ Desain Sistem & Arsitektur

Sistem ini memecah fungsi besar menjadi beberapa layanan mandiri yang berkomunikasi melalui **Pub/Sub (Publish/Subscribe)**:

1. **API Gateway:** Menjadi gerbang utama untuk permintaan dari user.
2. **Order Service:** Mengelola pembuatan pesanan awal.
3. **Inventory Service:** Mengelola ketersediaan stok secara otomatis.
4. **Payment Service:** Mensimulasikan validasi pembayaran.
5. **Delivery Service:** Menangani logistik, kurir, dan pelacakan koordinat.
6. **Notification Service:** Jembatan WebSocket yang mendorong status terbaru ke layar user.

---

### 🔄 Alur Kerja Utama (Event Flow)

Berdasarkan diagram yang kita bahas, berikut adalah urutan kejadian saat sebuah pesanan dibuat:

* **Langkah 1:** User membuat pesanan  **Order Service** menyimpan data (`PENDING`) dan melempar event `order.created`.
* **Langkah 2:** **Inventory** dan **Payment Service** mendengar event tersebut secara bersamaan.
* **Langkah 3:** Jika stok aman dan pembayaran sukses, event `payment.success` dilempar.
* **Langkah 4:** **Delivery Service** menangkap sinyal sukses, mengalokasikan kurir, dan mulai mengirim koordinat lokasi.
* **Langkah 5:** **Notification Service** memantau semua aktivitas di broker dan mengirimkan update detik-demi-detik ke aplikasi user melalui **WebSockets**.

---

### 💾 Desain Database (PostgreSQL)

Setiap layanan memiliki tanggung jawab datanya sendiri (**Logical Separation**):

* **Orders:** Menyimpan status pesanan (`id`, `user_id`, `status`).
* **Products:** Mengelola kuantitas stok.
* **Payments:** Mencatat histori transaksi.
* **Shipments:** Menyimpan data kurir dan koordinat (`lat`, `long`).

---

### 📁 Struktur Kode (Clean Architecture)

Proyek ini menggunakan standar **Standard Go Project Layout** untuk memastikan kode mudah diuji dan dikembangkan:

* `cmd/`: Titik masuk (entry point) untuk setiap service.
* `internal/`: Berisi logika bisnis (Service), akses database (Repository), dan Handler API.
* `pkg/`: Library bersama untuk koneksi database dan message broker.

---

## Dependency

* **Web Framework**: github.com/gin-gonic/gin
* **Database**: gorm.io/gorm dan gorm.io/driver/postgres
* **Message Broker (RabbitMQ)**
* **WebSockets**: github.com/gorilla/websocket (untuk notifikasi real-time)
* **Config Management**: github.com/spf13/viper (untuk membaca file .env)
* **UUID**: github.com/google/uuid (karena kita tidak ingin pakai ID integer berurutan)
* **Dan lain-lain**