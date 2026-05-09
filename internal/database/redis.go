package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnRedis(rdsURL string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     rdsURL,
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke Redis: %v", err)
	}

	fmt.Println("⚡ Berhasil terhubung ke Redis (In-Memory Cache)!")
	return client

}
