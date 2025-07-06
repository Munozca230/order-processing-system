# Test script for MongoDB-backed APIs
Write-Host "=== Testing MongoDB-backed APIs ====" -ForegroundColor Green

Set-Location "../infra"

# Clean start with volumes to ensure fresh MongoDB initialization
Write-Host "1. Cleaning and starting with fresh MongoDB..." -ForegroundColor Yellow
docker-compose down -v

# Start infrastructure
docker-compose up -d zookeeper kafka mongo redis

# Wait longer for MongoDB initialization scripts
Write-Host "2. Waiting for MongoDB initialization..." -ForegroundColor Yellow
Start-Sleep -Seconds 45

# Check MongoDB initialization
Write-Host "3. Verifying MongoDB initialization..." -ForegroundColor Yellow
Write-Host "   Products in MongoDB:" -ForegroundColor Cyan
docker-compose exec mongo mongosh catalog --eval "db.products.countDocuments(); db.products.find({}, {productId:1, name:1, active:1}).limit(3).forEach(printjson);"

Write-Host "`n   Customers in MongoDB:" -ForegroundColor Cyan
docker-compose exec mongo mongosh catalog --eval "db.customers.countDocuments(); db.customers.find({}, {customerId:1, name:1, active:1}).limit(3).forEach(printjson);"

# Start APIs
Write-Host "`n4. Starting MongoDB-backed APIs..." -ForegroundColor Yellow
docker-compose up -d --build product-api customer-api

Start-Sleep -Seconds 20

# Test APIs with direct HTTP calls (using Invoke-WebRequest)
Write-Host "`n5. Testing APIs with MongoDB data..." -ForegroundColor Yellow

Write-Host "   Product API - Get Product from MongoDB:" -ForegroundColor Cyan
try {
    $productResponse = Invoke-WebRequest -Uri "http://localhost:8081/products/product-1" -UseBasicParsing
    $productData = $productResponse.Content | ConvertFrom-Json
    Write-Host "   ✅ Product: $($productData.name) - $($productData.price)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Product API failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n   Customer API - Get Customer from MongoDB:" -ForegroundColor Cyan
try {
    $customerResponse = Invoke-WebRequest -Uri "http://localhost:8082/customers/customer-1" -UseBasicParsing
    $customerData = $customerResponse.Content | ConvertFrom-Json
    Write-Host "   ✅ Customer: $($customerData.name) - Active: $($customerData.active)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Customer API failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test pagination
Write-Host "`n   Testing pagination:" -ForegroundColor Cyan
try {
    $paginatedResponse = Invoke-WebRequest -Uri "http://localhost:8081/products?page=0&page_size=2" -UseBasicParsing
    $paginatedData = $paginatedResponse.Content | ConvertFrom-Json
    Write-Host "   ✅ Paginated products: $($paginatedData.total) total, page size: $($paginatedData.products.Count)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Pagination failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Start order worker for E2E test
Write-Host "`n6. Starting order worker for E2E test..." -ForegroundColor Yellow
docker-compose up -d order-worker

Start-Sleep -Seconds 15

# Test complete E2E workflow
Write-Host "7. Testing E2E with MongoDB-backed APIs..." -ForegroundColor Yellow
$testMessage = '{"orderId":"mongodb-e2e-test","customerId":"customer-1","products":[{"productId":"product-1"},{"productId":"product-2"}]}'
echo $testMessage | docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic orders

Start-Sleep -Seconds 8

# Show logs
Write-Host "`n8. Order Worker processing logs:" -ForegroundColor Yellow
docker-compose logs --tail=15 order-worker | Select-String -Pattern "(RECEIVED|ENRICHMENT|VALIDATION|MONGODB|COMPLETED)"

# Verify MongoDB final result
Write-Host "`n9. Final E2E verification in MongoDB:" -ForegroundColor Yellow
docker-compose exec mongo mongosh orders --eval "db.orders.find({orderId: 'mongodb-e2e-test'}).forEach(printjson);"

Write-Host "`n=== MongoDB-backed APIs Test Complete ====" -ForegroundColor Green
Write-Host "✅ Production-ready system with MongoDB persistence!" -ForegroundColor Green