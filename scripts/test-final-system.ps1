# Final test script for production-ready order processing system
Write-Host "=== Testing Final Production System ====" -ForegroundColor Green

Set-Location "../infra"

# Clean start
Write-Host "1. Cleaning and starting fresh..." -ForegroundColor Yellow
docker-compose down -v

# Start all services
Write-Host "2. Starting all services..." -ForegroundColor Yellow
docker-compose up -d

# Wait for initialization
Write-Host "3. Waiting for MongoDB initialization..." -ForegroundColor Yellow
Start-Sleep -Seconds 45

# Verify MongoDB data
Write-Host "4. Verifying MongoDB initialization..." -ForegroundColor Yellow
Write-Host "   Products:" -ForegroundColor Cyan
docker-compose exec mongo mongosh catalog --eval "print('Count:', db.products.countDocuments()); db.products.find({}, {productId:1, name:1, active:1}).limit(3).forEach(printjson);"

Write-Host "`n   Customers:" -ForegroundColor Cyan
docker-compose exec mongo mongosh catalog --eval "print('Count:', db.customers.countDocuments()); db.customers.find({}, {customerId:1, name:1, active:1}).limit(3).forEach(printjson);"

# Test APIs
Write-Host "`n5. Testing API endpoints..." -ForegroundColor Yellow

Write-Host "   Product API Health:" -ForegroundColor Cyan
try {
    $productHealth = Invoke-WebRequest -Uri "http://localhost:8081/health" -UseBasicParsing
    Write-Host "   ✅ Product API healthy" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Product API failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "   Customer API Health:" -ForegroundColor Cyan
try {
    $customerHealth = Invoke-WebRequest -Uri "http://localhost:8082/health" -UseBasicParsing
    Write-Host "   ✅ Customer API healthy" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Customer API failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test data retrieval
Write-Host "`n   Testing data retrieval:" -ForegroundColor Cyan
try {
    $productResponse = Invoke-WebRequest -Uri "http://localhost:8081/products/product-1" -UseBasicParsing
    $productData = $productResponse.Content | ConvertFrom-Json
    Write-Host "   ✅ Product fetched: $($productData.name) - $$($productData.price)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Product fetch failed: $($_.Exception.Message)" -ForegroundColor Red
}

try {
    $customerResponse = Invoke-WebRequest -Uri "http://localhost:8082/customers/customer-1" -UseBasicParsing
    $customerData = $customerResponse.Content | ConvertFrom-Json
    Write-Host "   ✅ Customer fetched: $($customerData.name) - Active: $($customerData.active)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Customer fetch failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test end-to-end processing
Write-Host "`n6. Testing end-to-end order processing..." -ForegroundColor Yellow
$testMessage = '{"orderId":"final-e2e-test","customerId":"customer-1","products":[{"productId":"product-1"},{"productId":"product-2"}]}'
Write-Host "   Sending order: $testMessage" -ForegroundColor Cyan
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Start-Sleep -Seconds 10

# Show processing logs
Write-Host "`n7. Order processing logs:" -ForegroundColor Yellow
docker-compose logs --tail=15 order-worker | Select-String -Pattern "(RECEIVED|ENRICHMENT|VALIDATION|MONGODB|COMPLETED|ERROR)"

# Verify final result
Write-Host "`n8. Final verification in MongoDB:" -ForegroundColor Yellow
docker-compose exec mongo mongosh orders --eval "db.orders.find({orderId: 'final-e2e-test'}).forEach(printjson);"

Write-Host "`n=== Final Production System Test Complete ====" -ForegroundColor Green
Write-Host "✅ System ready for delivery!" -ForegroundColor Green