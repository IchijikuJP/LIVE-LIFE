# LIVE LIFE Git Push 部署说明

这套 P0 部署流程用于实现：在本地 VSCode 里提交代码并 `git push` 后，自动同步到阿里云东京服务器。

```text
VSCode / 本地 Git
  -> git push aliyun master
  -> 阿里云服务器上的 bare Git 仓库（SSH 2222 端口）
  -> post-receive hook
  -> /opt/livelife/app
  -> 后台部署脚本
  -> 复制本地预构建 release 产物
  -> 重启后端进程
  -> Nginx 对外提供 live-life.asia
```

## 为什么这样做

- 符合当前需求：本地 push 后同步到阿里云。
- 暂时不需要 GitHub Actions、GitLab CI 或其他 CI/CD 平台。
- 当前服务器只有 1 GiB 内存，不适合在服务器上构建 Go 和 Node 前端产物。
- 服务器只运行本地已经构建好的 Linux 后端二进制和 Nginx 静态站点。
- 运行时数据不进 Git，保留在 `/opt/livelife/data`、`/opt/livelife/uploads`、`/opt/livelife/logs`。

## 服务器一次性配置

把 `scripts/setup-aliyun-git-deploy.sh` 上传到服务器，例如：

```text
/home/admin/setup-aliyun-git-deploy.sh
```

然后在服务器上运行：

```bash
bash ~/setup-aliyun-git-deploy.sh master
```

为了让本地 Git push 稳定连接服务器，使用单独的 SSH 端口 `2222`。`22` 保留给 Workbench / 普通 SSH，`2222` 给本地 Git 部署使用。

```bash
echo 'Port 22' | sudo tee /etc/ssh/sshd_config.d/99-livelife-git-deploy.conf
echo 'Port 2222' | sudo tee -a /etc/ssh/sshd_config.d/99-livelife-git-deploy.conf
sudo systemctl reload ssh
sudo ufw allow 2222/tcp
sudo ufw reload
sudo ss -lntp | grep -E ':22|:2222'
```

同时在阿里云实例防火墙里添加：

```text
TCP 2222 0.0.0.0/0
```

配置完成后，退出 Workbench / SSH，再重新连接一次。

## 本地一次性配置

在本地项目目录运行：

```powershell
git remote add aliyun ssh://admin@47.74.8.10:2222/opt/livelife/git/livelife.git
```

如果服务器登录用户不是 `admin`，替换成实际用户名。

当前项目已经配置了 `aliyun` 远程仓库，正常情况下不需要重复添加。

## 日常部署

在 VSCode 终端里运行：

```powershell
.\scripts\build-release.ps1
git status
git add .
git commit -m "Update LIVE LIFE release"
git push aliyun master
```

`git push` 应该很快返回，并提示部署已经在服务器后台启动。服务器部署时不会运行 `npm build`、`go build` 或 Docker 镜像构建。

如需查看服务器部署日志：

```bash
tail -f /opt/livelife/logs/deploy.log
```

## 部署后检查

服务器本机检查：

```bash
cd /opt/livelife/app
cat /opt/livelife/runtime/livelife-api.pid
curl -I http://127.0.0.1:8080/api/health
curl -I http://127.0.0.1
```

公网检查：

```text
https://live-life.asia
https://live-life.asia/api/health
```

## 以后什么时候升级到 GitHub Actions

当前阶段先用本地 push 到服务器的方式。等项目需要 Pull Request、远程自动测试、部署审批、多人协作或失败通知时，再升级到 GitHub Actions。

升级后流程会变成：

```text
本地 push 到 GitHub
  -> GitHub Actions 构建和测试
  -> GitHub Actions SSH 到阿里云
  -> 执行同一套发布命令
```
