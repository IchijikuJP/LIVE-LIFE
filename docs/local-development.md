# LiveLife 本地开发说明

当前本地优先目标：不依赖域名、不依赖阿里云公网，先把页面和 API 跑起来。

## 1. 当前可立即运行的方式

因为这台 Windows 当前没有可用的 `npm`，所以先用 Go 后端启动一个本地预览页。

```bash
cd backend
go run ./cmd/server
```

打开：

```text
http://localhost:8080
```

API：

```text
GET  http://localhost:8080/api/health
GET  http://localhost:8080/api/events
GET  http://localhost:8080/api/contents
POST http://localhost:8080/api/join
```

## 2. React/Vite 前端

正式前端源码在：

```text
frontend/
```

等 Node/npm 可用后执行：

```bash
cd frontend
npm install
npm run dev
```

打开：

```text
http://localhost:5173
```

Vite 会把 `/api` 代理到：

```text
http://127.0.0.1:8080
```

所以前端和后端本地联调时，需要同时运行：

```bash
cd backend
go run ./cmd/server
```

## 3. 当前目录分工

```text
backend/
  cmd/server/main.go      Go API
  static/                 临时本地预览页

frontend/
  src/                    React/Vite/Tailwind 正式前端源码

docs/
  alicloud-tokyo-p0-deployment.md
  alicloud-tokyo-p0-status.md
  local-development.md
```

## 4. 下一步建议

1. 先用 Go 本地预览页确认页面和 API 流程。
2. 安装可用的 Node/npm 后跑 React/Vite。
3. 给 Go 后端接 SQLite。
4. 把静态 seed 数据替换成数据库数据。
5. 增加后台管理页面。
