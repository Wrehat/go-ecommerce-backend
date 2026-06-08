package repository_test

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/repository"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestProductRepository_CreateProduct(t *testing.T) {

	// Jalankan container postgres
	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx, "postgres:15-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10*time.Second)),
	)

	require.NoError(t, err, "Gagal menyalakan container: ")

	defer pgContainer.Terminate(ctx)

	// Migration database
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	require.NoError(t, err)

	migrationPath, err := filepath.Abs("../../migrations")

	require.NoError(t, err)

	m, err := migrate.New("file://"+migrationPath, connStr)

	require.NoError(t, err, "Gagal membuat instance migrasi: ")

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Gagal menjalankan migrasi database: %v", err)
	}

	// Connection Pool & Init Repository
	dbPool, err := pgxpool.New(ctx, connStr)

	require.NoError(t, err, "Gagal membuat connection pool")

	defer dbPool.Close()

	// Table Driven Test
	repo := repository.NewProductRepository(dbPool)

	tests := []struct {
		name        string
		input       domain.Product
		expectError bool
		expectedErr error
	}{
		{
			name: "Skenario 1: Sukses membuat produk baru",
			input: domain.Product{
				SKU:   "SKU-ARSITEK-001",
				Name:  "Meja Gambar Arsitek",
				Price: decimal.NewFromInt(2500000),
				Stock: 15,
			},
			expectError: false,
			expectedErr: nil,
		},
		{
			name: "Skenario 2: Gagal karena duplikasi SKU",
			input: domain.Product{
				SKU:   "SKU-ARSITEK-001",
				Name:  "Meja Gambar Tiruan",
				Price: decimal.NewFromInt(1200000),
				Stock: 5,
			},
			expectError: true,
			expectedErr: domain.ErrSKUDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdProduct, err := repo.CreateProduct(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)

				if tt.expectedErr != nil {
					assert.True(t, errors.Is(err, tt.expectedErr), "Error harus berupa %v, tapi mendapat %v", tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)

				// t.Logf("ID Produk = %d, Dibuat Pada = %v", createdProduct.ID, createdProduct.CreatedAt)

				// assert.NotZero(t, createdProduct.ID)
				assert.NotZero(t, createdProduct.CreatedAt)
				// assert.Zero(t, createdProduct.CreatedAt)
				// assert.NotZero(t, createdProduct.UpdatedAt)
			}
		})
	}

}
