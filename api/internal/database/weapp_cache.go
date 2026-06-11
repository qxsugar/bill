package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	weappcache "github.com/medivhzhan/weapp/v3/cache"
)

// WeappCache 用 Redis 实现 weapp 的 cache.Cache 接口，
// 让多实例共享微信 access_token 缓存，避免各自重复拉取。
type WeappCache struct {
	cli    *goredis.Client
	prefix string
}

var _ weappcache.Cache = (*WeappCache)(nil)

// NewWeappCache 创建带前缀的 weapp 缓存适配器。
func NewWeappCache(cli *goredis.Client) *WeappCache {
	return &WeappCache{cli: cli, prefix: "bill:weapp:"}
}

// Set 写入缓存，timeout 为 0 时表示永不过期。
func (c *WeappCache) Set(key string, val interface{}, timeout time.Duration) {
	s, ok := val.(string)
	if !ok {
		return
	}
	c.cli.Set(context.Background(), c.prefix+key, s, timeout)
}

// Get 读取缓存，返回值及是否命中。
func (c *WeappCache) Get(key string) (interface{}, bool) {
	s, err := c.cli.Get(context.Background(), c.prefix+key).Result()
	if err != nil {
		return nil, false
	}
	return s, true
}
