# LIVE LIFE

LIVE LIFE MVP 项目。当前目标是先做一个可以本地预览的东京音乐入口：

- 我们自己的演出和推荐演出情报
- CD 严选：内部再分 CD 和黑胶
- 单品详情页里的外部购买按钮，例如跳转到 BASE
- 首页内容摘要
- Connect 联系入口

## 给协作者

目标读者：三丰老师（GitHub：`kirori-1`）

请先阅读：

- [docs/collaborator-manual-kirori-1.md](docs/collaborator-manual-kirori-1.md)
- [docs/github-collaboration-workflow.md](docs/github-collaboration-workflow.md)

当前协作规则：

- `main` 是稳定主线。
- `develop` 是日常开发和预览部署基线。
- 普通修改请从 `develop` 新建 `feature/*`、`fix/*` 或 `docs/*` 分支。
- 修改完成后向 `develop` 提 Pull Request，负责人 review 后合并。
- 阿里云 `aliyun` remote 只用于部署，不作为普通协作者入口。

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

## 当前设计审批重点

新的 UI 方向还在审批中，详见：

```text
docs/product-architecture-and-ui-approval.md
```

审批重点：

- 顶层导航不再放 Shop。
- `CD` 入口改为 `CD 严选`。
- `CD 严选` 内部分成 `CD` 和 `黑胶`。
- 单品详情页提供 `点击此处购买` 按钮，跳转外部 Shop。
- 首页纹理从代码/二进制改成乐队名罗列和音乐制作数据流。
