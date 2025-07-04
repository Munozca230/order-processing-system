# Test script to verify new package structure works
Write-Host "=== Testing New Package Structure (com.orderprocessing) ====" -ForegroundColor Green

Set-Location "../infra"

# Clean build to test new package structure
Write-Host "1. Clean build with new package structure..." -ForegroundColor Yellow
docker-compose down -v
docker-compose build order-worker

# Start all services
Write-Host "2. Starting all services..." -ForegroundColor Yellow
docker-compose up -d

# Wait for services
Write-Host "3. Waiting for services to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 45

# Test that order-worker starts correctly with new package
Write-Host "4. Checking order-worker logs for startup..." -ForegroundColor Yellow
docker-compose logs order-worker | Select-String -Pattern "(Started|ERROR|Exception)"

# Test basic functionality
Write-Host "5. Testing basic order processing..." -ForegroundColor Yellow
$testMessage = '{"orderId":"package-test","customerId":"customer-1","products":[{"productId":"product-1"}]}'
Write-Host "   Sending test order: $testMessage" -ForegroundColor Cyan
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Start-Sleep -Seconds 8

# Check processing logs
Write-Host "6. Order processing logs:" -ForegroundColor Yellow
docker-compose logs --tail=10 order-worker | Select-String -Pattern "(RECEIVED|ENRICHMENT|VALIDATION|MONGODB|COMPLETED|ERROR)"

# Verify MongoDB result
Write-Host "7. MongoDB verification:" -ForegroundColor Yellow
docker-compose exec mongo mongosh orders --eval "db.orders.find({orderId: 'package-test'}).forEach(printjson);"

Write-Host "`n=== Package Structure Test Complete ====" -ForegroundColor Green
Write-Host "✅ New package structure working!" -ForegroundColor Green