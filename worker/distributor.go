package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistribtor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistribtor struct {
	client *asynq.Client
}

func NewRedisTaskDistribtor(redisOpt asynq.RedisClientOpt) TaskDistribtor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistribtor{
		client: client,
	}
}