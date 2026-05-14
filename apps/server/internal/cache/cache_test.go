package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestCacheOperations(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Set instance manually for testing
	oldInstance := instance
	instance = &RedisCache{
		Client:     db,
		defaultTTL: time.Hour,
	}
	defer func() { instance = oldInstance }()

	ctx := context.Background()

	t.Run("Set operation", func(t *testing.T) {
		mock.ExpectSet("key1", "value1", instance.defaultTTL).SetVal("OK")
		err := Set(ctx, "key1", "value1")
		assert.NoError(t, err)
	})

	t.Run("Get operation", func(t *testing.T) {
		mock.ExpectGet("key1").SetVal("value1")
		val, err := Get(ctx, "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)
	})

	t.Run("Del operation", func(t *testing.T) {
		mock.ExpectDel("key1").SetVal(1)
		err := Del(ctx, "key1")
		assert.NoError(t, err)
	})

	t.Run("Get operation - Not Found", func(t *testing.T) {
		mock.ExpectGet("non-existent").RedisNil()
		_, err := Get(ctx, "non-existent")
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCacheUninitialized(t *testing.T) {
	oldInstance := instance
	instance = nil
	defer func() { instance = oldInstance }()

	ctx := context.Background()

	t.Run("Get should return Nil when uninitialized", func(t *testing.T) {
		_, err := Get(ctx, "any")
		assert.Error(t, err)
	})

	t.Run("Set should return nil when uninitialized", func(t *testing.T) {
		err := Set(ctx, "any", "value")
		assert.NoError(t, err)
	})

	t.Run("Del should return nil when uninitialized", func(t *testing.T) {
		err := Del(ctx, "any")
		assert.NoError(t, err)
	})
}
