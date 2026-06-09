# bill api

记账小程序后端服务，按 [jump-admin](https://github.com/jump-gate/jump-admin) 的分层规范组织。

## 目录结构

```
api/
├── main.go                  # 入口，调用 cmd.Execute()
├── cmd/
│   ├── root.go              # cobra 根命令，加载配置 + 启动 HTTP 应用
│   └── http/
│       ├── http.go          # Application：gin 引擎、中间件注册、路由注册、优雅退出
│       ├── wire.go          # wire 依赖注入声明（//go:build wireinject）
│       └── wire_gen.go      # 依赖装配实现
├── internal/
│   ├── config/              # viper 配置（环境变量前缀 BILL_）
│   ├── logger/              # zap 日志
│   ├── database/            # GORM + MySQL 连接
│   ├── middleware/          # cors / accesslog / recovery
│   ├── model/               # 数据模型（user / room / room_member / transaction / room_log）
│   ├── dao/                 # 数据访问层
│   ├── service/             # 业务逻辑层
│   ├── router/              # HTTP handler（返回 (any, error)，由 kit.TranslateFunc 统一包装）
│   └── provider.go          # wire ProviderSet
├── migrations/
│   └── 001_init.sql         # 建表脚本
└── sql/
    └── table.sql            # 表结构快照
```

分层调用链：`router → service → dao → model`，依赖通过 `wire` 注入。
响应格式与错误码复用 `github.com/qxsugar/pkg/kit`。

## 配置

通过环境变量覆盖（前缀 `BILL_`，`.` 替换为 `_`）：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `BILL_DEFAULT_DATABASE` | `bill:***@tcp(ppapi.cn:3306)/bill?...` | MySQL DSN |
| `BILL_SERVER_HOST` | `0.0.0.0` | 监听地址 |
| `BILL_SERVER_PORT` | `8080` | 监听端口 |
| `BILL_LOG_ENCODING` | `json` | `json` / `console` |

## 运行

```bash
go mod tidy
go run .            # 或 go build -o bill-api . && ./bill-api
```

健康检查：`GET /health` → `{"status":"ok"}`（含 DB ping）

## 数据库迁移

```bash
mysql -h ppapi.cn -P 3306 -u <user> -p bill < migrations/001_init.sql
```
