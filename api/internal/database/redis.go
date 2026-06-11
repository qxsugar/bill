package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewClient 按配置建立 Redis 连接并校验连通性。
func NewClient() (*redis.Client, func(), error) {
	logger := zap.S()
	addr := viper.GetString("redis.addr")
	if addr == "" {
		return nil, nil, fmt.Errorf("redis address (BILL_REDIS_ADDR) is required")
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	logger.Info("redis connection established successfully")
	return cli, func() { _ = cli.Close() }, nil
}
