# LIVE LIFE GitHub 多人协作工作流

状态：准备启用
最后更新：2026-06-11

## 1. 目标

多人协作以后，GitHub 负责：

- 保存主要代码。
- 让外部开发者通过 Pull Request 提交修改。
- 让负责人 review 后再合并。
- 用 GitHub Actions 自动跑后端测试和前端构建。

阿里云负责：

- 运行线上 / 预览环境。
- 保存 SQLite、uploads、logs、backups 等运行时数据。
- 接收已经构建好的 release 产物。

## 2. 当前远程仓库状态

当前本地已有：

```text
aliyun  ssh://admin@47.74.8.10:2222/opt/livelife/git/livelife.git
```

用途：

```text
部署到阿里云服务器
```

还缺：

```text
origin  GitHub 私有仓库
```

你现在可以做的事：

1. 登录 GitHub。
2. 新建一个 Private 私有仓库。
3. 仓库名建议：`live-life` 或 `LIVE-LIFE`。
4. 不要勾选 README / .gitignore / license。
5. 把仓库地址发给我。

## 3. 目标远程配置

最终应该是：

```text
origin  git@github.com:<OWNER>/<REPO>.git
aliyun  ssh://admin@47.74.8.10:2222/opt/livelife/git/livelife.git
```

区别：

```text
origin
  多人协作、Pull Request、代码 review、CI。

aliyun
  部署远程。只有负责人或部署流程使用。
```

## 4. 分支规则

```text
main
  主要代码基线。只放负责人确认过的阶段性稳定版本。

develop
  日常生产 / 预览部署分支。平时功能 PR 合并到这里。

feature/*
  功能开发分支。

fix/*
  修复分支。

docs/*
  文档分支。
```

默认 PR 目标：

```text
develop
```

main 更新方式：

```text
负责人确认阶段稳定后，再把 develop 合并到 main。
```

## 5. 外部开发者工作方式

开发者第一次：

```bash
git clone git@github.com:<OWNER>/<REPO>.git
cd live-life
git checkout develop
```

开始做功能：

```bash
git pull origin develop
git checkout -b feature/example
```

提交：

```bash
git add .
git commit -m "实现某个功能"
git push origin feature/example
```

然后在 GitHub 上开 Pull Request：

```text
base: develop
compare: feature/example
```

## 6. 负责人 Review 规则

合并 PR 前至少看：

- 是否改动了固定需求。
- 是否改动了后端 API 契约。
- 是否改动了数据库表结构。
- 是否影响部署 release 产物。
- GitHub Actions 是否通过。
- 是否有必要更新 docs。

如果改了后端业务逻辑，必须检查：

```text
backend/internal/domain
backend/internal/application
backend/internal/infrastructure
backend/internal/interfaces
docs/backend-detailed-design.md
docs/database-schema-draft.md
```

## 7. 现在工作流和 GitHub 工作流的区别

以前：

```text
单人本地改
  -> 本地 commit
  -> 本地构建 release
  -> push aliyun master
  -> 服务器部署
```

以后：

```text
多人各自开分支
  -> push GitHub
  -> Pull Request
  -> GitHub CI
  -> 负责人 review
  -> 合并 develop
  -> 负责人构建 release
  -> push aliyun develop
  -> 服务器部署
```

核心变化：

- GitHub 成为协作中心。
- 阿里云不再承担代码 review。
- 阿里云仍然只接收构建好的 release 产物。
- 1GB 内存限制仍然有效，服务器不跑前端构建。

## 8. 当前需要你做的事

现在你可以同时做：

1. 创建 GitHub 私有仓库。
2. 把仓库地址发给我。
3. 之后我会添加 `origin`。
4. 我会把当前代码推到 GitHub。
5. 再建立或确认 `main` / `develop` 分支。
6. 以后其他人从 GitHub 拉代码和提 PR。
