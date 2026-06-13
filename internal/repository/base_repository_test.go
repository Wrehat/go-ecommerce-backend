package repository_test

import (
	"context"
	"ecommerce/migrations"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Definisikan struct
type BaseRepositoryTestSuite struct {
	suite.Suite
	PgContainer *postgres.PostgresContainer
	DbPool      *pgxpool.Pool
}

// Buat SetupSuite()
func (ts *BaseRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	// 1. Setup container Postgres
	pgContainer, err := postgres.Run(
		ctx, "postgres:15-alpine",
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		postgres.WithDatabase("test_db"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(10*time.Second),
		),
	)

	ts.Require().NoError(err)
	ts.PgContainer = pgContainer

	// 2. Migrations
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	ts.Require().NoError(err)

	migrationURL := strings.Replace(connStr, "postgres://", "pgx5://", 1)

	srcDriver, err := iofs.New(migrations.FS, ".")
	ts.Require().NoError(err)

	m, err := migrate.NewWithSourceInstance("iofs", srcDriver, migrationURL)
	ts.Require().NoError(err)

	defer func() {
		srcErr, dbEr := m.Close()
		ts.Require().NoError(srcErr)
		ts.Require().NoError(dbEr)
	}()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		ts.Require().NoError(err)
	}

	// 3. Database Pool
	dbPool, err := pgxpool.New(ctx, connStr)
	ts.Require().NoError(err)
	ts.DbPool = dbPool
}

// Implementasi Tear Down
func (ts *BaseRepositoryTestSuite) TearDownSuite() {
	ts.DbPool.Close()
	ts.Require().NoError(ts.PgContainer.Terminate(context.Background()))
}

func (ts *BaseRepositoryTestSuite) SetupTest() {
	_, err := ts.DbPool.Exec(context.Background(), "TRUNCATE TABLE order_items, orders, products, users RESTART IDENTITY CASCADE")
	ts.Require().NoError(err)
}
