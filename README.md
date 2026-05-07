# Go E-Commerce Backend (Industrial Architecture)

Proyek ini adalah implementasi Backend menggunakan bahasa Go dengan standar **Clean Architecture** dan **Production-Ready Engineering**.

## 🚀 Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin Gonic
- **Database:** PostgreSQL
- **Driver:** pgx/v5 (with Connection Pooling)
- **Migration:** golang-migrate
- **Currency Handling:** shopspring/decimal
- **Unique ID:** Google UUID

## 🛠️ Arsitektur
Mengikuti prinsip Clean Architecture:
- `internal/domain`: Kontrak interface dan entitas pusat.
- `internal/repository`: Implementasi SQL murni & akses database.
- `internal/usecase`: Logika bisnis utama.
- `internal/handler`: Layer API & input validation (Gin).

## ✨ Fitur Utama (Implemented)
- **ACID Transactions:** Menjamin integritas data saat checkout (Stok berkurang otomatis, rollback jika gagal).
- **Soft Delete:** Menghapus produk tanpa benar-benar menghilangkan data dari database.
- **Precision Calculation:** Menggunakan tipe data Decimal untuk akurasi finansial.
- **Dependency Injection:** Kode yang modular dan mudah ditest.

## 📈 Progres (Bulan 2 - Minggu 6)
- [x] Setup Environment & Database Connection Pooling
- [x] Database Migration System
- [x] Product Module (CRUD with Soft Delete)
- [x] Order Module (ACID Transactions - Checkout Logic)
- [ ] Authentication & Security (JWT, Bcrypt) - *Next Task*

## ⚙️ Cara Menjalankan
1. Clone repository.
2. Setup file `.env` (DB_URI, PORT).
3. Jalankan migrasi & server: `go run main.go`.
4. Server berjalan di `localhost:8080`.
5. Akses Swagger UI: http://localhost:8080/swagger/index.html