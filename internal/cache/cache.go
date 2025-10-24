package cache

import (
	"context"
	"log"
	"time"

	"github.com/giridhar-m-a/custom_form_app/internal/utils"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

var instance *RedisCache

// Init initializes the package-level RedisCache singleton from environment variables and verifies connectivity.
// It returns immediately if the cache is already initialized. Reads configuration from these environment
// variables: REDIS_ADDR, REDIS_MAX_TIME_TO_KEEP_DATA, REDIS_USER, REDIS_PASSWORD, REDIS_DB, REDIS_POOL_SIZE,
// REDIS_MIN_IDLE_CONNS, and REDIS_DIAL_TIMEOUT. If REDIS_ADDR is empty or the initial Redis ping fails,
// Init logs a fatal error; on success it sets the client's default TTL and logs that Redis was initialized successfully.
func Init() {
	if instance != nil {
		return
	}

	addr := utils.GetEnv("REDIS_ADDR", "custom_form_app_redis:6379")
	if addr == "" {
		log.Fatal("REDIS_ADDR environment variable is not set")
	}

	maxTTL := utils.GetEnvAsInt("REDIS_MAX_TIME_TO_KEEP_DATA", 3600)
	user := utils.GetEnv("REDIS_USER", "")
	redisPassword := utils.GetEnv("REDIS_PASSWORD", "")

	instance = &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:         addr,
			Username:     user,
			Password:     redisPassword,
			DB:           utils.GetEnvAsInt("REDIS_DB", 0),
			PoolSize:     utils.GetEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: utils.GetEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			DialTimeout:  time.Duration(utils.GetEnvAsInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
		}),
		defaultTTL: time.Duration(maxTTL) * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := instance.client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis initialized successfully")
}

// Close closes the Redis client associated with the package singleton.
// If the cache has not been initialized, Close does nothing.
// Any error encountered while closing the client is logged.
func Close() {
	if instance != nil {
		if err := instance.client.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}
}

// ----------------------------
// Package-level helper methods
// Get retrieves the value associated with the given key from the Redis cache.
// It returns the value as a string. If the key does not exist, the returned error
// will be redis.Nil; other errors indicate Redis or connection failures.

func Get(ctx context.Context, key string) (string, error) {
	return instance.client.Get(ctx, key).Result()
}

// Set stores the given value under the specified key in the package Redis cache using the configured default TTL.
// It returns an error if the Redis SET operation fails.
func Set(ctx context.Context, key string, value interface{}) error {
	return instance.client.Set(ctx, key, value, instance.defaultTTL).Err()
}

// Del deletes the given key from the Redis-backed cache.
// It returns an error if the deletion fails.
func Del(ctx context.Context, key string) error {
	return instance.client.Del(ctx, key).Err()
}