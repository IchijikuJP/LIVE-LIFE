# LIVE LIFE 阿里云东京轻量服务器 P0 部署文档

本文对应架构图里的 MVP 版本，用来完成第一阶段服务器基础设施配置。

目标不是一次性把所有业务功能上线，而是先把下面这条链路打通：

```text
用户浏览器
  -> 域名 live-life.asia
  -> HTTPS 443
  -> 阿里云轻量应用服务器防火墙
  -> Ubuntu 上的 Nginx
  -> Docker Compose 里的前端和后端服务
  -> SQLite 数据库和上传文件目录
```

## 0. 当前服务器信息

根据架构图，目前计划使用：

- 云厂商：阿里云
- 地域：日本东京
- 产品：轻量应用服务器
- 系统：Ubuntu 24.04
- 配置：2 vCPU / 1 GiB 内存 / 30 GiB 磁盘
- 域名：`live-life.asia`
- 前端：React + Vite + TypeScript + Tailwind CSS
- 后端：Go + Gin
- 数据库：MVP 阶段先用 SQLite，后续再迁移 PostgreSQL
- 反向代理：Nginx
- HTTPS：Let's Encrypt + Certbot
- 部署方式：Docker Compose

备注：

- 1 GiB 内存可以跑 MVP，但余量很小，所以必须加 swap，并且暂时不建议把 PostgreSQL 也放在这台服务器上。
- SQLite 适合早期内容展示、活动信息、报名表单、后台录入等低并发场景。只要备份策略做好，MVP 阶段够用。
- 后续如果站内购买、订单、支付和多人后台管理成为核心功能，再升级服务器并迁移 PostgreSQL。

## 1. MVP 推荐部署结构

推荐使用“宿主机 Nginx + Docker Compose 应用服务”的结构：

```text
公网用户
  |
  | HTTPS 443
  v
阿里云安全组 / 防火墙
  |
  v
Ubuntu 24.04 宿主机
  |
  v
Nginx
  |-- /              -> 127.0.0.1:3000  前端容器
  |-- /api/          -> 127.0.0.1:8080  后端 Go + Gin 容器
  |-- /uploads/      -> 127.0.0.1:8080  后端上传文件服务

数据文件：
  SQLite: /opt/livelife/data/sqlite/livelife.db
  上传文件: /opt/livelife/uploads
  备份文件: /opt/livelife/backups
  日志文件: /opt/livelife/logs
```

为什么这样设计：

- Nginx 放在宿主机上，申请 HTTPS 证书和续期最简单。
- 前端和后端容器只绑定到 `127.0.0.1`，外网不能直接访问容器端口，安全性更好。
- 所有项目数据都集中在 `/opt/livelife`，后续备份、迁移、排查都更清楚。
- Docker Compose 足够支撑 MVP，之后要迁移到更复杂的部署方式也容易。

## 2. 阿里云控制台需要确认的项目

登录阿里云控制台后，优先确认这些信息。

### 2.1 服务器实例

需要确认：

- 实例地域是日本东京。
- 系统镜像是 Ubuntu 24.04。
- 实例状态是运行中。
- 记录公网 IPv4 地址。
- 记录登录用户名，一般是 `root` 或创建时指定的用户。

备注：

- 不要点击“重置系统”“释放实例”“更换镜像”这类操作，除非已经明确要重装服务器。
- 如果服务器不是 Ubuntu 24.04，也可以继续部署，但命令可能需要小调整。P0 文档默认以 Ubuntu 24.04 为准。

### 2.2 阿里云防火墙 / 安全组

轻量应用服务器通常有自己的防火墙规则页面。需要开放：

| 端口 | 协议 | 用途 | 是否必须 |
| --- | --- | --- | --- |
| 22 | TCP | SSH 登录服务器 | 必须 |
| 80 | TCP | HTTP、Let's Encrypt 域名验证 | 必须 |
| 443 | TCP | HTTPS 正式访问 | 必须 |

备注：

- `22` 端口前期可以对所有 IP 开放，配置完成后建议限制为你的固定 IP。
- `80` 端口即使最终全站 HTTPS，也要保留，用于证书续期和 HTTP 跳转。
- 不建议开放 `3000`、`8080`、数据库端口。前端和后端由 Nginx 反向代理访问。

### 2.3 域名和 DNS

如果域名 `live-life.asia` 已经注册，需要添加 DNS 解析：

| 主机记录 | 记录类型 | 记录值 |
| --- | --- | --- |
| `@` | A | 服务器公网 IPv4 |
| `www` | A | 服务器公网 IPv4 |

建议 TTL：

- 配置初期：`600` 秒
- 稳定后：可以保持 `600`，也可以改成更长

备注：

- `@` 代表根域名，也就是 `live-life.asia`。
- `www` 代表 `www.live-life.asia`。
- DNS 没生效之前不要申请 HTTPS 证书，否则 Certbot 会失败。

## 3. 第一次 SSH 登录服务器

在本地终端执行：

```bash
ssh root@服务器公网IP
```

例如：

```bash
ssh root@123.123.123.123
```

如果阿里云提供的是其他用户名，就替换 `root`。

首次登录后建议设置主机名：

```bash
sudo hostnamectl set-hostname livelife-tokyo-1
```

备注：

- 主机名只是方便以后识别服务器，不影响网站访问。
- 如果 SSH 提示是否信任服务器指纹，确认 IP 是你的阿里云服务器后输入 `yes`。
- 如果 SSH 登录失败，优先检查阿里云防火墙是否开放 `22`，以及服务器是否运行中。

## 4. 更新系统和安装基础软件

更新软件源和系统包：

```bash
sudo apt update
sudo apt upgrade -y
```

安装基础工具：

```bash
sudo apt install -y ca-certificates curl git ufw fail2ban nginx certbot python3-certbot-nginx docker.io docker-compose-v2 sqlite3
```

启用服务：

```bash
sudo systemctl enable --now docker
sudo systemctl enable --now nginx
sudo systemctl enable --now fail2ban
```

检查版本：

```bash
docker --version
docker compose version
nginx -v
sqlite3 --version
```

备注：

- `docker.io` 是 Ubuntu 官方仓库里的 Docker 包，MVP 阶段够用。
- `docker-compose-v2` 提供的是 `docker compose` 命令，不是旧版 `docker-compose`。
- `fail2ban` 用来降低 SSH 被暴力尝试的风险。
- 如果 `apt upgrade` 提示需要重启，可以在安装完成后执行 `sudo reboot`，等 1 分钟再重新 SSH。

## 5. 创建部署用户

建议不要长期直接用 `root` 部署应用。创建 `deploy` 用户：

```bash
sudo adduser deploy
sudo usermod -aG sudo deploy
sudo usermod -aG docker deploy
```

切换到 `deploy` 用户测试：

```bash
su - deploy
docker ps
```

备注：

- `deploy` 加入 `sudo` 组后，可以执行管理员命令。
- `deploy` 加入 `docker` 组后，可以运行 Docker 命令。
- 刚加入 `docker` 组后，可能需要重新登录 SSH 才生效。

## 6. 给 1 GiB 内存服务器添加 swap

1 GiB 内存跑 Docker、Node 构建、Go 后端都比较紧张。建议添加 2 GiB swap。

先检查是否已有 swap：

```bash
swapon --show
```

如果没有输出，创建 swap 文件：

```bash
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

写入开机自动挂载：

```bash
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

降低 swap 使用倾向：

```bash
echo 'vm.swappiness=10' | sudo tee /etc/sysctl.d/99-livelife.conf
sudo sysctl --system
```

检查结果：

```bash
free -h
swapon --show
```

备注：

- swap 不是性能优化，而是防止内存不够时进程直接崩掉。
- 如果以后升级到 2 GiB 或 4 GiB 内存，swap 仍然可以保留。

## 7. 配置 Ubuntu 防火墙 UFW

阿里云控制台防火墙负责云层面的入口，Ubuntu 自己的 UFW 负责系统层面的入口。两个都配置，比较稳。

执行：

```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

期望看到：

```text
22/tcp 或 OpenSSH ALLOW
80/tcp ALLOW
443/tcp ALLOW
```

备注：

- 执行 `sudo ufw enable` 前必须先允许 SSH，否则可能把自己锁在服务器外。
- 不要开放 `3000` 和 `8080` 到公网。
- 如果未来更改 SSH 端口，要先放行新端口，再关闭旧端口。

## 8. 创建项目目录

统一使用 `/opt/livelife`：

```bash
sudo mkdir -p /opt/livelife/app
sudo mkdir -p /opt/livelife/data/sqlite
sudo mkdir -p /opt/livelife/uploads
sudo mkdir -p /opt/livelife/backups
sudo mkdir -p /opt/livelife/logs
sudo chown -R deploy:deploy /opt/livelife
```

推荐目录结构：

```text
/opt/livelife
  app/
    frontend/
    backend/
    docker-compose.yml
    .env
    backup-sqlite.sh
  data/
    sqlite/
      livelife.db
  uploads/
  backups/
  logs/
```

备注：

- `app` 放代码和部署文件。
- `data/sqlite` 放数据库，不能随便删除。
- `uploads` 放活动海报、CD 封面、内容图片等上传文件。
- `backups` 放本机备份，但上线后还要定期下载到本地或对象存储。
- `logs` 放备份日志或应用日志。

## 9. 环境变量文件 `.env`

项目代码准备好后，在 `/opt/livelife/app/.env` 创建环境变量：

```bash
APP_ENV=production
APP_BASE_URL=https://live-life.asia
API_BASE_URL=https://live-life.asia/api

BACKEND_PORT=8080
DATABASE_PATH=/app/data/livelife.db
UPLOAD_DIR=/app/uploads

SESSION_SECRET=replace-with-long-random-string

STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=
```

生成 `SESSION_SECRET`：

```bash
openssl rand -base64 48
```

备注：

- `.env` 不要提交到公开 Git 仓库。
- `SESSION_SECRET` 用于登录态、后台会话等敏感逻辑，必须使用随机长字符串。
- Stripe 相关变量在 P3 支付阶段再填写。
- 如果先不做登录后台，也建议提前保留这些变量，后续接入更顺。

## 10. Docker Compose 基础模板

项目代码存在后，在 `/opt/livelife/app/docker-compose.yml` 使用类似结构：

```yaml
services:
  frontend:
    build:
      context: ./frontend
    restart: unless-stopped
    ports:
      - "127.0.0.1:3000:80"
    depends_on:
      - backend

  backend:
    build:
      context: ./backend
    restart: unless-stopped
    env_file:
      - .env
    ports:
      - "127.0.0.1:8080:8080"
    volumes:
      - ../data/sqlite:/app/data
      - ../uploads:/app/uploads
      - ../logs:/app/logs
```

部署命令：

```bash
cd /opt/livelife/app
docker compose up -d --build
docker compose ps
```

查看日志：

```bash
docker compose logs -f --tail=100
```

备注：

- `127.0.0.1:3000:80` 表示只有服务器本机可以访问前端容器端口。
- `127.0.0.1:8080:8080` 表示只有服务器本机可以访问后端容器端口。
- 外部用户访问的是 Nginx 的 `80/443`，不是容器端口。
- `restart: unless-stopped` 表示服务器重启后容器会自动恢复。

## 11. Nginx 反向代理配置

创建配置文件：

```bash
sudo nano /etc/nginx/sites-available/livelife
```

写入：

```nginx
server {
    listen 80;
    server_name live-life.asia www.live-life.asia;

    client_max_body_size 20m;

    location /api/ {
        proxy_pass http://127.0.0.1:8080/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8080/uploads/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用站点：

```bash
sudo ln -s /etc/nginx/sites-available/livelife /etc/nginx/sites-enabled/livelife
sudo nginx -t
sudo systemctl reload nginx
```

备注：

- `client_max_body_size 20m` 控制上传文件大小，MVP 阶段用于海报、封面图基本够用。
- `/api/` 转发到后端。
- `/uploads/` 转发到后端，由后端控制文件访问。
- `/` 转发到前端。
- 每次修改 Nginx 配置都先执行 `sudo nginx -t`，确认无误再 reload。

## 12. 申请 HTTPS 证书

先确认 DNS 已生效：

```bash
dig +short live-life.asia
dig +short www.live-life.asia
```

如果服务器没有 `dig`：

```bash
sudo apt install -y dnsutils
```

申请证书：

```bash
sudo certbot --nginx -d live-life.asia -d www.live-life.asia
```

测试自动续期：

```bash
sudo certbot renew --dry-run
```

备注：

- DNS 返回的 IP 必须是阿里云服务器公网 IP。
- Certbot 会自动修改 Nginx 配置，把 HTTP 升级到 HTTPS。
- Let's Encrypt 证书有效期通常是 90 天，Certbot 会自动续期。
- 如果证书申请失败，不要连续反复重试太多次，可能触发频率限制。先确认 DNS、80 端口、防火墙和 Nginx。

## 13. SQLite 和上传文件备份

创建备份脚本：

```bash
nano /opt/livelife/app/backup-sqlite.sh
```

写入：

```bash
#!/usr/bin/env bash
set -euo pipefail

DB_PATH="/opt/livelife/data/sqlite/livelife.db"
BACKUP_DIR="/opt/livelife/backups"
STAMP="$(date +%Y%m%d-%H%M%S)"

mkdir -p "$BACKUP_DIR"

if [ -f "$DB_PATH" ]; then
  sqlite3 "$DB_PATH" ".backup '$BACKUP_DIR/livelife-$STAMP.db'"
fi

tar -czf "$BACKUP_DIR/uploads-$STAMP.tar.gz" -C /opt/livelife uploads

find "$BACKUP_DIR" -type f -mtime +14 -delete
```

设置执行权限：

```bash
chmod +x /opt/livelife/app/backup-sqlite.sh
```

手动测试：

```bash
/opt/livelife/app/backup-sqlite.sh
ls -lh /opt/livelife/backups
```

添加定时任务：

```bash
crontab -e
```

添加：

```cron
30 3 * * * /opt/livelife/app/backup-sqlite.sh >> /opt/livelife/logs/backup.log 2>&1
```

备注：

- 这个脚本每天凌晨 03:30 备份一次。
- 本机只保留 14 天备份，避免 30 GiB 磁盘被备份占满。
- 本机备份不等于真正安全。上线后建议每周下载到本地电脑，或同步到阿里云 OSS。
- 如果未来接入订单和支付，备份频率需要提高，数据库也建议迁移 PostgreSQL。

## 14. 上线后的健康检查

容器状态：

```bash
cd /opt/livelife/app
docker compose ps
```

本机检查：

```bash
curl -I http://127.0.0.1:3000
curl -I http://127.0.0.1:8080/health
```

公网检查：

```bash
curl -I https://live-life.asia
curl -I https://live-life.asia/api/health
```

Nginx 状态：

```bash
sudo systemctl status nginx
sudo nginx -t
```

磁盘和内存：

```bash
df -h
free -h
```

备注：

- `docker compose ps` 里前端和后端都应该是 running。
- `/api/health` 需要后端代码提供健康检查接口。
- 如果本机 `127.0.0.1` 正常但公网不正常，优先查 Nginx、HTTPS、防火墙、DNS。
- 如果公网根域名能打开但 `/api/health` 不行，优先查后端容器和 Nginx `/api/` 配置。

## 15. P0 执行顺序

建议按这个顺序做：

1. 注册或确认 `live-life.asia` 域名。
2. 在阿里云控制台记录服务器公网 IP。
3. 在阿里云控制台开放 `22`、`80`、`443`。
4. 添加 DNS 解析：`@` 和 `www` 都指向服务器公网 IP。
5. SSH 登录服务器。
6. 更新 Ubuntu 系统。
7. 安装 Docker、Nginx、Certbot、UFW、Fail2ban、SQLite。
8. 创建 `deploy` 用户。
9. 添加 2 GiB swap。
10. 配置 Ubuntu UFW 防火墙。
11. 创建 `/opt/livelife` 目录结构。
12. 等项目代码准备好后放入 `/opt/livelife/app`。
13. 创建 `.env` 和 `docker-compose.yml`。
14. 启动 Docker Compose。
15. 配置 Nginx HTTP 反向代理。
16. DNS 生效后申请 HTTPS 证书。
17. 配置 SQLite 和上传文件备份。
18. 检查公网 HTTPS 和 `/api/health`。

## 16. 暂时不要做的事情

P0 阶段先不要急着做：

- 不要直接上 Kubernetes，当前项目规模不需要。
- 不要把数据库端口暴露公网。
- 不要一开始就在 1 GiB 服务器上跑 PostgreSQL、Redis、对象存储模拟服务等一堆服务。
- 不要在没有备份的情况下做系统重装、磁盘初始化、数据库迁移。
- 不要在阿里云控制台里随手点“重置密码”“重置系统”“释放实例”等高风险操作。

## 17. 什么时候升级服务器或数据库

继续使用当前服务器和 SQLite，直到出现这些情况：

- 后台需要多人同时编辑内容。
- 站内购买、订单、支付成为正式生产数据。
- 报名表单数据明显增长，恢复和导出变得重要。
- 图片上传量变大，30 GiB 磁盘压力明显。
- `free -h` 经常显示内存紧张，服务频繁被杀。
- 页面访问量明显增加，后端响应变慢。

升级建议：

- 第一阶段升级：2 vCPU / 2 GiB 或 4 GiB 内存。
- 数据库升级：SQLite -> PostgreSQL。
- 文件升级：本地 uploads -> 阿里云 OSS。
- 备份升级：本机备份 -> OSS 或其他异地备份。

## 18. P0 完成标准

P0 算完成，需要满足：

- `live-life.asia` 能解析到阿里云东京服务器。
- `https://live-life.asia` 能打开。
- `https://www.live-life.asia` 能打开或跳转到主域名。
- Nginx 正常运行。
- Docker 正常运行。
- 前端容器和后端容器正常运行。
- `/api/health` 返回正常。
- 服务器只开放必要公网端口。
- SQLite 数据和上传文件有每日备份。
- 所有关键部署文件集中在 `/opt/livelife`。

完成这些以后，再进入 P1：搭建前端页面骨架、后端 API 骨架、数据库表结构，以及第一个演出情报 / CD 严选列表页。
