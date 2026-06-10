package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Setup() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BILL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.encoding", "json")
	viper.SetDefault("log.development", false)
	viper.SetDefault("log.sampling", true)

	viper.SetDefault("cors.allowed_origins", "*")

	// 微信小程序登录（code2session）。留空时走 dev 兜底：openid = dev_<code>
	viper.SetDefault("wechat.appid", "")
	viper.SetDefault("wechat.secret", "")

	// JWT 签发密钥与有效期（小时）
	viper.SetDefault("jwt.secret", "bill-dev-secret-change-me")
	viper.SetDefault("jwt.expire_hours", 24*30)

	// 数据库连接串：BILL_DEFAULT_DATABASE
	// 例：bill:2ZhwKarp@tcp(ppapi.cn:3306)/bill?charset=utf8mb4&parseTime=True&loc=Local
	viper.SetDefault("default.database", "root:3aVhodFz@tcp(ppapi.cn:3306)/bill?charset=utf8mb4&parseTime=True&loc=Local")
}
