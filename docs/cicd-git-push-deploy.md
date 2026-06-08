# LIVE LIFE Git Push Deployment

This P0 deployment flow lets VSCode deploy to the Alibaba Cloud Tokyo server by pushing Git.

```text
VSCode / local Git
  -> git push aliyun master
  -> Alibaba Cloud server bare Git repo over SSH port 2222
  -> post-receive hook
  -> /opt/livelife/app
  -> background deploy script
  -> docker compose up -d --build
  -> Nginx serves live-life.asia
```

## Why This Flow

- It matches the current need: push locally and sync to Alibaba Cloud.
- It does not require GitHub Actions, GitLab CI, or another SaaS yet.
- The server already has Docker, Docker Compose, Nginx, and `/opt/livelife`.
- It keeps runtime data outside Git: `/opt/livelife/data`, `/opt/livelife/uploads`, `/opt/livelife/logs`.

## One-Time Server Setup

Upload `scripts/setup-aliyun-git-deploy.sh` to the server, for example as `/home/admin/setup-aliyun-git-deploy.sh`, then run:

```bash
bash ~/setup-aliyun-git-deploy.sh master
```

Also open a dedicated SSH port for Git deploy. Keep port `22` for Workbench/normal SSH and use `2222` for local Git push:

```bash
echo 'Port 22' | sudo tee /etc/ssh/sshd_config.d/99-livelife-git-deploy.conf
echo 'Port 2222' | sudo tee -a /etc/ssh/sshd_config.d/99-livelife-git-deploy.conf
sudo systemctl reload ssh
sudo ufw allow 2222/tcp
sudo ufw reload
sudo ss -lntp | grep -E ':22|:2222'
```

In the Alibaba Cloud instance firewall, add:

```text
TCP 2222 0.0.0.0/0
```

If the script warns that the current user cannot run Docker without sudo, run:

```bash
sudo usermod -aG docker admin
```

Then log out of Workbench/SSH and reconnect before deploying.

## One-Time Local Setup

From the local project folder:

```powershell
git remote add aliyun ssh://admin@47.74.8.10:2222/opt/livelife/git/livelife.git
```

If `admin` is not the server login user, replace it.

## Daily Deploy

In VSCode terminal:

```powershell
git status
git add .
git commit -m "Update LIVE LIFE"
git push aliyun master
```

The push output should show Docker Compose building and starting containers.
The push output should return quickly and say that deployment started in the background. Watch progress on the server:

```bash
tail -f /opt/livelife/logs/deploy.log
```

## Verify After Deploy

On the server:

```bash
cd /opt/livelife/app
docker compose ps
curl -I http://127.0.0.1:3000
curl -I http://127.0.0.1:8080/api/health
```

From a browser:

```text
https://live-life.asia
https://live-life.asia/api/health
```

## When To Upgrade To GitHub Actions

Move to GitHub Actions later when the project needs pull requests, remote CI tests, deployment approval, or a team workflow. The GitHub Actions version would push to GitHub first, then SSH into Alibaba Cloud and run the same deploy commands.
