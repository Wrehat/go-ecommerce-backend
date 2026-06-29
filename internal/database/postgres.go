package database

import (
	"context"
	"ecommerce/pkg/logger"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("Gagal parsing config %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat koneksi: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Gagal ping database: %w", err)
	}

	logger.Log.Info("Berhasil terhubung ke database PostgreSQL!")

	return pool, nil
}

func RunMigrations(dsn string) error {

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("Gagal instance migration: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Gagal menjalankan migration: %w", err)
	}

	logger.Log.Info("Migrasi Berhasil: Database sudah dalam versi terbaru!")
	return nil
}
