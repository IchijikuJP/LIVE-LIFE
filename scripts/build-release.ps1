$ErrorActionPreference = "Stop"

$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
$FrontendDir = Join-Path $Root "frontend"
$BackendDir = Join-Path $Root "backend"
$ReleaseDir = Join-Path $Root "deploy\release"
$BackendReleaseDir = Join-Path $ReleaseDir "backend"
$FrontendReleaseDir = Join-Path $ReleaseDir "frontend"

New-Item -ItemType Directory -Force -Path $BackendReleaseDir | Out-Null
New-Item -ItemType Directory -Force -Path $FrontendReleaseDir | Out-Null

$NodeCandidates = @()
$SystemNode = (Get-Command node -ErrorAction SilentlyContinue).Source
if ($SystemNode) {
  $NodeCandidates += $SystemNode
}
$BundledNode = Join-Path $env:USERPROFILE ".cache\codex-runtimes\codex-primary-runtime\dependencies\node\bin\node.exe"
if (Test-Path $BundledNode) {
  $NodeCandidates += $BundledNode
}

$Node = $null
foreach ($Candidate in $NodeCandidates) {
  try {
    & $Candidate --version | Out-Null
    $Node = $Candidate
    break
  }
  catch {
    $Node = $null
  }
}
if (-not $Node) {
  throw "Node.js was not found. Install Node.js or run from Codex desktop with bundled Node."
}

Push-Location $FrontendDir
try {
  & $Node "node_modules\typescript\bin\tsc" -b
  & $Node "node_modules\vite\bin\vite.js" build --config ".\vite.config.ts"
}
finally {
  Pop-Location
}

if (Test-Path $FrontendReleaseDir) {
  Remove-Item -Recurse -Force $FrontendReleaseDir
}
New-Item -ItemType Directory -Force -Path $FrontendReleaseDir | Out-Null
robocopy (Join-Path $FrontendDir "dist") $FrontendReleaseDir /MIR | Out-Host
if ($LASTEXITCODE -gt 7) {
  throw "robocopy failed with exit code $LASTEXITCODE"
}

Push-Location $BackendDir
try {
  $env:GOCACHE = Join-Path $Root ".cache\go-build"
  $env:GOOS = "linux"
  $env:GOARCH = "amd64"
  $env:CGO_ENABLED = "0"
  go build -o (Join-Path $BackendReleaseDir "livelife-api") ./cmd/server
}
finally {
  Pop-Location
  Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
  Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
  Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue
}

Write-Host "Release built:"
Get-ChildItem -Recurse $ReleaseDir | Select-Object FullName, Length
