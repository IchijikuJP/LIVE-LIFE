# LIVE LIFE

LIVE LIFE MVP 项目。当前目标是先做一个可以本地预览的东京音乐入口：

- 我们自己的演出和推荐演出情报
- CD 页面
- Shop 页面
- 首页内容摘要
- Connect 联系入口

## 本地预览

当前 Windows 环境暂时没有可用的 `npm`，所以先用 Go 后端直接提供本地静态预览页。

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\start-local-backend.ps1
```

打开：

```text
http://localhost:8080
```

健康检查：

```text
http://localhost:8080/api/health
```

## 当前技术栈

- 前端预览：Go 后端托管的静态页面
- 正式前端骨架：React + Vite + TypeScript + Tailwind CSS
- 后端：Go
- 数据库计划：MVP 先 SQLite，后续可升级 PostgreSQL
- 部署计划：Docker Compose + Nginx + 阿里云东京节点

## 项目目录

```text
backend/      Go API 和本地静态预览页
frontend/     React/Vite/Tailwind 正式前端骨架
docs/         部署文档、产品审批稿和本地开发说明
scripts/      本地启动脚本
```
