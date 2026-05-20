# Go E-Commerce Backend (Industrial Architecture)

Proyek ini adalah implementasi Backend menggunakan bahasa Go dengan standar **Clean Architecture** dan **Production-Ready Engineering**.

## 🚀 Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin Gonic
- **Database:** PostgreSQL (pgx/v5 with Connection Pooling)
- **Caching:** Redis (go-redis/v9)
- **Configuration:** Viper (Environment Management)
- **Security:** JWT (JSON Web Tokens) & golang.org/x/crypto/bcrypt
- **Observability & Resilience:** golang.org/x/time/rate (Rate Limiting), Google UUID
- **Migration:** golang-migrate
- **Currency Handling:** shopspring/decimal

## 🛠️ Arsitektur

Mengikuti prinsip Clean Architecture:

- `internal/domain`: Kontrak interface dan entitas pusat.
- `internal/repository`: Implementasi SQL murni & akses database.
- `internal/usecase`: Logika bisnis utama & Caching Strategy.
- `internal/handler`: Layer API, input validation (Gin), & Response Standard.
- `internal/middleware`: Lapisan pertahanan (Auth, RBAC, Rate Limiting, Logger).

## ✨ Fitur Utama (Implemented)

- **ACID Transactions:** Menjamin integritas data saat checkout.
- **High Performance Caching:** Implementasi Cache-Aside Pattern dengan Redis untuk katalog produk.
- **Environment Management:** Konfigurasi fleksibel menggunakan Viper (.env & OS variables).
- **Clean Dependency Injection:** Kode modular dengan pemisahan interface yang ketat.
- **Security Layer:** Autentikasi berbasis JWT statless dan otorisasi dengan Role-Based Access Control (RBAC).
- **Infrastructure Protection:** Proteksi server dengan per-IP Rate Limiting, penyebaran Request ID (UUID) untuk *tracing*, dan *Structured Logging*.
- **Graceful Shutdown:** Memastikan server mati dengan aman tanpa memutus transaksi.

## 📈 Progres (Bulan 2 - Minggu 8)

- [x] Setup Environment & Database Connection Pooling
- [x] Database Migration System
- [x] Product Module (CRUD with Soft Delete)
- [x] Order Module (ACID Transactions - Checkout Logic)
- [x] Modular Refactoring & Viper Configuration
- [x] Caching Strategy (Redis Implementation)
- [x] Authentication & Security (JWT, Bcrypt, & Simple RBAC)
- [x] Infrastructure Middleware (Request ID, Logging, Rate Limiting)
- [ ] Self-Review Final Project API E-Commerce (Day 56) - _Next Task_
- [ ] Unit Testing & Table-Driven Tests (Bulan 3) - _Upcoming_

## ⚙️ Cara Menjalankan

1. Clone repository.
2. Setup file `.env` (DB_URI, REDIS_URL, PORT, JWT_SECRET).
3. Jalankan Redis & Postgres via Docker.
4. Jalankan migrasi & server: `go run cmd/api/main.go`.
5. Server berjalan di `localhost:8080`.
6. Akses Swagger UI: http://localhost:8080/swagger/index.html