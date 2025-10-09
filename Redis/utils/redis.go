package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var REDIS *redis.Client
var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	// ✅ assign the created client to the global variable
	REDIS = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})

	_, err := REDIS.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect Redis: %v", err)
	}

	log.Println("✅ Connected to Redis successfully")
	return REDIS
}
