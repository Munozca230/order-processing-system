# Final test - comprehensive system validation
Write-Host "=== Final System Test ====" -ForegroundColor Green

Set-Location "../infra"

# Test multiple scenarios
Write-Host "1. Sending test message..." -ForegroundColor Yellow
$testMessage = '{"orderId":"final-validation","customerId":"customer-1","products":["product-1","product-2"]}'
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Start-Sleep -Seconds 8

# Show recent logs
Write-Host "2. Processing logs:" -ForegroundColor Yellow
docker-compose logs --tail=20 order-worker

# Check MongoDB directly
Write-Host "3. MongoDB verification:" -ForegroundColor Yellow
docker-compose exec mongo mongosh orders --eval "db.orders.find({orderId: 'final-validation'}).forEach(printjson);"

Write-Host "=== Final Test Complete ====" -ForegroundColor Green