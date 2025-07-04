# Debug test script with enhanced logging
Write-Host "=== Debug Test: Enhanced Logging ====" -ForegroundColor Green

# Change to infra directory  
Set-Location "../infra"

# Rebuild services
Write-Host "1. Rebuilding services..." -ForegroundColor Yellow
docker-compose down
docker-compose up -d --build

# Wait for services to be ready
Write-Host "2. Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Send test message
Write-Host "3. Sending test message..." -ForegroundColor Yellow
$testMessage = '{"orderId":"debug-test","customerId":"customer-1","products":["product-1","product-2"]}'
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

# Wait for processing
Write-Host "4. Waiting for processing..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Show detailed logs
Write-Host "5. Order Worker logs:" -ForegroundColor Yellow
docker-compose logs --tail=50 order-worker

# Check MongoDB
Write-Host "6. MongoDB check:" -ForegroundColor Yellow
docker-compose exec mongo mongosh --eval "use orders; db.orders.find().pretty();"

Write-Host "=== Debug Test Complete ====" -ForegroundColor Green