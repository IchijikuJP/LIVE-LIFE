# LIVE LIFE 本地开发说明

当前本地优先目标：不依赖域名、不依赖阿里云公网，先把页面结构、API 和活动内容跑起来。

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

- 把 Go 编译缓存放到项目内 `.cache/go-build`。
- 把运行日志写入 `logs/local-backend.out.log`。
- 切换到 `backend/` 目录，因为后端会从 `backend/static/` 读取本地预览页。

## 2. 当前 API

```text
GET  http://localhost:8080/api/health
GET  http://localhost:8080/api/events
GET  http://localhost:8080/api/cd-items
GET  http://localhost:8080/api/shop-items
GET  http://localhost:8080/api/contents
POST http://localhost:8080/api/connect
```

兼容保留：

```text
POST http://localhost:8080/api/join
```

说明：

- `/api/events` 会返回 `ownedEvents` 和 `recommendedEvents`，方便前端把 LIVE LIFE 自主演出固定在最上面。
- `/api/cd-items` 和 `/api/shop-items` 已经拆开，避免把演出、CD、Shop 混在一个页面里。
- `/api/connect` 是统一联系入口，不再叫 Join us。它可以承接票务、付款后未收到货、发货、投稿和合作问题。
- 所有对外展示数据都带 `brand: "LIVE LIFE"`，三语言文案使用 `titleI18n`、`summaryI18n` 等字段。

## 3. 三语言策略

默认语言是中文，同时支持：

```text
中文 / 日本語 / English
```

后端返回三语言数据，前端负责根据当前语言选择展示。这样以后接数据库时，只要数据库继续保存同样结构，页面语言切换逻辑不用重写。

## 4. React/Vite 前端

正式前端源代码在：

```text
frontend/
```

当前机器暂时没有可用的 `npm`，所以 React/Vite 还不能本地跑。等 Node/npm 可用后执行：

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

目前先不做登录注册。原因是 Shop 的购买流程还在讨论中：

- 如果只跳外部购买链接，可能不需要注册。
- 如果站内支付，就需要订单、支付状态、发货状态和售后入口。
- 如果要处理“付款了但没收到货”，至少需要订单查询或人工客服表单。

所以现在先把 Shop 和 Connect 拆出来，等购买流程确认后再设计账户系统。
