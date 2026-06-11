package weapp

import (
	"github.com/medivhzhan/weapp/v3"
	billredis "github.com/qxsugar/bill/api/internal/database"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// NewClient 按微信小程序配置创建 weapp 客户端，
// 并以 Redis 作为 access_token 缓存。未配置 appid/secret 时返回 nil，
// 由上层走 dev 兜底逻辑。
func NewClient(rdb *goredis.Client) *weapp.Client {
	appid := viper.GetString("wechat.appid")
	secret := viper.GetString("wechat.secret")
	if appid == "" || secret == "" {
		return nil
	}
	return weapp.NewClient(appid, secret, weapp.WithCache(billredis.NewWeappCache(rdb)))
}
