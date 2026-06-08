package database

import (
	"context"
	"ecommerce/pkg/logger"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func ConnectDB(dsn string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		// log.Fatalf("❌ Gagal parsing config: %v\n", err)
		logger.Log.Error(" Gagal parsing config", zap.Error(err))
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		// log.Fatalf("❌ Gagal membuat koneksi ke database: %v\n", err)
		logger.Log.Error(" Gagal membuat koneksi ke database:", zap.Error(err))
	}

	err = pool.Ping(context.Background())
	if err != nil {
		// log.Fatalf("❌ Gagal ping database: %v\n", err)
		logger.Log.Error(" Gagal ping database:", zap.Error(err))
	}

	// fmt.Println("✅ Berhasil terhubung ke database PostgreSQL!")
	logger.Log.Info("Berhasil terhubung ke database PostgreSQL!")

	return pool
}

func RunMigrations(dsn string) {

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		// log.Fatalf("❌ Gagal membuat instance migrasi: %v\n", err)
		logger.Log.Error(" Gagal membuat instance migrasi:", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// log.Fatalf("❌ Gagal menjalankan migrasi: %v\n", err)
		logger.Log.Error(" Gagal menjalankan migrasi:", zap.Error(err))
	}

	defer m.Close()

	// fmt.Println("🚀 Migrasi Berhasil: Database sudah dalam versi terbaru!")
	logger.Log.Info("Migrasi Berhasil: Database sudah dalam versi terbaru!")

}
