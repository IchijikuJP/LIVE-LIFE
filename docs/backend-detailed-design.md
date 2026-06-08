# LIVE LIFE 后端详细设计

状态：当前开发基准  
最后更新：2026-06-08

## 1. 后端目标

后端先保持轻量、稳定、易迁移。

当前阶段的职责：

- 提供本地预览 API。
- 固定前后端数据契约。
- 返回三语言展示字段。
- 承接 Connect 表单验证。
- 为未来数据库、管理后台、部署上线留下清晰边界。

视觉设计会从 V2、V3 继续演化，但后端 API 暂时不跟着视觉版本改变。

## 2. 当前技术栈

```text
语言：Go
入口：backend/cmd/server/main.go
静态文件：backend/static/
本地端口：http://localhost:8080
```

当前是内存种子数据，不接数据库。

后续可以把 seed 数据迁移到：

- SQLite。
- GORM ORM。
- SQL migration 文件。
- 后续管理后台或 CMS。

数据库与 ORM 的详细表设计见：

```text
docs/database-schema-draft.md
```

当前建议：

```text
P1 数据库：SQLite
Go ORM：GORM
SQLite Driver：优先 github.com/glebarez/sqlite
Migration：SQL migration 文件，建议 goose
```

这个选择的重点是：当前 LIVE LIFE 是内容型站点和外部购买入口，不是高并发站内交易系统，所以先用单文件 SQLite 降低部署、备份和维护成本。

## 3. API 总览

```text
GET  /api/health
GET  /api/events
GET  /api/cd-items
GET  /api/contents
POST /api/connect
POST /api/join
```

说明：

- `/api/join` 只是兼容旧表单命名，内部和 `/api/connect` 使用同一套逻辑。
- 没有 `/api/shop-items`。
- 没有顶层 Shop API。
- CD/黑胶商业路径统一收敛到 `/api/cd-items`。

## 4. 通用字段约定

所有对外展示数据都应该带：

```json
{
  "brand": "LIVE LIFE"
}
```

多语言字段采用两种形式。

### 4.1 I18n 后缀字段

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

### 4.2 直接三语言对象字段

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

前端读取优先级：

```text
当前语言 -> 中文 -> 默认按钮文案
```

## 5. /api/health

用途：

- 前端检查 API 是否可用。
- 本地预览显示连接状态。

返回示例：

```json
{
  "brand": "LIVE LIFE",
  "service": "LIVE LIFE API",
  "status": "ok",
  "time": "2026-06-08T00:00:00Z"
}
```

## 6. /api/events

用途：

- 返回演出情报。
- 区分 LIVE LIFE 自主演出和推荐/档案演出。

核心字段：

```json
{
  "id": "redhair-2026-july",
  "brand": "LIVE LIFE",
  "owned": true,
  "title": "Fallback title",
  "titleI18n": {},
  "date": "2026.07.10 / 2026.07.14",
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

- `owned = true` 的活动固定在 Shows 页面最上方。
- `owned = false` 的活动放在推荐或 Archive 风格区域。
- `ticketNote` 用于说明外部票站或票务待确认。

未来数据库建议：

```text
events
event_translations
event_lineup
event_tags
```

## 7. /api/cd-items

用途：

- 返回 CD 严选列表。
- 同时支持 CD 和黑胶分类。

核心字段：

```json
{
  "id": "redhair-demo-cd",
  "brand": "LIVE LIFE",
  "format": "cd",
  "artist": "紅髪少年殺人事件",
  "title": "Fallback title",
  "titleI18n": {},
  "summary": "Fallback summary",
  "summaryI18n": {},
  "tracks": [],
  "status": "external shop",
  "price": "TBD",
  "imageUrl": "/assets/events/redhair-2026-july.jpg",
  "purchaseUrl": "https://thebase.com/",
  "purchaseText": {
    "zh": "点击此处购买",
    "ja": "こちらから購入",
    "en": "BUY HERE"
  }
}
```

返回结构：

```json
{
  "brand": "LIVE LIFE",
  "items": [],
  "cd": [],
  "vinyl": []
}
```

说明：

- `items` 是完整列表。
- `cd` 是 CD 分类。
- `vinyl` 是黑胶分类。
- 前端可以直接用 `items` 自己筛选，也可以使用后端分组。

购买路径：

```text
单品卡片 -> purchaseUrl -> 外部 Shop
```

后端不处理支付、不创建订单。

未来数据库建议：

```text
catalog_items
catalog_item_translations
catalog_item_tracks
```

## 8. /api/contents

用途：

- 返回 Archive 或说明性内容。
- 不和演出、CD 商品混在一起。

核心字段：

```json
{
  "id": "about-live-life",
  "brand": "LIVE LIFE",
  "type": "note",
  "title": "Fallback title",
  "titleI18n": {},
  "summary": "Fallback summary",
  "summaryI18n": {}
}
```

未来可扩展：

- `imageUrl`
- `sourceUrl`
- `publishedAt`
- `relatedEventId`
- `tags`

## 9. /api/connect

用途：

- 接收联系表单。
- 当前只做本地验证。
- 后续可接邮件、数据库或客服系统。

请求结构：

```json
{
  "nickname": "LIVE LIFE 朋友",
  "email": "you@example.com",
  "topic": "ticket",
  "message": "想咨询票务。"
}
```

当前验证规则：

- `nickname` 必填。
- `email` 必填。
- `topic` 必填。
- `message` 当前可以为空，但建议前端引导填写。

返回示例：

```json
{
  "brand": "LIVE LIFE",
  "message": "LIVE LIFE received your message.",
  "topic": "ticket"
}
```

未来数据库建议：

```text
connect_messages
```

字段建议：

- `id`
- `nickname`
- `email`
- `topic`
- `message`
- `status`
- `created_at`
- `handled_at`

## 10. 错误处理

当前 API 返回 JSON 错误：

```json
{
  "error": "nickname is required"
}
```

建议保持：

- 前端可读。
- 不暴露内部堆栈。
- HTTP status 与错误类型匹配。

## 11. 后端不跟随设计版本变化

未来 Review 模式会出现：

```text
V2 当前版 / V3 Band Signal / V4 ...
```

但后端仍然保持同一套 API。

原因：

- 设计版本只是视觉 Renderer。
- 数据结构不应该因为颜色、动效、布局变化而改变。
- 客户 Review 可以对比视觉，但内容数据必须一致。

## 12. 下一步建议

待审批后可以补：

- 数据库表结构文档。
- API 示例 JSON 文件。
- Connect 表单邮件通知方案。
- 外部链接安全策略。
- 管理后台最小字段设计。
