package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/raifuki/task-management/internal/config"
)

func NewRedis(cfg *config.Config) *redis.Client {
	var client *redis.Client

	if cfg.RedisURL != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %v", err)
		}
		client = redis.NewClient(opt)
		log.Println("📦 Using REDIS_URL from environment")
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
		log.Println("📦 Using individual REDIS_* env vars (local mode)")
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	log.Println("✅ Connected to Redis")
	return client
}
