@echo off
setlocal

rem LIVE LIFE local backend launcher for Windows cmd.
rem Keep cache and runtime paths inside this repository.

set "SCRIPT_DIR=%~dp0"
for %%I in ("%SCRIPT_DIR%..") do set "REPO_ROOT=%%~fI"
set "GOCACHE=%REPO_ROOT%\.cache\go-build"
set "LOG_DIR=%REPO_ROOT%\logs"
set "GO_EXE=C:\Program Files\Go\bin\go.exe"

if not exist "%GOCACHE%" mkdir "%GOCACHE%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

cd /d "%REPO_ROOT%\backend"
"%GO_EXE%" run ./cmd/server
