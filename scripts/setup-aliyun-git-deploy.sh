#!/usr/bin/env bash
set -euo pipefail

DEPLOY_BRANCH="${1:-master}"
DEPLOY_USER="$(whoami)"
SSH_PORT="${SSH_PORT:-2222}"
BASE_DIR="/opt/livelife"
APP_DIR="$BASE_DIR/app"
REPO_DIR="$BASE_DIR/git/livelife.git"
DEPLOY_DIR="$BASE_DIR/deploy"
DEPLOY_SCRIPT="$DEPLOY_DIR/deploy.sh"
START_BACKEND_SCRIPT="$DEPLOY_DIR/start-backend.sh"
DEPLOY_LOG="$BASE_DIR/logs/deploy.log"
BACKEND_LOG="$BASE_DIR/logs/backend.log"
RUNTIME_DIR="$BASE_DIR/runtime"
WWW_DIR="$BASE_DIR/www"

mkdir -p "$BASE_DIR/git" "$APP_DIR" "$DEPLOY_DIR" "$RUNTIME_DIR" "$WWW_DIR" "$BASE_DIR/data/sqlite" "$BASE_DIR/uploads" "$BASE_DIR/logs" "$BASE_DIR/backups"

if [ ! -d "$REPO_DIR" ]; then
  git init --bare "$REPO_DIR"
fi

cat > "$START_BACKEND_SCRIPT" <<START
#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$BASE_DIR"
RUNTIME_DIR="$RUNTIME_DIR"
BACKEND_BIN="\$RUNTIME_DIR/livelife-api"
PID_FILE="\$RUNTIME_DIR/livelife-api.pid"
BACKEND_LOG="$BACKEND_LOG"

mkdir -p "\$RUNTIME_DIR" "\$BASE_DIR/data/sqlite" "\$BASE_DIR/uploads" "\$BASE_DIR/logs"

if [ ! -x "\$BACKEND_BIN" ]; then
  echo "[\$(date -Is)] backend binary missing: \$BACKEND_BIN" >> "\$BACKEND_LOG"
  exit 1
fi

if [ -f "\$PID_FILE" ]; then
  OLD_PID="\$(cat "\$PID_FILE" || true)"
  if [ -n "\$OLD_PID" ] && kill -0 "\$OLD_PID" 2>/dev/null; then
    kill "\$OLD_PID" || true
    for _ in 1 2 3 4 5; do
      if ! kill -0 "\$OLD_PID" 2>/dev/null; then
        break
      fi
      sleep 1
    done
  fi
fi

nohup env \\
  PORT=8080 \\
  BACKEND_PORT=8080 \\
  APP_ENV=production \\
  DATABASE_PATH="\$BASE_DIR/data/sqlite/livelife.db" \\
  "\$BACKEND_BIN" >> "\$BACKEND_LOG" 2>&1 &

echo "\$!" > "\$PID_FILE"
echo "[\$(date -Is)] backend started pid=\$(cat "\$PID_FILE")" >> "\$BACKEND_LOG"
START

chmod +x "$START_BACKEND_SCRIPT"

cat > "$DEPLOY_SCRIPT" <<DEPLOY
#!/usr/bin/env bash
set -euo pipefail

APP_DIR="$APP_DIR"
DEPLOY_LOG="$DEPLOY_LOG"
RUNTIME_DIR="$RUNTIME_DIR"
WWW_DIR="$WWW_DIR"
START_BACKEND_SCRIPT="$START_BACKEND_SCRIPT"
LOCK_DIR="/tmp/livelife-deploy.lock"

mkdir -p "\$(dirname "\$DEPLOY_LOG")"
exec >> "\$DEPLOY_LOG" 2>&1

if ! mkdir "\$LOCK_DIR" 2>/dev/null; then
  echo "[\$(date -Is)] deploy skipped: another deployment is already running"
  exit 0
fi
trap 'rmdir "\$LOCK_DIR"' EXIT

echo "[\$(date -Is)] deploy start"
cd "\$APP_DIR"
if [ ! -f "\$APP_DIR/deploy/release/backend/livelife-api" ]; then
  echo "missing release backend: \$APP_DIR/deploy/release/backend/livelife-api"
  exit 1
fi
if [ ! -f "\$APP_DIR/deploy/release/frontend/index.html" ]; then
  echo "missing release frontend: \$APP_DIR/deploy/release/frontend/index.html"
  exit 1
fi

install -m 0755 "\$APP_DIR/deploy/release/backend/livelife-api" "\$RUNTIME_DIR/livelife-api"
rm -rf "\$WWW_DIR"
mkdir -p "\$WWW_DIR"
cp -a "\$APP_DIR/deploy/release/frontend/." "\$WWW_DIR/"
"\$START_BACKEND_SCRIPT"
echo "[\$(date -Is)] deploy done"
DEPLOY

chmod +x "$DEPLOY_SCRIPT"

cat > "$REPO_DIR/hooks/post-receive" <<HOOK
#!/usr/bin/env bash
set -euo pipefail

DEPLOY_BRANCH="$DEPLOY_BRANCH"
APP_DIR="$APP_DIR"
REPO_DIR="$REPO_DIR"
DEPLOY_SCRIPT="$DEPLOY_SCRIPT"
DEPLOY_LOG="$DEPLOY_LOG"
SHOULD_DEPLOY=0

while read -r oldrev newrev refname; do
  if [ "\$refname" = "refs/heads/\$DEPLOY_BRANCH" ]; then
    SHOULD_DEPLOY=1
  fi
done

if [ "\$SHOULD_DEPLOY" != "1" ]; then
  echo "No deployment branch update. Expected refs/heads/\$DEPLOY_BRANCH."
  exit 0
fi

echo "Deploying LIVE LIFE from \$DEPLOY_BRANCH to \$APP_DIR"
mkdir -p "\$APP_DIR"
GIT_DIR="\$REPO_DIR" GIT_WORK_TREE="\$APP_DIR" git checkout -f "\$DEPLOY_BRANCH"

nohup "\$DEPLOY_SCRIPT" >/dev/null 2>&1 &
echo "Deployment started in background."
echo "Watch logs on the server with:"
echo "  tail -f \$DEPLOY_LOG"
HOOK

chmod +x "$REPO_DIR/hooks/post-receive"

cat > /tmp/livelife-nginx.conf <<NGINX
server {
    listen 80;
    server_name live-life.asia www.live-life.asia;

    root $WWW_DIR;
    index index.html;
    client_max_body_size 20m;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        try_files \$uri \$uri/ /index.html;
    }
}
NGINX

if command -v sudo >/dev/null 2>&1; then
  sudo cp /tmp/livelife-nginx.conf /etc/nginx/sites-available/livelife
  sudo ln -sf /etc/nginx/sites-available/livelife /etc/nginx/sites-enabled/livelife
  sudo rm -f /etc/nginx/sites-enabled/default
  sudo nginx -t
  sudo systemctl reload nginx
else
  cp /tmp/livelife-nginx.conf /etc/nginx/sites-available/livelife
  ln -sf /etc/nginx/sites-available/livelife /etc/nginx/sites-enabled/livelife
  rm -f /etc/nginx/sites-enabled/default
  nginx -t
fi

(crontab -l 2>/dev/null | grep -v "$START_BACKEND_SCRIPT" || true; echo "@reboot $START_BACKEND_SCRIPT") | crontab -

cat <<DONE

LIVE LIFE git deployment is ready.

Server bare repo:
  $REPO_DIR

Deploy branch:
  $DEPLOY_BRANCH

From your local project, add the server remote:
  git remote add aliyun ssh://$DEPLOY_USER@47.74.8.10:$SSH_PORT$REPO_DIR

Deploy:
  git push aliyun $DEPLOY_BRANCH

This no-build deployment expects local release artifacts in:
  deploy/release/backend/livelife-api
  deploy/release/frontend/index.html

DONE
