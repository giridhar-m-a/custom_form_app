package scheduler

import (
	"github.com/giridhar-m-a/custom_form_app/internal/cache"
	"github.com/hibiken/asynq"
)

func NewRedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     cache.Client.Options().Addr,
		Password: cache.Client.Options().Password,
		DB:       cache.Client.Options().DB,
	}
}

func NewClient() *asynq.Client {
	opts := NewRedisClientOpt()
	return asynq.NewClient(opts)
}
