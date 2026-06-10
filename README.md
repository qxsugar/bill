# 打牌记账

打牌记账是一款微信小程序，帮助多人打牌时记录积分、收支和账目。

## 功能模块

| 模块 | 说明 |
|------|------|
| 登录 | 微信 code2session + JWT |
| 房间 | 创建/加入/离开/结算，房主码 4 位优先含连号 |
| 支出 | 单人/均分/统一三种模式，可撤销，可感谢 |
| 实时 | WebSocket 广播房间事件（成员变动、支出、结算） |
| 记牌器 | 按牌型统计剩余张数，支持多副牌，单击扣减/双击归还 |
| 日志 | 房间操作流水，分页加载 |
| 个人 | 昵称 + emoji 头像选择 |

## 技术栈

**后端** (`/api`)
- Go 1.24 · Gin · GORM + MySQL · zap · Wire DI
- JWT HS256 · gorilla/websocket
- `github.com/qxsugar/pkg/kit` 统一响应格式

**前端** (`/mp`)
- 微信小程序原生（WXML / WXSS / JS）
- 无第三方 UI 库，CSS 变量统一设计 token

## 目录结构

```
bill/
├── api/
│   ├── cmd/http/          # 入口 + wire_gen.go
│   ├── internal/
│   │   ├── config/        # 配置读取
│   │   ├── dao/           # 数据访问层
│   │   ├── middleware/    # JWT 鉴权
│   │   ├── model/         # GORM 模型
│   │   ├── router/        # HTTP 路由处理
│   │   ├── service/       # 业务逻辑
│   │   └── ws/            # WebSocket Hub + Client
│   └── migrations/        # SQL 建表脚本
└── mp/
    ├── app.js             # 全局：登录、请求封装
    ├── api.js             # 所有接口调用
    ├── ws.js              # WebSocket 封装（自动重连）
    ├── components/        # user-avatar, expense-popup
    └── pages/
        ├── index/         # 首页：创建/加入房间
        ├── room/          # 房间主页
        ├── settle/        # 结算预览
        ├── settle-done/   # 结算完成
        ├── room-log/      # 房间日志
        ├── room-invite/   # 邀请好友（房间码+分享）
        ├── card-tracker/  # 记牌器
        ├── profile/       # 个人信息
        └── profile-avatar/# 选择 emoji 头像
```

## 本地运行

### 后端

```bash
# 配置环境变量
export BILL_WECHAT_APPID=your_appid
export BILL_WECHAT_SECRET=your_secret
export BILL_DB_DSN="user:pass@tcp(127.0.0.1:3306)/bill?charset=utf8mb4&parseTime=True"

# 初始化数据库
mysql -u root bill < api/migrations/001_init.sql

# 启动
cd api && go run ./cmd/http
```

### 前端

用微信开发者工具打开 `mp/` 目录，在 `app.js` 中将 `baseUrl` 改为后端地址。

## API 概览

```
POST /api/v1/auth.login          # 微信登录
GET  /api/v1/user/detail         # 个人信息
PUT  /api/v1/user/update         # 更新昵称/头像
GET  /api/v1/user/preset-avatars # 预设 emoji 列表

POST /api/v1/room/create         # 创建房间
POST /api/v1/room/join           # 加入房间
POST /api/v1/room/leave          # 离开房间
GET  /api/v1/room/detail         # 房间快照
POST /api/v1/room/settle         # 结算
GET  /api/v1/room/logs           # 日志分页

POST /api/v1/transaction/expense # 记一笔支出
POST /api/v1/transaction/revoke  # 撤销
POST /api/v1/transaction/thank   # 感谢

GET  /api/v1/card/detail         # 记牌器状态
POST /api/v1/card/adjust         # 调整某牌数量
POST /api/v1/card/reset          # 重置
POST /api/v1/card/set-deck       # 设置牌副数

GET  /ws/room?room_id=&token=    # WebSocket 连接
```

## WebSocket 消息格式

服务端推送 JSON：

```json
{ "event": "room_updated", "data": { ...RoomSnapshot } }
{ "event": "settled",      "data": { ...RoomSnapshot } }
```
