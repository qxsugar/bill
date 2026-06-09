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

	// 数据库连接串：BILL_DEFAULT_DATABASE
	// 例：bill:2ZhwKarp@tcp(ppapi.cn:3306)/bill?charset=utf8mb4&parseTime=True&loc=Local
	viper.SetDefault("default.database", "bill:2ZhwKarp@tcp(ppapi.cn:3306)/bill?charset=utf8mb4&parseTime=True&loc=Local")
}
