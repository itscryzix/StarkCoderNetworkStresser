# Windows (PowerShell) build script. Replaces `make all` on machines without make.
# Usage: powershell -ExecutionPolicy Bypass -File build.ps1
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $Root

$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

Write-Host "==> go mod tidy" -ForegroundColor Cyan
go mod tidy

Write-Host "==> building web client" -ForegroundColor Cyan
Push-Location web-client
npm install --no-audit --no-fund
npm run build
Pop-Location

Write-Host "==> copying web assets to bin/web-client" -ForegroundColor Cyan
$BinWeb = Join-Path $Root "bin/web-client"
if (Test-Path $BinWeb) { Remove-Item -Recurse -Force $BinWeb }
New-Item -ItemType Directory -Force -Path $BinWeb | Out-Null
Copy-Item -Recurse -Force (Join-Path $Root "web-client/dist/public/*") $BinWeb

Write-Host "==> building mmb-cli and mmb-server" -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path (Join-Path $Root "bin") | Out-Null
go build -o (Join-Path $Root "bin/mmb-cli.exe") ./cmd/mmb-cli
go build -o (Join-Path $Root "bin/mmb-server.exe") ./cmd/mmb-server

Write-Host "Build complete: bin/mmb-server.exe, bin/mmb-cli.exe" -ForegroundColor Green
