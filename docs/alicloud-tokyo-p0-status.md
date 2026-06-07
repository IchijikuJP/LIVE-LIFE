# LIVE LIFE 阿里云东京服务器 P0 配置状态

更新时间：2026-06-06

本文记录已经在阿里云 Workbench 里实际完成的服务器初始化结果。

## 已确认的服务器信息

- 登录用户：`admin`
- 系统版本：Ubuntu 24.04.2 LTS
- 内核版本：Linux 6.8.0-63-generic
- 云厂商：Alibaba Cloud
- 主机名：`livelife-tokyo-1`
- 内存：约 894 MiB
- 系统盘：约 30 GiB，当前已用约 2.9 GiB
- Workbench 内部网卡地址：`172.19.4.46`

备注：

- Workbench 页面没有直接显示公网 IP，公网 IP 仍需在阿里云轻量应用服务器控制台里确认。
- 当前服务器是新环境，适合作为 LIVE LIFE MVP 的 P0 基础服务器。

## 已完成

### 1. 软件源刷新

已执行：

```bash
sudo apt update
```

结果：

- 阿里云 Ubuntu 软件源可访问。
- 软件包索引刷新成功。

### 2. 基础软件安装

已安装：

- `ca-certificates`
- `curl`
- `git`
- `ufw`
- `fail2ban`
- `nginx`
- `certbot`
- `python3-certbot-nginx`
- `docker.io`
- `docker-compose-v2`
- `sqlite3`
- `dnsutils`

验收结果：

```text
Docker version 29.1.3
Docker Compose version 2.40.3
nginx version: nginx/1.24.0 (Ubuntu)
certbot 2.9.0
```

### 3. 服务启用

已启用并确认 active：

- Docker
- Nginx
- Fail2ban

验收命令：

```bash
systemctl is-active docker nginx fail2ban
```

结果：

```text
active
active
active
```

### 4. Docker 用户组

已执行：

```bash
sudo usermod -aG docker admin
```

备注：

- 当前 Workbench 会话未重新登录前，`admin` 可能仍需要用 `sudo docker ...`。
- 重新登录 SSH/Workbench 后，`admin` 用户应可以直接使用 Docker。

### 5. swap 配置

已创建 2 GiB swap：

```text
/swapfile file 2G
```

验收结果：

```text
Mem: 894Mi
Swap: 2.0Gi
```

同时已设置：

```bash
vm.swappiness=10
```

备注：

- 阿里云镜像中原有 sysctl 配置会把 `swappiness` 覆盖成 0。
- 已新增 `/etc/sysctl.d/zz-livelife.conf`，并手动确认当前值为 `10`。

### 6. 项目目录

已创建：

```text
/opt/livelife
/opt/livelife/app
/opt/livelife/data/sqlite
/opt/livelife/uploads
/opt/livelife/backups
/opt/livelife/logs
```

权限：

```text
owner: admin
group: admin
```

### 7. Ubuntu UFW 防火墙

已启用 UFW，当前只放行：

```text
OpenSSH
80/tcp
443/tcp
```

验收结果：

```text
Status: active
OpenSSH ALLOW Anywhere
80/tcp ALLOW Anywhere
443/tcp ALLOW Anywhere
```

## 部分完成 / 需要注意

### 系统升级

已尝试执行系统升级：

```bash
sudo DEBIAN_FRONTEND=noninteractive apt upgrade -y
```

大部分包已升级，但阿里云镜像源有 3 个包返回 404：

- `libapparmor1`
- `apparmor`
- `python3-urllib3`

判断：

- 这是镜像源同步问题，不是服务器配置错误。
- 不影响 Docker、Nginx、Certbot、UFW、Fail2ban 当前使用。

后续建议：

```bash
sudo apt update
sudo DEBIAN_FRONTEND=noninteractive apt upgrade -y --fix-missing
```

如果仍然 404，等阿里云镜像同步后再执行即可。

## 尚未完成

这些需要下一步继续：

1. 在阿里云轻量应用服务器控制台确认公网 IP。
2. 在阿里云控制台确认云防火墙已开放 `22`、`80`、`443`。
3. 注册或确认 `live-life.asia`。
4. 添加 DNS A 记录：
   - `@` -> 服务器公网 IP
   - `www` -> 服务器公网 IP
5. 项目代码准备好后，放入 `/opt/livelife/app`。
6. 创建生产 `.env`。
7. 创建 `docker-compose.yml`。
8. 启动前端和后端容器。
9. 配置 Nginx 站点反向代理。
10. DNS 生效后申请 HTTPS 证书。
11. 创建 SQLite 和上传文件备份脚本。

## 当前结论

服务器 P0 基础环境已经可用。

现在可以进入下一阶段：

- 如果域名已准备好：先做 DNS 和 HTTPS。
- 如果代码还没开始：先在本地项目里搭 React/Vite 前端和 Go/Gin 后端骨架。
- 如果要先放一个临时页面：可以先用 Nginx 默认页或一个静态占位页验证公网访问。
