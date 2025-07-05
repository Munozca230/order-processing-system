# Simple Reset Script
# Use this when you need fresh data

Write-Host "=== Simple Reset ===" -ForegroundColor Green

$currentDir = Get-Location
Set-Location "../infra"

Write-Host "1. Stopping services and cleaning up..." -ForegroundColor Yellow
docker-compose down -v --remove-orphans

Write-Host "2. Force cleanup if needed..." -ForegroundColor Yellow
docker container prune -f 2>$null
docker network prune -f 2>$null

Write-Host "3. Starting with fresh data..." -ForegroundColor Yellow
docker-compose up -d

Write-Host "4. Waiting for core services to be healthy..." -ForegroundColor Yellow
Start-Sleep -Seconds 20  # Reduced since healthchecks are faster now

Set-Location $currentDir

Write-Host "`n✅ Reset complete!" -ForegroundColor Green
Write-Host "🌐 Frontend: http://localhost:8080" -ForegroundColor Cyan