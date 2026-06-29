package database

import (
	"context"
	"ecommerce/pkg/logger"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnRedis(rdsURL string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rdsURL,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal ping ke redis: %w", err)
	}

	logger.Log.Info("Berhasil terhubung ke Redis (In-Memory Cache)!")

	return client, nil

}
