# Simple test script with better error handling
Write-Host "=== Simple Debug Test ====" -ForegroundColor Green

# Change to infra directory  
Set-Location "../infra"

# Clean start
Write-Host "1. Cleaning previous containers..." -ForegroundColor Yellow
docker-compose down -v

# Start infrastructure first
Write-Host "2. Starting infrastructure..." -ForegroundColor Yellow
docker-compose up -d zookeeper kafka mongo redis

# Wait for infrastructure
Write-Host "3. Waiting for infrastructure to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

# Start applications
Write-Host "4. Starting applications..." -ForegroundColor Yellow
docker-compose up -d product-api customer-api

Start-Sleep -Seconds 10

# Start order worker
Write-Host "5. Starting order worker..." -ForegroundColor Yellow
docker-compose up -d order-worker

Start-Sleep -Seconds 10

# Check status
Write-Host "6. Checking service status..." -ForegroundColor Yellow
docker-compose ps

# Send test message
Write-Host "7. Sending test message..." -ForegroundColor Yellow
$testMessage = '{"orderId":"debug-simple","customerId":"customer-1","products":["product-1"]}'
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

# Wait for processing
Write-Host "8. Waiting for processing..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Show detailed logs
Write-Host "9. Order Worker logs:" -ForegroundColor Yellow
docker-compose logs order-worker

Write-Host "=== Simple Debug Test Complete ====" -ForegroundColor Green