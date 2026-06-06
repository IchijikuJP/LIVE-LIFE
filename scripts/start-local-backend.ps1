$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$env:GOCACHE = Join-Path $repoRoot ".cache\go-build"
$logDir = Join-Path $repoRoot "logs"
$logFile = Join-Path $logDir "local-backend.out.log"
$goExe = "C:\Program Files\Go\bin\go.exe"

New-Item -ItemType Directory -Force -Path $env:GOCACHE | Out-Null
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
Set-Location (Join-Path $repoRoot "backend")

"[$(Get-Date -Format o)] starting LiveLife backend" | Out-File -FilePath $logFile -Append -Encoding utf8
& $goExe run ./cmd/server 2>&1 | Tee-Object -FilePath $logFile -Append
