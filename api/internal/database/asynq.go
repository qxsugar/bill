package database

import (
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

// NewAsynqRedisOpt 复用现有 redis 配置（BILL_REDIS_ADDR 等）构建 asynq 连接选项。
func NewAsynqRedisOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
}

// NewAsynqClient 创建 asynq 任务投递客户端。
func NewAsynqClient(opt asynq.RedisClientOpt) (*asynq.Client, func()) {
	client := asynq.NewClient(opt)
	return client, func() { _ = client.Close() }
}
