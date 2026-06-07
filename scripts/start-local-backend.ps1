$ErrorActionPreference = "Stop"

# 本脚本用于启动 LIVE LIFE 本地 Go 后端。
# 使用方式：
#   powershell -ExecutionPolicy Bypass -File .\scripts\start-local-backend.ps1
#
# 为什么不直接写 go run：
# 1. 统一把 Go 编译缓存放到项目内 .cache，避免权限或系统目录问题。
# 2. 统一把本地运行日志写入 logs/local-backend.out.log，方便排查。
# 3. 自动切到 backend 目录，因为 Go 服务会从 backend/static 读取本地预览页。
$repoRoot = Split-Path -Parent $PSScriptRoot
$env:GOCACHE = Join-Path $repoRoot ".cache\go-build"
$logDir = Join-Path $repoRoot "logs"
$logFile = Join-Path $logDir "local-backend.out.log"
$goExe = "C:\Program Files\Go\bin\go.exe"

New-Item -ItemType Directory -Force -Path $env:GOCACHE | Out-Null
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
Set-Location (Join-Path $repoRoot "backend")

"[$(Get-Date -Format o)] starting LIVE LIFE backend" | Out-File -FilePath $logFile -Append -Encoding utf8
& $goExe run ./cmd/server 2>&1 | Tee-Object -FilePath $logFile -Append
