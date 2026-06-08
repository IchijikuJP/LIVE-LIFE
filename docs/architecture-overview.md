# LIVE LIFE 整体架构图

状态：当前开发基准  
最后更新：2026-06-08

这个文档用于从最高层说明 LIVE LIFE 的整体结构。  
设计版本以后会继续变化，但产品入口、API 逻辑和未来部署方向先保持稳定。

## 1. 总体架构

```mermaid
flowchart TD
  Audience["用户<br/>观众 / 买 CD 的人 / 合作方"] --> Browser["浏览器<br/>手机 / 桌面"]

  Browser --> ReviewShell["前端 Review Shell<br/>右上角可选 V2 / V3 / 未来版本"]

  ReviewShell --> VariantState["designVariant<br/>URL 参数 / 下拉框 / localStorage"]

  VariantState --> V2["V2 当前版<br/>网格系统 / Schedule 面板"]
  VariantState --> V3["V3 Band Signal<br/>经典摇滚乐队名密集纹理<br/>LIVE LIFE 字母轮播<br/>抽象音轨背景"]
  VariantState --> FutureDesign["未来设计版本<br/>V4 / V5"]

  V2 --> SharedPages["共享页面入口"]
  V3 --> SharedPages
  FutureDesign --> SharedPages

  SharedPages --> Shows["Shows<br/>演出情报"]
  SharedPages --> CDSelect["CD 严选<br/>CD / 黑胶"]
  SharedPages --> Archive["Archive<br/>历史档案 / 公开资料"]
  SharedPages --> Connect["Connect<br/>票务 / 售后 / 合作 / 投稿"]

  Shows --> APIEvents["GET /api/events"]
  CDSelect --> APICD["GET /api/cd-items"]
  Archive --> APIContents["GET /api/contents"]
  Connect --> APIConnect["POST /api/connect"]

  APIEvents --> Backend["Go Backend API"]
  APICD --> Backend
  APIContents --> Backend
  APIConnect --> Backend

  Backend --> SeedData["当前阶段<br/>Go 内存 seed 数据"]
  Backend --> FutureDB["P1 数据库<br/>SQLite + GORM ORM<br/>SQL migration"]

  CDSelect --> ExternalShop["外部 Shop<br/>BASE 等"]
  Shows --> ExternalTicket["外部票站<br/>票务代理"]
  Connect --> FutureMail["未来通知<br/>邮件 / 数据库 / 客服系统"]

  Backend --> StaticPreview["本地静态预览<br/>backend/static"]
  Backend --> Deploy["未来部署<br/>阿里云东京轻量服务器"]
```

## 2. 架构原则

### 2.1 设计版本和业务逻辑分离

V2、V3、未来 V4/V5 只是前端视觉 Renderer。

它们共享：

- 同一套页面入口。
- 同一套语言策略。
- 同一套 API。
- 同一套购买路径。
- 同一套 Connect 表单逻辑。

这样客户 Review 时可以自由切换设计，但不会影响后端数据结构。

### 2.2 顶层没有 Shop

当前顶层入口固定为：

```text
Shows / CD 严选 / Archive / Connect
```

购买路径放在 CD/黑胶单品卡片里：

```text
CD 严选 -> 单品卡片 -> 点击此处购买 -> 外部 Shop
```

演出票务也可以跳外部票站。

### 2.3 Connect 是统一问题入口

Connect 承接：

- 票务问题。
- 外部购买后未收到货。
- CD/黑胶购买咨询。
- 发货问题。
- 合作。
- 投稿。

未来如果做站内支付或订单系统，再从 Connect 扩展到完整售后流程。

### 2.4 数据库选型原则

当前 P1 选择：

```text
SQLite + GORM + SQL migration
```

原因：

- 当前数据主要是演出、CD/黑胶、Archive、Connect 消息，读多写少。
- 购买和票务都先跳外部平台，不需要站内订单和库存锁定。
- 阿里云轻量服务器上，SQLite 单文件更容易备份和迁移。
- Go 后端可以通过 GORM 快速接入表结构，同时保留未来迁移 PostgreSQL 的可能。

未来升级 PostgreSQL 的触发条件：

- 做站内支付、订单、库存。
- 多人后台频繁同时编辑。
- 多实例部署。
- 复杂报表或全文搜索。
- SQLite 写入锁或备份策略无法满足运营。

## 3. 前端版本切换图

```mermaid
flowchart LR
  Query["?design=v2 / ?design=v3"] --> Resolver["版本解析"]
  Dropdown["右上角版本下拉框"] --> Resolver
  LocalStorage["本地记忆上次版本"] --> Resolver

  Resolver --> Current["当前 designVariant"]

  Current --> RenderV2["渲染 V2"]
  Current --> RenderV3["渲染 V3"]
  Current --> RenderFuture["渲染未来版本"]

  RenderV2 --> SameData["共享 API 数据"]
  RenderV3 --> SameData
  RenderFuture --> SameData
```

## 4. 后端 API 图

```mermaid
flowchart TB
  Frontend["Frontend<br/>V2 / V3 / Future"] --> Health["GET /api/health"]
  Frontend --> Events["GET /api/events"]
  Frontend --> CDItems["GET /api/cd-items"]
  Frontend --> Contents["GET /api/contents"]
  Frontend --> ConnectPost["POST /api/connect"]

  Events --> Owned["ownedEvents<br/>LIVE LIFE 自主演出置顶"]
  Events --> Recommended["recommendedEvents<br/>推荐 / 档案演出"]

  CDItems --> AllItems["items<br/>完整列表"]
  CDItems --> CD["cd<br/>CD 分类"]
  CDItems --> Vinyl["vinyl<br/>黑胶分类"]

  Contents --> ArchiveNotes["Archive notes<br/>历史资料 / 说明"]

  ConnectPost --> Validation["表单验证"]
  Validation --> Accepted["当前返回 accepted"]
  Accepted --> FutureStorage["未来写入数据库或发送邮件"]
```

## 5. 部署演进图

```mermaid
flowchart TD
  Local["当前本地开发<br/>localhost:8080"] --> Static["Go 托管静态页面<br/>backend/static"]
  Local --> API["Go API"]

  API --> Seed["内存 seed 数据"]

  Seed --> P1["P1<br/>SQLite 持久化<br/>GORM + migration 文件"]
  Static --> P1

  P1 --> Server["阿里云东京轻量服务器<br/>/opt/livelife/app"]
  Server --> Nginx["Nginx 反向代理"]
  Server --> Backup["SQLite / uploads 备份"]

  Nginx --> Domain["未来域名 / HTTPS"]
```

## 6. 当前文档关系

```mermaid
flowchart LR
  Overview["architecture-overview.md"] --> Requirements["requirements-analysis.md"]
  Overview --> BackendDoc["backend-detailed-design.md"]
  Overview --> FrontendDoc["frontend-design-variants.md"]
  Overview --> LocalDev["local-development.md"]
  Overview --> Deployment["alicloud-tokyo-p0-deployment.md"]

  FrontendDoc --> V2Doc["product-architecture-and-ui-approval.md"]
  FrontendDoc --> V3Doc["v3-design-approval.md"]

  Overview --> Roadmap["documentation-roadmap-for-approval.md"]
```
