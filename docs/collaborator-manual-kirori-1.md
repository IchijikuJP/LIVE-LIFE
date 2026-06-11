# LIVE LIFE 协作者 Manual：三丰老师（kirori-1）

目标读者：三丰老师（GitHub：`kirori-1`）
最后更新：2026-06-11

## 1. 这份 Manual 的目的

这份文档用于说明三丰老师参与 LIVE LIFE 项目时的 GitHub 协作方式。

当前项目已经从单人本地开发，切换为 GitHub 多人协作模式。以后代码修改不直接推到主分支，而是通过 Pull Request 交给负责人 review 后合并。

## 2. 仓库地址

GitHub 仓库：

```text
https://github.com/IchijikuJP/LIVE-LIFE
```

本地建议使用 HTTPS clone：

```bash
git clone https://github.com/IchijikuJP/LIVE-LIFE.git
cd LIVE-LIFE
```

如果习惯 SSH，也可以使用：

```bash
git clone git@github.com:IchijikuJP/LIVE-LIFE.git
cd LIVE-LIFE
```

## 3. 分支规则

当前固定分支如下：

```text
main
  稳定主线。只放负责人确认过的阶段性稳定版本。

develop
  日常开发和预览部署基线。普通功能都从这里拉分支，也合并回这里。

feature/*
  新功能分支。

fix/*
  bug 修复分支。

docs/*
  文档修改分支。
```

三丰老师平时不要直接改 `main` 或 `develop`。建议每次都从 `develop` 新建自己的工作分支。

## 4. 第一次开始开发

接受 GitHub 邀请后，先 clone 仓库：

```bash
git clone https://github.com/IchijikuJP/LIVE-LIFE.git
cd LIVE-LIFE
```

切到日常开发基线：

```bash
git switch develop
git pull origin develop
```

确认当前分支：

```bash
git branch --show-current
```

应该看到：

```text
develop
```

如果是第一次在这台电脑使用 Git，建议先设置提交身份：

```bash
git config --global user.name "kirori-1"
git config --global user.email "你的 GitHub 邮箱"
```

## 5. 每次做新任务的流程

先更新 `develop`：

```bash
git switch develop
git pull origin develop
```

从 `develop` 新建分支：

```bash
git switch -c feature/short-task-name
```

分支名例子：

```text
feature/cd-select-detail
feature/connect-form-copy
fix/api-connect-validation
docs/update-readme
```

修改代码或文档后，查看改动：

```bash
git status
git diff
```

提交：

```bash
git add .
git commit -m "说明这次修改的内容"
```

推送自己的分支：

```bash
git push origin feature/short-task-name
```

然后到 GitHub 页面创建 Pull Request：

```text
base: develop
compare: feature/short-task-name
```

Pull Request 创建后，请在描述里写清楚：

- 这次改了什么。
- 改动影响哪些页面、API 或数据结构。
- 自己跑过哪些测试。
- 希望负责人重点 review 哪里。

## 6. 本地测试命令

后端测试：

```bash
cd backend
go test ./...
```

前端第一次安装依赖：

```bash
cd frontend
npm ci
```

前端本地开发：

```bash
cd frontend
npm run dev
```

前端生产构建检查：

```bash
cd frontend
npm run build
```

如果只改文档，可以不跑前后端构建，但 PR 描述里请写明“仅文档修改”。

## 7. 当前项目结构

```text
backend/
  Go 后端 API。现在已经按 Clean Architecture 拆分。

frontend/
  React + Vite + TypeScript 前端。

docs/
  需求、架构、数据库、设计方案、协作流程和部署说明。

scripts/
  本地启动、部署辅助脚本。

deploy/
  部署相关文件。release 产物由负责人控制。
```

后端重点目录：

```text
backend/internal/domain
  业务实体和规则。

backend/internal/application
  应用服务和用例。

backend/internal/infrastructure
  SQLite / GORM 等外部实现。

backend/internal/interfaces
  HTTP API 层。
```

## 8. 需求边界

当前固定业务方向：

- 品牌展示统一为 `LIVE LIFE`。
- 语言为中文、日本語、English，默认中文。
- 顶层导航为 Shows / CD 严选 / Archive / Connect。
- 不设置顶层 Shop 页面。
- CD 严选内部区分 CD 和黑胶。
- 单品详情页通过外部购买按钮跳转 BASE 等外部 Shop。
- Shows 里 LIVE LIFE 自主演出固定展示在推荐演出上方。
- Connect 是统一联系入口，可用于票务、购买未收到货、合作、投稿等问题。

如果修改会影响这些固定需求，请先在 PR 里说明，不要直接改掉。

## 9. 后端 API 注意事项

当前前端依赖这些 API：

```text
GET  /api/health
GET  /api/events
GET  /api/cd-items
GET  /api/contents
POST /api/connect
POST /api/join
```

其中：

- `/api/events` 返回演出情报，并区分 `ownedEvents` 和 `recommendedEvents`。
- `/api/cd-items` 返回 CD 严选数据，并区分 `cd` 和 `vinyl`。
- `/api/connect` 保存联系表单消息。
- `/api/shop-items` 已不再作为顶层 API 使用。

如果要改 API 字段名、删除字段、改变返回结构，请先和负责人确认，因为这会影响前端和未来后台。

## 10. 数据库和部署注意事项

当前 MVP 使用：

```text
SQLite + GORM
```

这样做是因为当前数据量小，阿里云服务器内存只有 1GB，MVP 阶段先保持简单。

普通协作者不要直接操作阿里云服务器，也不要使用 `aliyun` remote。`aliyun` 只用于负责人部署。

协作者日常只需要：

```text
origin = GitHub 仓库
```

不要提交这些运行时内容：

```text
data/*.db
logs/
deploy/release/
node_modules/
frontend/dist/
```

## 11. PR 合并前检查

提交 PR 前建议确认：

```bash
git status
```

工作区不应该有无关文件。

如果改了后端：

```bash
cd backend
go test ./...
```

如果改了前端：

```bash
cd frontend
npm run build
```

如果改了文档：

请确认链接路径能从 GitHub 页面点开。

## 12. 常见问题

### 我应该把 PR 合并到 main 还是 develop？

普通功能统一合并到 `develop`。`main` 由负责人在阶段稳定后更新。

### 我可以直接 push 到 develop 吗？

不建议。保护规则开启后也不应该允许直接 push。请开分支并提交 PR。

### 我需要知道阿里云服务器密码吗？

不需要。普通协作只走 GitHub。

### 设计 V2 / V3 / 之后 V4 会影响后端吗？

原则上不影响。视觉方案可以变化，但 Shows、CD 严选、Archive、Connect 和 API 数据结构要保持稳定。

### 如果我只改一句文案，也要 PR 吗？

是的。多人协作后，所有改动都通过 PR 留记录。
