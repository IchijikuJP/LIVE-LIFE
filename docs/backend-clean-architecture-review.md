# 后端 Clean Architecture 审查记录

状态：初版拆分完成
最后更新：2026-06-11

## 1. 本次为什么拆分

之前后端主要逻辑集中在 `backend/cmd/server/main.go`。

问题：

- HTTP handler、数据库 model、种子数据、业务校验、启动逻辑混在一个文件里。
- 多人协作时容易互相冲突。
- 后续换数据库或加后台管理时，很难判断该改哪里。
- 业务逻辑没有固定边界，容易被前端视觉改版带着一起乱动。

本次目标：

- 把启动、HTTP、业务用例、数据库实现拆开。
- 固定业务实体和 API 契约。
- 让未来新增数据或换数据库时，不影响前端和业务层。

## 2. 拆分后的目录

```text
backend/cmd/server
  程序入口。只负责读取配置、创建 Store、创建 Service、创建 HTTP Server。

backend/internal/domain
  业务实体和基础规则。

backend/internal/application
  用例服务。包含 Connect 提交等业务流程。

backend/internal/infrastructure/sqlite
  SQLite + GORM 的具体实现。

backend/internal/interfaces/httpapi
  HTTP API 适配层。
```

## 3. 已完成的边界

- `domain` 不依赖 HTTP、GORM、SQLite。
- `application` 不依赖 HTTP、GORM、SQLite。
- `application` 只通过 Repository 接口访问数据。
- `sqlite` 实现 Repository。
- `httpapi` 调用 application，不直接操作 GORM。
- `cmd/server` 不放业务规则。

## 4. 保持不变的外部行为

API 保持：

```text
GET  /api/health
GET  /api/events
GET  /api/cd-items
GET  /api/contents
POST /api/connect
POST /api/join
```

前端不需要因为后端拆层而改代码。

## 5. 当前测试覆盖

已覆盖：

- `/api/health` 返回 LIVE LIFE。
- `/api/events` 返回 ownedEvents。
- `/api/cd-items` 返回 CD / vinyl / purchaseUrl。
- `/api/shop-items` 仍然不存在。
- `/api/connect` 校验邮箱。
- `/api/connect` 成功后写入 SQLite。

## 6. 审查结论

本次拆分是结构性重构，不改变业务目标。

可以继续作为多人协作的后端基线。

后续新增业务时建议：

- 新业务实体先放 `domain`。
- 新用例先放 `application`。
- 新数据库实现放 `infrastructure`。
- 新 HTTP 路由放 `interfaces/httpapi`。
- 不要把业务逻辑写回 `cmd/server/main.go`。
