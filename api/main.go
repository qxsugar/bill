package main

import "github.com/qxsugar/bill/api/cmd"

// ref: https://github.com/swaggo/swag/blob/master/README_zh-CN.md
//
//	@title						Bill API
//	@version					1.0
//	@description				记账小程序后端 API 文档
//	@contact.name				Developer
//	@contact.url				https://github.com/qxsugar
//	@contact.email				qxsugar@gmail.com
//	@BasePath					/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				需登录接口：在请求头携带 "Bearer {token}"
func main() {
	cmd.Execute()
}
