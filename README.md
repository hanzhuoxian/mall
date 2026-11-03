# Mall 多商家商城

Go / Gin / gRPC / MySQL / Redis + React / Umi / Ant Design。

[完整功能规划与阶段边界](docs/商城功能规划.md) · [接口与验收说明](docs/接口与验收.md)

当前实现可持久化的多商家交易核心：注册、商家入驻和审核、分类、商品审核、库存、收货地址、购物车、跨店拆单、订单履约、基础售后、佣金统计与审计。真实支付、分账打款、营销、多规格等尚未完成，详见功能矩阵。

## 本地启动

需要 Go 1.26.3、Node.js 22 或更新版本、MySQL 和 Redis。使用仓库已有依赖版本。

```sh
docker compose -f deployments/docker-compose.dev.yaml up -d mysql redis elasticsearch
```

新 MySQL 数据卷会自动创建 `mall` 和 `mall_user` 数据库。已有数据卷不执行初始化 SQL，需要自行创建缺失的数据库；不要删除现有数据卷。

迁移会为用户名和邮箱增加唯一索引。已有数据库请先备份并检查重复值，不能直接忽略迁移失败：

```sql
SELECT username, COUNT(*) FROM mall_user.sys_user GROUP BY username HAVING COUNT(*) > 1;
SELECT email, COUNT(*) FROM mall_user.sys_user GROUP BY email HAVING COUNT(*) > 1;
```

```sh
go run ./cmd/mallctl migrate
# 替换为自己的密码，管理员用户名使用 admin。
go run ./cmd/mallctl init-admin --password 'Replace-With-Your-Strong-Password' --email admin@example.com
```

在两个终端分别启动用户服务和 API 服务。默认连接本机 MySQL（root / 123456）和 Redis；这些是现有开发环境凭据，正式环境必须通过启动配置替换。

```sh
go run ./cmd/mall-user-server
```

```sh
# 临时随机密钥只用于本次本地运行；重启后旧令牌失效。
# 正式环境从密钥管理系统提供固定的 32 字节以上密钥。
go run ./cmd/mall-api-server --server.jwt-secret "$(openssl rand -hex 32)" --sandbox-payment --elasticsearch-url http://localhost:9200
```

首次接入 Elasticsearch 或需要修复索引时，执行全量重建（命令会先清空 `mall-products` 索引，再从 MySQL 回填当前在售商品）：

```sh
go run ./cmd/mallctl reindex-search --elasticsearch-url http://localhost:9200
```

```sh
cd web/mall-admin
npm ci
npm run dev
```

前端开发代理指向 API `9091`。生产前端改用同源请求，反向代理须将 `/v1/`、`/login`、`/logout`、`/refresh`、`/register`、`/captcha` 转发至 API 服务。

## 验收路径

1. 管理员登录，新增商品分类。
2. 打开 `/user/register`，注册两个商家账号和一个消费者账号；商家登录后提交入驻申请。
3. 管理员审核两家店铺，分别设置佣金比例。
4. 商家发布商品，平台审核通过后消费者可在「逛商城」看到。
5. 消费者添加收货地址，将两家店的商品加入购物车并提交订单，系统生成两个店铺订单。
6. 开发环境开启沙箱时，消费者可逐单模拟付款；商家选择「我的店铺订单」填写物流并发货；消费者确认收货。
7. 消费者可发起售后，商家或平台处理；检查库存回补、结算统计和审计记录。

关闭 `--sandbox-payment` 后不能模拟付款/退款，系统也不会假称真实支付成功。发布模式显式使用 `--server.mode=release`，该模式拒绝开启沙箱支付。

## 日志开发约定

`pkg/logger` 负责 zap 配置、初始化和请求上下文关联，业务日志直接使用 zap API：

```go
// 服务启动时：先初始化 telemetry，再初始化日志；退出时同步日志。
if err := logger.Init(logOptions); err != nil {
    return err
}
defer func() { _ = zap.L().Sync() }()

zap.L().Info("server started", zap.String("address", address))
logger.L(ctx).Info("request handled", zap.String("operation", operation))
```

HTTP 处理中使用 `c.Request.Context()`；自定义请求字段通过 `logger.WithFields(ctx, zap.String(...))` 传递。`logger.L(ctx)` 返回原生 `*zap.Logger`，同时关联本地日志和 OTLP 的 trace/span ID。GORM 适配器通过 `NewGormLogger(zapLogger, level)` 显式接收 logger。

`logger.Init` 安装 `zap.L()`、`zap.S()` 和标准库 `log` 的全局输出；需要独立实例时使用无全局副作用的 `logger.New`。原有 `log.*` 配置保留，本地与 OTLP 共用日志级别和采样策略：每秒相同级别、相同消息保留前 100 条，此后每 100 条保留 1 条。

## 自动验证

```sh
go test ./cmd/... ./internal/... ./pkg/... ./proto/...
cd web/mall-admin
npm run tsc
npm test
npm run build
```

默认交易集成测试使用临时 SQLite 数据库，验证持久化、回滚、拆单、幂等、权限和状态流转。MySQL 的真实行锁并发测试需指向**独立测试数据库**，不要使用业务数据库：

```sh
MALL_TEST_MYSQL_DSN='user:password@tcp(127.0.0.1:3306)/mall_test?parseTime=true' go test ./internal/marketplace -v
```

该测试变量启用 MySQL 自动迁移并写入随机测试数据；测试数据库由测试环境负责清理。
