# Go E-Commerce Backend (Industrial Architecture)

Proyek ini adalah implementasi Backend menggunakan bahasa Go dengan standar **Clean Architecture** dan **Production-Ready Engineering**.

## 🚀 Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin Gonic
- **Database:** PostgreSQL
- **Driver:** pgx/v5 (with Connection Pooling)
- **Migration:** golang-migrate
- **Currency Handling:** shopspring/decimal

## 🛠️ Arsitektur
Mengikuti prinsip Clean Architecture:
- `internal/domain`: Kontrak interface dan entitas pusat.
- `internal/repository`: Implementasi SQL murni & akses database.
- `internal/usecase`: Logika bisnis utama.
- `internal/handler`: Layer API & input validation (Gin).

## 📈 Progres (Bulan 2 - Minggu 6)
- [x] Setup Environment & Database Connection Pooling
- [x] Database Migration System
- [x] Product Module (CRUD with Soft Delete)
- [ ] Order Module (ACID Transactions) - *Current Task*

## ⚙️ Cara Menjalankan
1. Clone repository.
2. Setup file `.env` (DB_URI, PORT).
3. Jalankan migrasi: `go run main.go` (otomatis memicu RunMigrations).
4. Server berjalan di `localhost:8080`.