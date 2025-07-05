# Docker Cleanup Script - PROJECT SPECIFIC
# This cleans up ONLY this project's resources safely

Write-Host "=== Project Docker Cleanup ===" -ForegroundColor Green
Write-Host "🎯 This will only clean resources from this project" -ForegroundColor Cyan

$currentDir = Get-Location
Set-Location "../infra"

Write-Host "1. Stopping project services..." -ForegroundColor Yellow
docker-compose down -v --remove-orphans

Write-Host "2. Removing project containers (if any stuck)..." -ForegroundColor Yellow
docker ps -a --filter "name=infra-" --format "{{.ID}}" | ForEach-Object { 
    if ($_) { docker rm -f $_ }
}

Write-Host "3. Removing project network (if stuck)..." -ForegroundColor Yellow
docker network rm infra_default 2>$null

Write-Host "4. Removing project volumes..." -ForegroundColor Yellow
docker volume rm infra_mongo-data infra_kafka-data infra_redis-data 2>$null

Write-Host "5. Removing project images (to force rebuild)..." -ForegroundColor Yellow
docker images --filter "reference=infra-*" --format "{{.ID}}" | ForEach-Object { 
    if ($_) { docker rmi -f $_ }
}

Set-Location $currentDir

Write-Host "`n✅ Project cleanup complete!" -ForegroundColor Green
Write-Host "💡 Only this project's resources were affected" -ForegroundColor Cyan
Write-Host "🚀 Now you can run fresh-start.ps1 safely" -ForegroundColor Cyan