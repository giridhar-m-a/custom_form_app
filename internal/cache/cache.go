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

// Init initializes the Redis cache
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

// Close closes the redis client
func Close() {
	if instance != nil {
		if err := instance.client.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}
}

// ----------------------------
// Package-level helper methods
// ----------------------------

func Get(ctx context.Context, key string) (string, error) {	
	return instance.client.Get(ctx, key).Result()
}

func Set(ctx context.Context, key string, value interface{}) error {
	return instance.client.Set(ctx, key, value, instance.defaultTTL).Err()
}

func Del(ctx context.Context, key string) error {
	return instance.client.Del(ctx, key).Err()
}
