# Marketplace Structure Implementation Plan

> **For agentic workers:** Use executing-plans for the backend sequence and dispatching-parallel-agents for the independent frontend implementation. Track completion below.

**Goal:** 按现有工程分层整理商城代码，保持业务与接口兼容。

**Architecture:** controller → service → store；model 为三层共享类型。前端按业务页面拆分，通过共享组件和类型化 service 复用。

**Tech Stack:** Go、Gin、GORM、Wire、React、TypeScript、Ant Design、Vitest、Biome。

**Spec:** `docs/superpowers/specs/2026-09-17-marketplace-structure-design.md`

## Global Constraints

- 保持接口路径、请求字段、状态码、业务状态流转、数据库表及约束。
- 保留用户已有未提交改动，不增加依赖，不自动提交。
- 继续在当前工作目录落实用户指定的现有代码调整。

## Tasks

- [x] 基线：`GOCACHE=/private/tmp/mall-go-cache go test ./internal/marketplace ./internal/apiserver/... ./internal/pkg/response ./internal/userserver/...`。
- [x] 将 `internal/marketplace/model.go` 的表模型迁至 `model/`，增加共享请求、分页、筛选类型；保持原 GORM tag。
- [x] `store/store.go` 定义带 context 的类型化存储接口；`store/mysql/` 封装原 SQL 和事务、锁、分页。事务使用 `Transaction(ctx, func(Store) error)`。
- [x] 将 `service.go` 和原 handler 中的业务逻辑迁至 `service/` 的商品、商家、购物车、订单等文件，所有持久化经过 Store。
- [x] `controller/` 按业务拆分 HTTP 适配，统一响应工具；更新 API 构造、Wire、路由、生命周期、迁移及测试引用。
- [ ] 前端拆分商城业务页面和共享组件，显式路由，提取类型化 `service.ts`，保持现有操作与提示。
- [ ] 保留既有集成测试断言，添加响应和权限覆盖；运行 Go 测试/编译、前端测试/类型检查。
- [ ] 对新增和修改文件运行 gofmt/goimports/Biome，审查差异及静态检查结果，记录未执行的外部集成验证。
