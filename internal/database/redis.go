package database

import (
	"context"
	"ecommerce/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnRedis(rdsURL string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     rdsURL,
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		// log.Fatalf("❌ Gagal terhubung ke Redis: %v", err)
		logger.Log.Error("Gagal terhubung ke Redis", zap.Error(err))
	}

	// fmt.Println("⚡ Berhasil terhubung ke Redis (In-Memory Cache)!")
	logger.Log.Info("Berhasil terhubung ke Redis (In-Memory Cache)!")

	return client

}
