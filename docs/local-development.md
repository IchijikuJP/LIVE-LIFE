# LIVE LIFE 本地开发说明

当前本地优先目标：不依赖域名、不依赖阿里云公网，先把页面结构、API、活动内容和 CD 严选流程跑起来。

## 1. 本地启动

在项目根目录运行：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\start-local-backend.ps1
```

打开：

```text
http://localhost:8080
```

这个脚本会自动做几件事：

- 把 Go 编译缓存放到项目内 `.cache/go-build`，避免写到系统目录时遇到权限问题。
- 把运行日志写入 `logs/local-backend.out.log`。
- 切换到 `backend/` 目录，因为后端会从 `backend/static/` 读取本地预览页面。

## 2. 当前 API

```text
GET  http://localhost:8080/api/health
GET  http://localhost:8080/api/events
GET  http://localhost:8080/api/cd-items
GET  http://localhost:8080/api/contents
POST http://localhost:8080/api/connect
```

兼容保留：

```text
POST http://localhost:8080/api/join
```

说明：

- `/api/events` 会返回 `ownedEvents` 和 `recommendedEvents`，方便前端把 LIVE LIFE 自主演出固定在最上面。
- 顶层不再有独立 Shop 入口，也不再提供 `/api/shop-items`。
- `CD 严选` 内部分为 `CD` 和 `黑胶` 两类。
- 单品详情卡片里提供 `点击此处购买` 按钮，跳转到 BASE 等外部 Shop。
- `/api/connect` 是统一联系入口，不再叫 Join us。它可以承接票务、外部购买后未收到货、发货、投稿、售后和合作问题。
- 所有对外展示数据都带 `brand: "LIVE LIFE"`，三语言文案使用 `titleI18n`、`summaryI18n` 等字段。

## 3. 三语言策略

默认语言是中文，同时支持：

```text
中文 / 日本語 / English
```

后端返回三语言数据，前端负责根据当前语言选择展示。这样以后接数据库时，只要数据库继续保存同样结构，页面语言切换逻辑不用重写。

英语界面的导航和主要入口文案使用全大写，例如：

```text
SHOWS / CD SELECT / ARCHIVE / CONNECT
```

## 4. React/Vite 前端

正式前端源码在：

```text
frontend/
```

当前静态预览由 Go 后端直接托管 `backend/static/`。等 Node/npm 环境可用后，可以运行 React/Vite：

```bash
cd frontend
npm install
npm run dev
```

Vite 会把 `/api` 代理到：

```text
http://127.0.0.1:8080
```

## 5. 登录注册暂缓

目前先不做登录注册。原因是购买流程改为“外部 Shop 跳转优先”：

- 如果单品详情页只跳 BASE 等外部 Shop，LIVE LIFE 站内可能不需要注册。
- 如果之后改成站内支付，就需要订单、支付状态、发货状态和售后入口。
- 如果用户在外部平台付款后没有收到货，现阶段先通过 `Connect` 表单人工处理。

所以现在先把 `CD 严选 -> 单品详情 -> 外部购买按钮` 这条路径设计清楚，等确认是否站内支付后再设计账户系统。

## 6. V2 UI 审批方向

审批文档：

```text
docs/product-architecture-and-ui-approval.md
```

V2 方向：

- 首页参考 Nintendo Systems 的网格、Schedule 面板和强品牌节奏，但不照抄。
- 不使用代码墙，也不使用二进制纹理。
- 背景纹理改成经典英摇/另类摇滚乐队与艺人名罗列，用作文化参考文本，页面上避免写成合作或授权关系。
- 音乐制作数据流改成抽象音轨图、波形块、采样片段和轨道线，不写具体歌曲名。
- 色系以未来 LIVE LIFE 左侧黄色图标为前提，使用米纸底、黑、亮黄、酸蓝和红色形成冲突感。
