#!/usr/bin/env pwsh

Write-Host "Clean System Restart - Complete Docker Environment Reset" -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan

# Navigate to infra directory
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$infraPath = Join-Path $projectRoot "infra"

if (-not (Test-Path $infraPath)) {
    Write-Host "ERROR: Infrastructure directory not found: $infraPath" -ForegroundColor Red
    exit 1
}

Set-Location $infraPath
Write-Host "Working directory: $(Get-Location)" -ForegroundColor Gray

Write-Host "`nThis will completely reset the Docker environment:" -ForegroundColor Yellow
Write-Host "  - Stop all containers" -ForegroundColor Gray
Write-Host "  - Remove all volumes (MongoDB data, Kafka data, Redis data)" -ForegroundColor Gray
Write-Host "  - Clean orphaned volumes" -ForegroundColor Gray
Write-Host "  - Remove unused Docker images" -ForegroundColor Gray

$confirmation = Read-Host "`nContinue? (y/N)"
if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
    Write-Host "Operation cancelled" -ForegroundColor Yellow
    exit 0
}

Write-Host "`nStep 1: Stopping all services and removing volumes..." -ForegroundColor Yellow
try {
    docker-compose --profile frontend --profile backend down -v
    Write-Host "Services stopped and volumes removed" -ForegroundColor Green
} catch {
    Write-Host "No services were running" -ForegroundColor Yellow
}

Write-Host "`nStep 2: Cleaning orphaned volumes..." -ForegroundColor Yellow
try {
    docker volume prune -f
    Write-Host "Orphaned volumes cleaned" -ForegroundColor Green
} catch {
    Write-Host "No orphaned volumes found" -ForegroundColor Yellow
}

Write-Host "`nStep 3: Cleaning unused Docker images..." -ForegroundColor Yellow
try {
    docker image prune -f
    Write-Host "Unused images cleaned" -ForegroundColor Green
} catch {
    Write-Host "No unused images found" -ForegroundColor Yellow
}

Write-Host "`nStep 4: Cleaning Docker build cache..." -ForegroundColor Yellow
try {
    docker builder prune -f
    Write-Host "Build cache cleaned" -ForegroundColor Green
} catch {
    Write-Host "No build cache to clean" -ForegroundColor Yellow
}

Write-Host "`nDocker environment cleaned successfully!" -ForegroundColor Green
Write-Host "===============================================" -ForegroundColor Green

Write-Host "`nNext steps:" -ForegroundColor Cyan
Write-Host "  Backend only: .\deploy-backend.ps1" -ForegroundColor Gray
Write-Host "  Frontend full: .\deploy-frontend.ps1" -ForegroundColor Gray
Write-Host "  Manual: cd ..\infra && docker-compose up -d" -ForegroundColor Gray

Write-Host "`nSystem is ready for a fresh deployment!" -ForegroundColor Green