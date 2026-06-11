# LIVE LIFE 部署流程：GitHub 协作后如何兼容 1GB 阿里云服务器

状态：目标流程说明
最后更新：2026-06-11

## 1. 以前的工作流

之前是单人本地开发，所以流程比较直接：

```text
本地修改代码
  -> 本地构建 release 产物
  -> git push aliyun master
  -> 阿里云 bare repo post-receive hook
  -> 复制 deploy/release 产物
  -> 重启 Go 后端
  -> Nginx 提供前端静态文件
```

本地命令：

```powershell
.\scripts\build-release.ps1
git status
git add .
git commit -m "Update LIVE LIFE release"
git push aliyun master
```

这个流程的重点是：服务器不构建。

原因：

- 阿里云东京轻量服务器只有 1GB 内存。
- 不适合在服务器上跑 `npm run build`。
- 不适合在服务器上跑较重的 Go/Node 构建任务。
- 服务器只运行已经构建好的 Linux Go 二进制和前端静态文件。

## 2. 现在准备升级的多人协作工作流

多人协作后，代码不再直接从每个人电脑推到阿里云。

目标流程：

```text
开发者本地改代码
  -> push 到自己的 GitHub 分支
  -> 提 Pull Request 到 develop
  -> GitHub CI 自动测试
  -> 负责人 Review
  -> 合并到 develop
  -> 负责人构建 release 产物
  -> push 到 aliyun develop
  -> 阿里云服务器只复制产物并重启
```

服务器仍然不构建。

## 3. 分支规则

按当前项目要求：

```text
main
  主要代码基线。只放负责人确认过的阶段性稳定版本。

develop
  日常生产 / 预览部署分支。平时多人 PR 合并到这里，阿里云部署也以这里为准。

feature/*
  新功能分支，例如 feature/backend-clean-architecture。

fix/*
  修复分支，例如 fix/connect-validation。

docs/*
  文档分支，例如 docs/github-workflow。
```

规则：

- 普通协作者不能直接 push 到 `main`。
- 普通协作者不能直接 push 到 `develop`。
- 所有功能通过 Pull Request 给负责人 review。
- `develop` 通过 CI 后才能部署。
- `main` 只在阶段确认后由负责人从 `develop` 合并或 cherry-pick。

## 4. GitHub 远程配置

当前本地已经有阿里云远程：

```text
aliyun  ssh://admin@47.74.8.10:2222/opt/livelife/git/livelife.git
```

还需要新增 GitHub 远程：

```powershell
git remote add origin git@github.com:<OWNER>/<REPO>.git
```

或 HTTPS：

```powershell
git remote add origin https://github.com/<OWNER>/<REPO>.git
```

新增后应当是：

```text
origin  GitHub 私有仓库，用于多人协作和 Pull Request
aliyun  阿里云 bare repo，用于部署 release 产物
```

## 5. GitHub Actions 的职责

GitHub Actions 先只做 CI，不直接部署。

当前 `.github/workflows/ci.yml` 会做：

```text
Pull Request 到 main/develop
  -> Go 后端测试：go test ./...
  -> React 前端构建：npm ci && npm run build
```

为什么暂时不让 GitHub Actions 自动部署：

- 现阶段服务器部署还依赖阿里云 bare repo 和 post-receive hook。
- 需要先确认 GitHub 仓库、Secrets、服务器 SSH 权限。
- 1GB 服务器不能构建，所以即使未来自动部署，也必须部署预构建产物。

## 6. 阿里云部署实际推送的内容

`scripts/build-release.ps1` 会生成：

```text
deploy/release/
  backend/
    livelife-api
  frontend/
    index.html
    assets/
      *.js
      *.css
      图片和静态资源
```

具体内容：

- `deploy/release/backend/livelife-api`
  - Linux amd64 Go 二进制。
  - 已经把 Clean Architecture 后端代码编译进一个可执行文件。
  - 服务器只运行这个文件，不需要 Go 源码参与运行。

- `deploy/release/frontend/`
  - Vite 构建后的静态文件。
  - Nginx 直接以静态站点方式提供。

不会推送到服务器运行时目录的数据：

```text
/opt/livelife/data/sqlite/livelife.db
/opt/livelife/uploads
/opt/livelife/logs
/opt/livelife/backups
```

这些是服务器运行时数据，不能进 Git。

## 7. 迁移到 develop 部署分支

当前历史部署文档里使用过 `master`。

多人协作后建议改成：

```bash
bash ~/setup-aliyun-git-deploy.sh develop
```

这会让阿里云 post-receive hook 只在收到 `develop` 分支更新时部署。

本地部署命令也变成：

```powershell
.\scripts\build-release.ps1
git add .
git commit -m "Build LIVE LIFE release"
git push aliyun develop
```

注意：

- 修改仓库里的脚本不会自动改变服务器已经安装好的 hook。
- 要让服务器从 `master` 改为 `develop`，需要在服务器上重新执行 setup 脚本或手动改 hook。

## 8. 未来自动部署升级方案

等 GitHub 协作稳定后，可以升级成：

```text
合并 PR 到 develop
  -> GitHub Actions 构建 release 产物
  -> GitHub Actions 通过 SSH 上传产物到阿里云
  -> 服务器复制产物并重启
```

关键原则仍然不变：

```text
构建发生在 GitHub Actions，不发生在 1GB 阿里云服务器。
```
