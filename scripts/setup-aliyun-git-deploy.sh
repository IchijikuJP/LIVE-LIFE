#!/usr/bin/env bash
set -euo pipefail

DEPLOY_BRANCH="${1:-master}"
DEPLOY_USER="$(whoami)"
BASE_DIR="/opt/livelife"
APP_DIR="$BASE_DIR/app"
REPO_DIR="$BASE_DIR/git/livelife.git"

mkdir -p "$BASE_DIR/git" "$APP_DIR" "$BASE_DIR/data/sqlite" "$BASE_DIR/uploads" "$BASE_DIR/logs" "$BASE_DIR/backups"

if [ ! -d "$REPO_DIR" ]; then
  git init --bare "$REPO_DIR"
fi

cat > "$REPO_DIR/hooks/post-receive" <<HOOK
#!/usr/bin/env bash
set -euo pipefail

DEPLOY_BRANCH="$DEPLOY_BRANCH"
APP_DIR="$APP_DIR"
REPO_DIR="$REPO_DIR"
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

cd "\$APP_DIR"
docker compose up -d --build
docker compose ps
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
  git remote add aliyun $DEPLOY_USER@47.74.8.10:$REPO_DIR

Deploy:
  git push aliyun $DEPLOY_BRANCH

DONE
