// Package db manages database connections
package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

// InitRedis establishes a connection to Redis instance
func InitRedis() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Println("REDIS_URL not found, using default localhost:6379")
		url = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	RDB = redis.NewClient(opt)

	// Ping to verify connection
	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}

	log.Println("Redis connection successful.")
}

// GetOrSet is a generic helper to implement the Cache-Aside pattern.
func GetOrSet[T any](ctx context.Context, key string, expiration time.Duration, dbFetch func() (*T, error)) (*T, error) {
    // Attempt to get from Redis
    val, err := RDB.Get(ctx, key).Result()
    if err == nil {
        var result T
        if err := json.Unmarshal([]byte(val), &result); err == nil {
            return &result, nil
        }
    }

    // Cache Miss: Fetch from DB
    data, err := dbFetch()
    if err != nil {
        return nil, err
    }

    // Capture Marshalling Error
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Printf("CACHE_ERROR: Failed to marshal data for key %s: %v", key, err)
        return data, nil // Return DB data
    }

    // Capture Redis SET Error
    if err := RDB.Set(ctx, key, jsonData, expiration).Err(); err != nil {
        log.Printf("CACHE_ERROR: Failed to set key %s in Redis: %v", key, err)
    }

    return data, nil
}