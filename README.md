# Go E-Commerce Backend (Industrial Architecture)

Proyek ini adalah implementasi Backend menggunakan bahasa Go dengan standar **Clean Architecture** dan **Production-Ready Engineering**.

## 🚀 Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin Gonic
- **Database:** PostgreSQL (pgx/v5 with Connection Pooling)
- **Caching:** Redis (go-redis/v9)
- **Configuration:** Viper (Environment Management)
- **Security:** JWT (JSON Web Tokens) & golang.org/x/crypto/bcrypt
- **Observability & Resilience:** - Uber Zap (Structured Logging)
  - OpenTelemetry (Distributed Tracing)
  - golang.org/x/time/rate (Rate Limiting)
- **Testing:** testify (Assertions), testcontainers-go (Docker-based Integration Testing)
- **Migration:** golang-migrate
- **Currency Handling:** shopspring/decimal
- **API Documentation:** Swagger (swaggo/swag)

## 🛠️ Arsitektur

Mengikuti prinsip Clean Architecture:

- `internal/domain`: Kontrak interface dan entitas pusat.
- `internal/repository`: Implementasi SQL murni & akses database.
- `internal/usecase`: Logika bisnis utama & Caching Strategy.
- `internal/handler`: Layer API, input validation (Gin), & Response Standard.
- `internal/middleware`: Lapisan pertahanan (Auth, RBAC, Rate Limiting, Logger, OTel).

## ✨ Fitur Utama (Implemented)

- **ACID Transactions:** Menjamin integritas data saat checkout.
- **High Performance Caching:** Implementasi Cache-Aside Pattern dengan Redis untuk katalog produk.
- **Environment Management:** Konfigurasi fleksibel menggunakan Viper (.env & OS variables).
- **Clean Dependency Injection:** Kode modular dengan pemisahan interface yang ketat.
- **Security Layer:** Autentikasi berbasis JWT statless dan otorisasi dengan Role-Based Access Control (RBAC).
- **Infrastructure Protection:** Proteksi server dengan per-IP Rate Limiting, penyebaran Request ID (UUID) untuk _tracing_, dan _Structured Logging_.
- **Distributed Tracing:** Terintegrasi dengan OpenTelemetry (OTel) untuk pemantauan latensi API secara presisi.
- **Graceful Shutdown:** Memastikan server mati dengan aman tanpa memutus transaksi.
- **Automated Testing:** Pengujian logika murni (Table-Driven Tests) dan Integration Testing database tersolasi secara _on-the-fly_ menggunakan Docker & Testcontainers.

## 📈 Progres (Bulan 3 - Minggu 9)

- [x] Setup Environment & Database Connection Pooling
- [x] Database Migration System
- [x] Product Module (CRUD with Soft Delete)
- [x] Order Module (ACID Transactions - Checkout Logic)
- [x] Modular Refactoring & Viper Configuration
- [x] Caching Strategy (Redis Implementation)
- [x] Authentication & Security (JWT, Bcrypt, & Simple RBAC)
- [x] Infrastructure Middleware (Request ID, Logging, Rate Limiting)
- [x] Self-Review Final Project API E-Commerce
- [x] Unit Testing & Table-Driven Tests (Day 57-58)
- [x] Integration Testing dengan testcontainers-go (Day 59-60)
- [x] Structured Logging dengan Uber Zap (Day 61)
- [x] OpenTelemetry (OTel) Minimal Instrumentation (Day 61)
- [x] Profiling Performa (pprof) (Day 62)

## 🚀 Fase Selanjutnya (Deployment & DevOps)

- [x] Dockerisasi Aplikasi (Dockerfile & Multi-stage Build) *(Day 29-30)*
- [ ] Docker Compose (Orkestrasi Multi-Container)
- [ ] CI/CD Pipeline (GitHub Actions untuk Automated Testing)
- [ ] Cloud Deployment (Deploy API ke VPS / PaaS)

## ⚙️ Cara Menjalankan

### 1. Menyiapkan Proyek

1. Clone repository.
2. Buat file `.env` dari contoh di bawah (isi sesuai konfigurasi lokal kamu):

```env
DB_URI=postgresql://user:password@host.docker.internal:5433/toko_db?sslmode=disable
PORT=8080
REDIS_URL=host.docker.internal:6379
JWT_SECRET=your_secret_key
```

3. Pastikan **Docker Desktop** dalam keadaan menyala.
4. Pastikan container **PostgreSQL** dan **Redis** sudah berjalan.

### 2. Menjalankan Server API (Development Mode)

Proyek ini menggunakan `Makefile` untuk mempermudah eksekusi infrastruktur dan aplikasi secara bersamaan.

```bash
make run
```

### 3. Menjalankan via Docker Container

Aplikasi ini sudah di-Dockerisasi menggunakan **Multi-Stage Build** untuk menghasilkan image yang ringan berbasis Alpine Linux.

**Build image:**
```bash
docker build -t tokokana-api:v1 .
```

**Jalankan container:**
```bash
docker run -d --name tokokana-server -p 8080:8080 --env-file .env tokokana-api:v1
```

> **Catatan:** Saat menggunakan Docker, `REDIS_URL` dan `DB_URI` harus menggunakan `host.docker.internal` (bukan `localhost`) agar container API bisa menjangkau service yang berjalan di Mac host.

**Cek logs:**
```bash
docker logs tokokana-server
```

**Akses API:**
- Swagger Docs: `http://localhost:8080/swagger/index.html`
- Health check: `http://localhost:8080/ping`

### 4. Analisis Performa (Profiling dengan pprof)

Aplikasi ini sudah dilengkapi dengan sensor `pprof` untuk membedah performa CPU dan memori secara _real-time_.

1. Pastikan server sedang berjalan dan sedang menerima trafik (bisa gunakan Postman atau alat Load Test seperti `hey`).
2. Buka terminal baru dan jalankan salah satu perintah berikut:

**Melihat Profil CPU (Selama 10 detik):**
```bash
go tool pprof -http=:8081 "http://localhost:8080/debug/pprof/profile?seconds=10"
```

**Melihat Profil Memori (Heap/RAM):**
```bash
go tool pprof -http=:8081 "http://localhost:8080/debug/pprof/heap"
```

_(Catatan: Anda membutuhkan Graphviz yang terinstal di OS Anda untuk melihat visualisasi Flame Graph di browser)._
