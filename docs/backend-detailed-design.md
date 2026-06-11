# LIVE LIFE 后端详细设计

状态：Clean Architecture 分层后基线
最后更新：2026-06-11

## 1. 后端目标

后端先保持轻量、稳定、容易迁移。

当前职责：

- 提供 LIVE LIFE 前端所需 API。
- 固定前后端数据契约。
- 返回中文 / 日本語 / English 三语言展示字段。
- 保存 Connect 表单消息。
- 通过 GORM 操作 SQLite。
- 为未来后台管理、PostgreSQL/MySQL、登录、订单、邮件通知留下清楚边界。

重要原则：

- 视觉方案可以从 V2、V3 继续演化。
- 后端 API 不跟随视觉版本变化。
- 多人协作时，业务逻辑和 API 契约必须先看本文档再修改。

## 2. 当前技术栈

```text
语言：Go
HTTP：net/http
ORM：GORM
数据库：SQLite
SQLite driver：github.com/glebarez/sqlite
入口：backend/cmd/server/main.go
静态兜底文件：backend/static/
本地后端端口：http://localhost:8080
```

## 3. Clean Architecture 分层

当前后端目录：

```text
backend/
  cmd/server/
    main.go
    main_test.go
  internal/
    domain/
      types.go
    application/
      service.go
    infrastructure/sqlite/
      models.go
      seed.go
      store.go
    interfaces/httpapi/
      server.go
```

职责：

```text
domain
  业务实体、字段定义、Connect 校验规则、品牌常量。

application
  用例服务。只依赖 domain 和 Repository 接口，不知道 SQLite/GORM/HTTP。

infrastructure/sqlite
  GORM model、SQLite 连接、迁移、种子数据、Repository 实现。

interfaces/httpapi
  HTTP handler、JSON 返回、CORS、静态文件兜底。

cmd/server
  程序启动和依赖组装，不放业务逻辑。
```

依赖方向：

```text
cmd/server
  -> interfaces/httpapi
  -> application
  -> domain

infrastructure/sqlite
  -> domain

application
  -> domain
```

application 通过接口依赖数据层：

```go
type Repository interface {
    ListEvents(ctx context.Context) ([]domain.Event, error)
    ListCatalogItems(ctx context.Context) ([]domain.CatalogItem, error)
    ListContents(ctx context.Context) ([]domain.ContentItem, error)
    SaveConnectMessage(ctx context.Context, message domain.ConnectMessage) error
}
```

这样未来如果从 SQLite 换 PostgreSQL，只需要新增一个 infrastructure/postgres 实现，不需要改 HTTP handler 和业务用例。

## 4. API 总览

```text
GET  /api/health
GET  /api/events
GET  /api/cd-items
GET  /api/contents
POST /api/connect
POST /api/join
```

说明：

- `/api/join` 只是兼容旧命名，内部复用 `/api/connect`。
- 没有 `/api/shop-items`。
- 没有顶层 Shop API。
- CD / 黑胶商业路径统一收敛到 `/api/cd-items`。

## 5. 通用字段约定

所有对外展示数据都必须带：

```json
{
  "brand": "LIVE LIFE"
}
```

多语言字段使用两类形式。

### 5.1 I18n 后缀字段

适合标题、摘要、备注：

```json
{
  "title": "Fallback title",
  "titleI18n": {
    "zh": "中文标题",
    "ja": "日本語タイトル",
    "en": "ENGLISH TITLE"
  }
}
```

前端读取优先级：

```text
当前语言 -> 中文 -> fallback 字段 -> 空字符串
```

### 5.2 直接三语言对象字段

适合按钮文案：

```json
{
  "purchaseText": {
    "zh": "点击此处购买",
    "ja": "こちらから購入",
    "en": "BUY HERE"
  }
}
```

## 6. /api/events

用途：

- 返回演出情报。
- 区分 LIVE LIFE 自主演出和推荐 / 档案演出。

核心字段：

```json
{
  "id": "redhair-japan-2026-july",
  "brand": "LIVE LIFE",
  "owned": true,
  "title": "Fallback title",
  "titleI18n": {},
  "date": "2026-07-10 / 2026-07-14",
  "time": "OPEN / START",
  "venue": "Tokyo",
  "price": "TBD",
  "lineup": [],
  "tags": [],
  "summary": "Fallback summary",
  "summaryI18n": {},
  "ticketNote": "Fallback ticket note",
  "ticketNoteI18n": {},
  "sourceNote": "Fallback source note",
  "sourceNoteI18n": {},
  "imageUrl": "/assets/events/example.jpg"
}
```

前端展示规则：

- `owned = true` 固定在 Shows 页面最上方。
- `owned = false` 放在推荐演出或 Archive 风格区域。
- `ticketNote` 用于说明外部票站或票务待确认。

数据库表：

```text
event_models
event_translation_models
event_lineup_models
event_tag_models
```

## 7. /api/cd-items

用途：

- 返回 CD 严选列表。
- 同时支持 CD 和黑胶分类。

商业路径：

```text
CD 严选 -> 单品卡片 / 详情 -> purchaseUrl -> 外部 Shop，例如 BASE
```

后端当前不处理支付，不创建站内订单。

返回结构：

```json
{
  "brand": "LIVE LIFE",
  "items": [],
  "cd": [],
  "vinyl": []
}
```

## 8. /api/contents

用途：

- 返回 Archive 或说明性内容。
- 不和演出、CD 商品混在一起。

未来可扩展字段：

- `imageUrl`
- `sourceUrl`
- `publishedAt`
- `relatedEventId`
- `tags`

## 9. /api/connect

用途：

- 接收联系表单。
- 保存到 SQLite。
- 后续可接邮件通知、客服系统、后台管理列表。

请求结构：

```json
{
  "nickname": "LIVE LIFE 朋友",
  "email": "you@example.com",
  "topic": "ticket",
  "message": "想咨询票务。"
}
```

校验规则：

- `nickname` 必填。
- `email` 必填并且包含 `@`。
- `topic` 必填。
- `message` 当前可为空。

返回结构：

```json
{
  "accepted": true,
  "brand": "LIVE LIFE",
  "message": "LIVE LIFE received your message.",
  "messageId": "conn_..."
}
```

## 10. 错误处理

API 返回 JSON 错误：

```json
{
  "error": "nickname is required"
}
```

约定：

- 前端可读。
- 不暴露内部堆栈。
- HTTP status 与错误类型匹配。

## 11. 数据库选型

当前：

```text
SQLite + GORM
```

原因：

- 当前是内容型站点和外部购买入口。
- 1GB 阿里云服务器资源有限。
- SQLite 文件备份、迁移、部署简单。
- 当前没有站内高并发订单交易。

未来升级条件：

- 需要多人后台同时编辑大量内容。
- 需要站内登录、订单、支付、库存。
- Connect 消息量明显增长。
- 需要更复杂的数据分析或权限管理。

未来可升级：

```text
PostgreSQL 或 MySQL
```

由于应用层依赖 Repository 接口，升级数据库时优先新增基础设施层实现，不改业务用例。
