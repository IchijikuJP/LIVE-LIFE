# LIVE LIFE GitHub 多人协作工作流

状态：已启用
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

当前本地已有两个远程：

```text
origin  https://github.com/IchijikuJP/LIVE-LIFE.git
aliyun  ssh://admin@47.74.8.10:2222/opt/livelife/git/livelife.git
```

用途：

```text
origin
  GitHub 协作、Pull Request、代码 review、CI。

aliyun
  部署到阿里云服务器。普通协作者不使用。
```

## 3. 分支规则

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

## 4. 外部开发者工作方式

开发者第一次：

```bash
git clone https://github.com/IchijikuJP/LIVE-LIFE.git
cd LIVE-LIFE
git switch develop
```

开始做功能：

```bash
git pull origin develop
git switch -c feature/example
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

三丰老师（`kirori-1`）的详细手册：

```text
docs/collaborator-manual-kirori-1.md
```

## 5. 负责人 Review 规则

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

## 6. 现在工作流和 GitHub 工作流的区别

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

## 7. 当前已完成

```text
main
  已推送到 GitHub。

develop
  已推送到 GitHub。

origin
  已指向 https://github.com/IchijikuJP/LIVE-LIFE.git。

aliyun
  保留为部署远程。
```

## 8. GitHub 分支保护建议

建议给 `main` 和 `develop` 分别添加保护规则。

最低要求：

- Require a pull request before merging
- Require approvals：至少 1 个
- Require status checks to pass before merging
- 选择 GitHub Actions 里的 `Go backend` 和 `React frontend`
- Require conversation resolution before merging
- 不允许 force push
- 不允许删除受保护分支

如果 GitHub 页面暂时还不能选择具体 status check，先让任意 PR 跑一次 GitHub Actions，再回到保护规则里选择对应检查项。

## 9. 保护规则设置步骤

在 GitHub 页面上：

1. 打开 `IchijikuJP/LIVE-LIFE` 仓库。
2. 点顶部 `Settings`。
3. 左侧点 `Branches`。
4. 在 `Branch protection rules` 区域点 `Add branch protection rule`。
5. `Branch name pattern` 填 `main`。
6. 勾选 `Require a pull request before merging`。
7. `Required approvals` 设置为 `1`。
8. 勾选 `Dismiss stale pull request approvals when new commits are pushed`。
9. 勾选 `Require status checks to pass before merging`。
10. 勾选 `Require branches to be up to date before merging`。
11. 在 status checks 里选择：
    - `Go backend`
    - `React frontend`
12. 勾选 `Require conversation resolution before merging`。
13. 不要勾选 `Allow force pushes`。
14. 不要勾选 `Allow deletions`。
15. 点 `Create` 或 `Save changes`。

然后重复一次，`Branch name pattern` 改成：

```text
develop
```

## 10. main 和 develop 的建议差异

`main` 可以更严格：

- 必须 PR。
- 必须 1 个 approval。
- 必须 CI 通过。
- 必须解决所有 conversation。
- 不允许 force push。
- 不允许删除。

`develop` 也建议保持同样规则。当前团队人数少，不建议为了方便直接推 `develop`，否则 PR review 的习惯很容易被破坏。
