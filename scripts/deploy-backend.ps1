#!/usr/bin/env pwsh
#
# 🚀 ALTERNATIVA CROSS-PLATFORM: make backend
# Este script mantiene funcionalidad completa, pero para simplicidad usa: make backend

param(
    [switch]$PreserveData,
    [switch]$Help
)

if ($Help) {
    Write-Host "Backend Profile Deployment - Usage:" -ForegroundColor Cyan
    Write-Host "  .\deploy-backend.ps1           # Clean deployment (default - recommended)"
    Write-Host "  .\deploy-backend.ps1 -PreserveData  # Preserve existing data volumes"
    Write-Host "  .\deploy-backend.ps1 -Help          # Show this help"
    exit 0
}

Write-Host "Backend Profile Deployment (APIs + Core Services Only)" -ForegroundColor Cyan
Write-Host "=======================================================" -ForegroundColor Cyan

if ($PreserveData) {
    Write-Host "MODE: Preserving existing data volumes" -ForegroundColor Yellow
} else {
    Write-Host "MODE: Clean deployment (volumes will be reset)" -ForegroundColor Green
}

# Check if Docker is running
Write-Host "Checking Docker status..." -ForegroundColor Yellow
try {
    $dockerVersion = docker version --format json 2>$null | ConvertFrom-Json
    if (-not $dockerVersion.Server) {
        throw "Docker daemon not responding"
    }
    Write-Host "✅ Docker is running (Client: $($dockerVersion.Client.Version), Server: $($dockerVersion.Server.Version))" -ForegroundColor Green
} catch {
    Write-Host "❌ ERROR: Docker is not running or not accessible" -ForegroundColor Red
    Write-Host "" -ForegroundColor Yellow
    Write-Host "Please ensure Docker Desktop is running:" -ForegroundColor Yellow
    Write-Host "  1. Start Docker Desktop application" -ForegroundColor Gray
    Write-Host "  2. Wait for the green 'Engine running' status" -ForegroundColor Gray
    Write-Host "  3. Verify with: docker --version" -ForegroundColor Gray
    Write-Host "" -ForegroundColor Yellow
    Write-Host "If Docker Desktop is not installed:" -ForegroundColor Yellow
    Write-Host "  Download from: https://www.docker.com/products/docker-desktop" -ForegroundColor Gray
    exit 1
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

# Start backend profile (default)
Write-Host "`nStarting Backend Profile services..." -ForegroundColor Yellow
Write-Host "Services included:" -ForegroundColor Cyan
Write-Host "  - Zookeeper (Kafka coordination)" -ForegroundColor Gray
Write-Host "  - Kafka (Message broker)" -ForegroundColor Gray
Write-Host "  - MongoDB (Database + sample data)" -ForegroundColor Gray
Write-Host "  - Redis (Cache & distributed locks)" -ForegroundColor Gray
Write-Host "  - Order Worker (Java - Core processing)" -ForegroundColor Gray
Write-Host "  - Product API (Go - Product catalog)" -ForegroundColor Gray
Write-Host "  - Customer API (Go - Customer management)" -ForegroundColor Gray
Write-Host "`nServices NOT included:" -ForegroundColor Yellow
Write-Host "  - Order API (Frontend bridge)" -ForegroundColor Gray
Write-Host "  - Nginx Frontend (Web interface)" -ForegroundColor Gray

try {
    docker-compose up -d
    Write-Host "`nBackend services started successfully!" -ForegroundColor Green
} catch {
    Write-Host "`nFailed to start backend services" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Wait for services to be ready
Write-Host "`nWaiting for services to be healthy..." -ForegroundColor Yellow
$maxWait = 45
$waited = 0
$checkInterval = 5

do {
    Start-Sleep -Seconds $checkInterval
    $waited += $checkInterval
    
    try {
        $status = docker-compose ps --format "table {{.Name}}\t{{.Status}}"
        $healthyCount = ($status | Select-String "healthy").Count
        
        # Count services that should have health checks (excluding kafka-setup which is one-time)
        $totalServices = ($status | Select-String -NotMatch "kafka-setup").Count - 1  # -1 for header line
        
        if ($healthyCount -eq $totalServices) {
            Write-Host "All services are healthy!" -ForegroundColor Green
            break
        } else {
            Write-Host "Healthy: $healthyCount/$totalServices (waited ${waited}s)" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "Checking service status... (waited ${waited}s)" -ForegroundColor Yellow
    }
} while ($waited -lt $maxWait)

if ($waited -ge $maxWait) {
    Write-Host "Timeout reached. Some services may still be starting." -ForegroundColor Yellow
}

# Check service status
Write-Host "`nChecking service status..." -ForegroundColor Cyan
try {
    $status = docker-compose ps
    Write-Host $status -ForegroundColor Gray
} catch {
    Write-Host "Could not retrieve service status" -ForegroundColor Yellow
}

# Verify critical services
Write-Host "`nHealth checks..." -ForegroundColor Cyan

$healthChecks = @(
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
Write-Host "`nBackend Profile Ready!" -ForegroundColor Green
Write-Host "======================" -ForegroundColor Green

Write-Host "`nAvailable Endpoints:" -ForegroundColor Cyan
Write-Host "  Product API: http://localhost:8081" -ForegroundColor Gray
Write-Host "  Customer API: http://localhost:8082" -ForegroundColor Gray
Write-Host "  MongoDB: mongodb://localhost:27017" -ForegroundColor Gray
Write-Host "  Redis: redis://localhost:6379" -ForegroundColor Gray
Write-Host "  Kafka: kafka://localhost:9092" -ForegroundColor Gray

Write-Host "`nTesting Options:" -ForegroundColor Cyan
Write-Host "  Postman Collection: Import postman/*.json files" -ForegroundColor Gray
Write-Host "  Automated Tests: scripts/test-final-system.ps1" -ForegroundColor Gray
Write-Host "  Manual Kafka:" -ForegroundColor Gray
Write-Host "    echo '{`"orderId`":`"test`",`"customerId`":`"customer-1`",`"products`":[{`"productId`":`"product-1`"}]}' | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders" -ForegroundColor DarkGray

Write-Host "`nMonitoring:" -ForegroundColor Cyan
Write-Host "  Service Status: docker-compose ps" -ForegroundColor Gray
Write-Host "  Order Worker Logs: docker-compose logs -f order-worker" -ForegroundColor Gray
Write-Host "  MongoDB Data: docker-compose exec mongo mongosh orders --eval `"db.orders.find().forEach(printjson)`"" -ForegroundColor Gray

Write-Host "`nTo add Frontend later:" -ForegroundColor Yellow
Write-Host "  docker-compose --profile frontend up -d" -ForegroundColor Gray

Write-Host "`nTo stop services:" -ForegroundColor Yellow
Write-Host "  docker-compose down" -ForegroundColor Gray

Write-Host "`nBackend Profile deployment completed!" -ForegroundColor Green