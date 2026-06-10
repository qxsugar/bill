# Bill API 常用命令

.PHONY: dev run build-docs build fmt vet tidy test clean

# 默认：生成文档并运行
dev: build-docs run

# 运行 HTTP API 服务
run:
	@go run main.go

# 生成 swagger 文档（解析依赖中的引用类型）
build-docs:
	@swag fmt
	@swag init --parseDependency --parseInternal

# 格式化代码
fmt:
	@go fmt ./...

# 静态检查
vet:
	@go vet ./...

# 整理依赖
tidy:
	@go mod tidy

# 运行测试
test:
	@go test -v ./...

# 构建二进制
build:
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o app

# 清理构建产物
clean:
	@rm -f app
