# Windows (PowerShell) run script. Starts the server with no-proxy allowed by default.
# Usage: powershell -ExecutionPolicy Bypass -File run-server.ps1
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $Root
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

$Server = Join-Path $Root "bin/mmb-server.exe"
if (!(Test-Path $Server)) {
  Write-Host "bin/mmb-server.exe not found, building first..." -ForegroundColor Yellow
  powershell -ExecutionPolicy Bypass -File (Join-Path $Root "build.ps1")
}

if (-not $env:ALLOW_NO_PROXY) { $env:ALLOW_NO_PROXY = "true" }
Write-Host "Starting Stark Coder on http://localhost:3000 (ALLOW_NO_PROXY=$env:ALLOW_NO_PROXY)" -ForegroundColor Cyan
& $Server
