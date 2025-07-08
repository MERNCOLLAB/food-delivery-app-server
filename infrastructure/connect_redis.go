package infrastructure

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func ConnectRedis() {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		log.Fatal("Environmental variable REDIS_URL does not exist")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse the REDIS_URL: %v", err)
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to the Redis Upstash: %v", err)
	}

	log.Printf("🍪 Connected successfully to Redis Client")
}
