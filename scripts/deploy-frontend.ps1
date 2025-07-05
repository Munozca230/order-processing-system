#!/usr/bin/env pwsh

param(
    [switch]$PreserveData,
    [switch]$Help
)

if ($Help) {
    Write-Host "Frontend Profile Deployment - Usage:" -ForegroundColor Cyan
    Write-Host "  .\deploy-frontend.ps1           # Clean deployment (default - recommended)"
    Write-Host "  .\deploy-frontend.ps1 -PreserveData  # Preserve existing data volumes"
    Write-Host "  .\deploy-frontend.ps1 -Help          # Show this help"
    exit 0
}

Write-Host "Frontend Profile Deployment (Full Stack + Web Interface)" -ForegroundColor Cyan
Write-Host "=========================================================" -ForegroundColor Cyan

if ($PreserveData) {
    Write-Host "MODE: Preserving existing data volumes" -ForegroundColor Yellow
} else {
    Write-Host "MODE: Clean deployment (volumes will be reset)" -ForegroundColor Green
}

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

# Stop any existing services
if ($PreserveData) {
    Write-Host "`nStopping existing services (preserving data)..." -ForegroundColor Yellow
    try {
        docker-compose --profile frontend --profile backend down
        Write-Host "Existing services stopped (data preserved)" -ForegroundColor Green
    } catch {
        Write-Host "No existing services to stop" -ForegroundColor Yellow
    }
} else {
    Write-Host "`nStopping existing services and cleaning volumes..." -ForegroundColor Yellow
    try {
        docker-compose --profile frontend --profile backend down -v
        Write-Host "Existing services stopped and volumes cleaned" -ForegroundColor Green
    } catch {
        Write-Host "No existing services to stop" -ForegroundColor Yellow
    }

    # Clean any orphaned volumes
    Write-Host "Cleaning orphaned Docker volumes..." -ForegroundColor Yellow
    try {
        docker volume prune -f
        Write-Host "Orphaned volumes cleaned" -ForegroundColor Green
    } catch {
        Write-Host "No orphaned volumes to clean" -ForegroundColor Yellow
    }
}

# Start frontend profile (includes all backend + frontend services)
Write-Host "`nStarting Frontend Profile services..." -ForegroundColor Yellow
Write-Host "Backend Services:" -ForegroundColor Cyan
Write-Host "  - Zookeeper (Kafka coordination)" -ForegroundColor Gray
Write-Host "  - Kafka (Message broker)" -ForegroundColor Gray
Write-Host "  - MongoDB (Database + sample data)" -ForegroundColor Gray
Write-Host "  - Redis (Cache & distributed locks)" -ForegroundColor Gray
Write-Host "  - Order Worker (Java - Core processing)" -ForegroundColor Gray
Write-Host "  - Product API (Go - Product catalog)" -ForegroundColor Gray
Write-Host "  - Customer API (Go - Customer management)" -ForegroundColor Gray
Write-Host "`nFrontend Services:" -ForegroundColor Green
Write-Host "  - Order API (Node.js - Frontend bridge to Kafka)" -ForegroundColor Gray
Write-Host "  - Nginx Frontend (Web server + proxy)" -ForegroundColor Gray

try {
    docker-compose --profile frontend up -d
    Write-Host "`nFrontend profile services started successfully!" -ForegroundColor Green
} catch {
    Write-Host "`nFailed to start frontend services" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Wait for services to be ready
Write-Host "`nWaiting for services to be healthy (60 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 60

# Check service status
Write-Host "`nChecking service status..." -ForegroundColor Cyan
try {
    $status = docker-compose ps
    Write-Host $status -ForegroundColor Gray
} catch {
    Write-Host "Could not retrieve service status" -ForegroundColor Yellow
}

# Verify critical services including frontend
Write-Host "`nHealth checks..." -ForegroundColor Cyan

$healthChecks = @(
    @{Name="Frontend Web"; Url="http://localhost:8080/health"},
    @{Name="Order API"; Url="http://localhost:3000/health"},
    @{Name="Product API"; Url="http://localhost:8081/health"},
    @{Name="Customer API"; Url="http://localhost:8082/health"}
)

foreach ($check in $healthChecks) {
    try {
        $response = Invoke-RestMethod -Uri $check.Url -TimeoutSec 5
        Write-Host "  $($check.Name) is healthy" -ForegroundColor Green
    } catch {
        Write-Host "  $($check.Name) is not responding" -ForegroundColor Red
    }
}

# Display useful information
Write-Host "`nFrontend Profile Ready!" -ForegroundColor Green
Write-Host "========================" -ForegroundColor Green

Write-Host "`nPrimary Access Point:" -ForegroundColor Cyan
Write-Host "  Frontend Web Interface: http://localhost:8080" -ForegroundColor Yellow
Write-Host "     Complete visual interface for order testing" -ForegroundColor Gray

Write-Host "`nAPI Endpoints:" -ForegroundColor Cyan
Write-Host "  Order API: http://localhost:3000" -ForegroundColor Gray
Write-Host "  Product API: http://localhost:8081" -ForegroundColor Gray
Write-Host "  Customer API: http://localhost:8082" -ForegroundColor Gray

Write-Host "`nDatabase Access:" -ForegroundColor Cyan
Write-Host "  MongoDB: mongodb://localhost:27017" -ForegroundColor Gray
Write-Host "  Redis: redis://localhost:6379" -ForegroundColor Gray
Write-Host "  Kafka: kafka://localhost:9092" -ForegroundColor Gray

Write-Host "`nTesting Options:" -ForegroundColor Cyan
Write-Host "  Web Interface: Open http://localhost:8080 in browser" -ForegroundColor Green
Write-Host "  Postman Collection: Import postman/*.json files" -ForegroundColor Gray
Write-Host "  Automated Tests: scripts/test-order-structure.ps1" -ForegroundColor Gray
Write-Host "  Direct API Test:" -ForegroundColor Gray
Write-Host "    curl -X POST http://localhost:3000/api/orders -H 'Content-Type: application/json' -d '{`"orderId`":`"test`",`"customerId`":`"customer-1`",`"products`":[{`"productId`":`"product-1`"}]}'" -ForegroundColor DarkGray

Write-Host "`nMonitoring:" -ForegroundColor Cyan
Write-Host "  Service Status: docker-compose ps" -ForegroundColor Gray
Write-Host "  Order Worker Logs: docker-compose logs -f order-worker" -ForegroundColor Gray
Write-Host "  Frontend Logs: docker-compose logs -f nginx-frontend" -ForegroundColor Gray
Write-Host "  Order API Logs: docker-compose logs -f order-api" -ForegroundColor Gray
Write-Host "  MongoDB Data: docker-compose exec mongo mongosh orders --eval `"db.orders.find().forEach(printjson)`"" -ForegroundColor Gray

Write-Host "`nQuick Start Guide:" -ForegroundColor Yellow
Write-Host "  1. Open http://localhost:8080 in your browser" -ForegroundColor Gray
Write-Host "  2. Click 'Verificar Estado' to check all services" -ForegroundColor Gray
Write-Host "  3. Fill out the order form and click 'Enviar Orden'" -ForegroundColor Gray
Write-Host "  4. Watch real-time processing in the interface" -ForegroundColor Gray
Write-Host "  5. Verify results in MongoDB using the provided commands" -ForegroundColor Gray

Write-Host "`nTo switch to Backend only:" -ForegroundColor Yellow
Write-Host "  docker-compose down && docker-compose up -d" -ForegroundColor Gray

Write-Host "`nTo stop all services:" -ForegroundColor Yellow
Write-Host "  docker-compose --profile frontend down" -ForegroundColor Gray

# Try to open browser automatically
Write-Host "`nAttempting to open web interface..." -ForegroundColor Cyan
try {
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        Start-Process "http://localhost:8080"
    } elseif ($IsMacOS) {
        & open "http://localhost:8080"
    } elseif ($IsLinux) {
        & xdg-open "http://localhost:8080"
    }
    Write-Host "Web interface should open in your default browser" -ForegroundColor Green
} catch {
    Write-Host "Could not auto-open browser. Please manually visit: http://localhost:8080" -ForegroundColor Yellow
}

Write-Host "`nFrontend Profile deployment completed!" -ForegroundColor Green
Write-Host "Access your order processing system at: http://localhost:8080" -ForegroundColor Yellow