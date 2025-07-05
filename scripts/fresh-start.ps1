# Fresh Start Script - Simple database reset
# Use this after git pull to ensure consistent data

Write-Host "=== Fresh Start ===" -ForegroundColor Green

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

Write-Host "5. Testing APIs..." -ForegroundColor Yellow
try {
    $customers = Invoke-WebRequest -Uri "http://localhost:8082/customers" -UseBasicParsing -TimeoutSec 10
    $products = Invoke-WebRequest -Uri "http://localhost:8081/products" -UseBasicParsing -TimeoutSec 10
    
    if ($customers.StatusCode -eq 200 -and $products.StatusCode -eq 200) {
        Write-Host "   ✅ APIs working correctly" -ForegroundColor Green
    }
} catch {
    Write-Host "   ⚠️  APIs still starting up (normal)" -ForegroundColor Yellow
}

Set-Location $currentDir

Write-Host "`n✅ Fresh start complete!" -ForegroundColor Green
Write-Host "🌐 Frontend: http://localhost:8080" -ForegroundColor Cyan