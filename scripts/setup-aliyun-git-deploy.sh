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
DEPLOY_LOG="$BASE_DIR/logs/deploy.log"

mkdir -p "$BASE_DIR/git" "$APP_DIR" "$DEPLOY_DIR" "$BASE_DIR/data/sqlite" "$BASE_DIR/uploads" "$BASE_DIR/logs" "$BASE_DIR/backups"

if [ ! -d "$REPO_DIR" ]; then
  git init --bare "$REPO_DIR"
fi

cat > "$DEPLOY_SCRIPT" <<DEPLOY
#!/usr/bin/env bash
set -euo pipefail

APP_DIR="$APP_DIR"
DEPLOY_LOG="$DEPLOY_LOG"
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
docker compose up -d --build
docker compose ps
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

if ! docker info >/dev/null 2>&1; then
  cat <<WARN

WARNING: user '$DEPLOY_USER' cannot run docker without sudo yet.
Run this once, then log out and reconnect:

  sudo usermod -aG docker $DEPLOY_USER

The git push deploy hook needs plain 'docker compose' to work non-interactively.
WARN
fi

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

DONE
