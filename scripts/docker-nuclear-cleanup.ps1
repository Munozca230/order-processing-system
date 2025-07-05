# Nuclear Docker Cleanup - ALL DOCKER RESOURCES
# ⚠️ WARNING: This affects ALL Docker containers/images/networks on your system

Write-Host "=== ☢️  NUCLEAR DOCKER CLEANUP ☢️  ===" -ForegroundColor Red
Write-Host "⚠️  WARNING: This will remove ALL Docker resources on your system!" -ForegroundColor Red
Write-Host "📋 This includes:" -ForegroundColor Yellow
Write-Host "   • All containers (from all projects)" -ForegroundColor Yellow
Write-Host "   • All unused images" -ForegroundColor Yellow
Write-Host "   • All unused networks" -ForegroundColor Yellow
Write-Host "   • All unused volumes" -ForegroundColor Yellow
Write-Host ""

$confirmation1 = Read-Host "Are you absolutely sure? Type 'NUCLEAR' to confirm"
if ($confirmation1 -ne 'NUCLEAR') {
    Write-Host "Cleanup cancelled - Smart choice!" -ForegroundColor Green
    exit
}

$confirmation2 = Read-Host "Last chance - this will affect OTHER Docker projects too. Type 'YES' to proceed"
if ($confirmation2 -ne 'YES') {
    Write-Host "Cleanup cancelled" -ForegroundColor Green
    exit
}

Write-Host "🚨 Starting nuclear cleanup..." -ForegroundColor Red

Write-Host "1. Stopping ALL Docker containers..." -ForegroundColor Yellow
docker stop $(docker ps -aq) 2>$null

Write-Host "2. Removing ALL containers..." -ForegroundColor Yellow
docker rm $(docker ps -aq) 2>$null

Write-Host "3. Cleaning ALL Docker system..." -ForegroundColor Yellow
docker system prune -f --volumes

Write-Host "4. Removing ALL networks..." -ForegroundColor Yellow
docker network prune -f

Write-Host "5. Removing ALL unused images..." -ForegroundColor Yellow
docker image prune -a -f

Write-Host "`n☢️  Nuclear cleanup complete!" -ForegroundColor Red
Write-Host "💡 ALL Docker resources have been removed" -ForegroundColor Yellow
Write-Host "🚀 You can now start fresh with any Docker project" -ForegroundColor Green