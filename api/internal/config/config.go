package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Setup() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BILL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 默认配置项 BILL_SERVER_HOST、BILL_SERVER_PORT 等
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	// 日志配置 BILL_LOG_LEVEL、BILL_LOG_ENCODING、BILL_LOG_DEVELOPMENT、BILL_LOG_SAMPLING 等
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.encoding", "json")
	viper.SetDefault("log.development", false)
	viper.SetDefault("log.sampling", true)

	// CORS 配置 BILL_CORS_ALLOWED_ORIGINS 等
	viper.SetDefault("cors.allowed_origins", "*")

	// 微信小程序登录（code2session）。留空时走 dev 兜底：openid = dev_<code>
	// BILL_WECHAT_APPID、BILL_WECHAT_SECRET
	viper.SetDefault("wechat.appid", "")
	viper.SetDefault("wechat.secret", "")

	// JWT 签发密钥与有效期（小时）
	// BILL_JWT_SECRET、BILL_JWT_EXPIRE_HOURS
	viper.SetDefault("jwt.secret", "2EB47418-CCA9-4FAF-892D-F4117ABB06DA")
	viper.SetDefault("jwt.expire_hours", 24*30)

	// 数据库连接串：BILL_DEFAULT_DATABASE
	viper.SetDefault("default.database", "root:3aVhodFz@tcp(ppapi.cn:3306)/bill?charset=utf8mb4&parseTime=True&loc=Local")

	// Redis 连接配置，用于缓存微信 access_token 等
	// BILL_REDIS_ADDR、BILL_REDIS_PASSWORD、BILL_REDIS_DB
	viper.SetDefault("redis.addr", "ppapi.cn:6379")
	viper.SetDefault("redis.password", "jRz6pzFv")
	viper.SetDefault("redis.db", 0)
}
